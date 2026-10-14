package gewe

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kilimro/sillyGirl/core"
	"github.com/kilimro/sillyGirl/core/storage"
)

// Gewe 第三方微信适配器
//
// 架构：执行节点 + 事件回调（webhook 被动模式，与 yyw 同类）。
// 接入不需要 web 扫码登录：用户在 GeWe 后台扫码登录微信拿到 appId，
// 再把 token / app_id / api_base / 回调地址填进本项目配置即可。
//
// 接收：GeWe 服务器将消息以 HTTP POST(JSON) 推送到配置的回调地址。
//   兼容两版回调结构：
//     v1.0（嵌套）: { "TypeName":"AddMsg","Appid","Wxid","Data":{ "MsgType":1,"FromUserName":{"string":..},... } }
//     v2.0（扁平）: { "appid","wxid","content","fromUser","toUser","msgType":"TEXT","newMsgId","isSelf","createTime" }
//   必须实现：
//     - 去重 Key = Appid + NewMsgId
//     - isSelf（自己发送的消息）必须忽略，防自回复循环
//     - 回调必须在 3 秒内返回（可返回空字符串）
// 发送：POST /gewe/v2/api/message/post{Text,Image,File,Voice,Video}，Header X-GEWE-TOKEN。

const (
	defaultAPIBase = "https://api.geweapi.com"
	dedupTTL       = 10 * time.Minute
)

var (
	cqAnyPattern = regexp.MustCompile(`\[CQ:(\w+),([^\]]+)\]`)

	gewe = core.MakeBucket("gewe")

	geweAdapter     *core.Factory
	geweAdapterOnce sync.Once
	runtimeMu       sync.RWMutex
	apiBase         string
	token           string
	appID           string
	selfWxid        string
	callbackURL     string

	dedupMu   sync.Mutex
	dedupSeen = make(map[string]time.Time)

	callbackMu         sync.Mutex
	lastCallbackURL    string
	callbackScheduled  = false
)

func init() {
	core.GinApi(core.POST, "/api/gewe/webhook", receiveWebhook)
	refreshConfig()
	for _, key := range []string{"api_base", "token", "app_id", "enable", "self_wxid"} {
		storage.Watch(gewe, key, func(old, new, key string) *storage.Final {
			go refreshConfig()
			return nil
		})
	}
	go func() {
		time.Sleep(3 * time.Second)
		refreshConfig()
	}()
}

func enabled() bool {
	return gewe.GetBool("enable", false)
}

func refreshConfig() {
	runtimeMu.Lock()
	apiBase = strings.TrimRight(gewe.GetString("api_base", defaultAPIBase), "/")
	token = gewe.GetString("token")
	appID = gewe.GetString("app_id")
	selfWxid = gewe.GetString("self_wxid")
	callbackURL = gewe.GetString("callback_url")
	runtimeMu.Unlock()
	if enabled() && token != "" && apiBase != "" {
		initAdapter()
		// 方案B：自动设置 Gewe 回调地址（setCallback）。失败不阻塞，可后台手动（方案A兼容）。
		if callbackURL != "" {
			go ensureCallback()
		}
	}
}

func ensureCallback() {
	runtimeMu.RLock()
	base, tk := apiBase, token
	cb := callbackURL
	runtimeMu.RUnlock()
	if base == "" || tk == "" || cb == "" {
		return
	}
	callbackMu.Lock()
	if callbackScheduled {
		callbackMu.Unlock()
		return
	}
	callbackScheduled = true
	callbackMu.Unlock()

	if err := callSetCallback(base, tk, cb); err != nil {
		core.Logs.Warn("gewe 自动设置回调地址失败（可后台手动填写）: %v", err)
	}
	callbackMu.Lock()
	callbackScheduled = false
	callbackMu.Unlock()
}

// callSetCallback 调用 Gewe 登录模块 setCallback 接口设置回调地址。
// POST /gewe/v2/api/login/setCallback, Header X-GEWE-TOKEN, Body {token, callbackUrl}
func callSetCallback(apiBase, tk, callbackURL string) error {
	payload := map[string]interface{}{"token": tk, "callbackUrl": callbackURL}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiBase+"/gewe/v2/api/login/setCallback", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GEWE-TOKEN", tk)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	core.Logs.Info("gewe 设置回调地址(%s) resp=%s", callbackURL, string(respBody))
	return nil
}

func initAdapter() {
	geweAdapterOnce.Do(func() {
		geweAdapter = &core.Factory{}
		geweAdapter.Init("gewe", "gewe-bot", nil)
		geweAdapter.SetReplyHandler(func(msg map[string]interface{}) string {
			return doReply(msg)
		})
	})
}

// ----- 接收（webhook） -----

type geweIncoming struct {
	appid     string
	wxid      string
	msgType   string // 归一化类型：text/image/voice/video/file/other
	fromUser  string
	toUser    string
	content   string
	newMsgID  string
	msgID     string // 微信数字 msgId（v1.0 Data.MsgId），语音下载用
	isSelf    bool
	createTime int64
}

func receiveWebhook(c *gin.Context) {
	if !enabled() {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "gewe disabled"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "read body error"})
		return
	}
	cb, err := parseCallback(body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "bad json"})
		return
	}
	// 掉线等系统事件不进入核心，但需要应答成功
	if cb.msgType == "other" {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	// 忽略自己发送的消息（防自回复循环）
	if cb.isSelf {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	// 去重：Appid + NewMsgId
	if cb.newMsgID != "" && isDuplicate(cb.appid, cb.newMsgID) {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	// 群聊判断：fromUser 以 @chatroom 结尾
	isGroup := strings.HasSuffix(cb.fromUser, "@chatroom")

	userID := cb.fromUser
	chatID := cb.fromUser
	if !isGroup {
		userID = cb.fromUser
		chatID = cb.fromUser
	} else {
		// 群消息：chatID = 群 wxid；尝试从 content 前缀解析群成员 wxid
		chatID = cb.fromUser
		if member, rest := splitGroupContent(cb.content); member != "" {
			userID = member
			cb.content = rest
		}
	}
	// 自动记录机器人 wxid
	if cb.wxid != "" {
		runtimeMu.Lock()
		if selfWxid == "" {
			selfWxid = cb.wxid
		}
		runtimeMu.Unlock()
	}

	// 内容转换为供核心消费的文本（多媒体生成 CQ 码）
	text := toCoreContent(cb.msgType, cb.content, cb.msgID)
	if text == "" {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	core.Logs.Info("gewe收到消息 from=%s chat=%s group=%v type=%s", userID, chatID, isGroup, cb.msgType)

	initAdapter()
	geweAdapter.Receive(map[string]interface{}{
		core.USER_ID:    userID,
		core.CHAT_ID:    chatID,
		core.CONETNT:    text,
		"user_name":     userID,
		"chat_name":     chatID,
		"gewe_is_group": isGroup,
	})

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// parseCallback 兼容 v1.0 嵌套结构与 v2.0 扁平结构
func parseCallback(body []byte) (*geweIncoming, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	cb := &geweIncoming{}
	cb.appid = str(raw["appid"])
	if cb.appid == "" {
		cb.appid = str(raw["Appid"])
	}
	cb.wxid = str(raw["wxid"])
	if cb.wxid == "" {
		cb.wxid = str(raw["Wxid"])
	}

	// 取 Data 子树（v1.0 存在，v2.0 无则用顶层）
	data := raw
	if d, ok := raw["Data"].(map[string]interface{}); ok {
		data = d
	}

	// msgType：优先数字（v1.0），否则字符串（v2.0）
	msgType := ""
	if n, ok := data["MsgType"].(float64); ok {
		msgType = normalizeMsgType(int(n))
	} else {
		msgType = normalizeMsgTypeString(str(data["msgType"]))
	}
	cb.msgType = msgType

	cb.fromUser = nestedStr(data, "FromUserName", "string")
	if cb.fromUser == "" {
		cb.fromUser = str(data["fromUser"])
	}
	cb.toUser = nestedStr(data, "ToUserName", "string")
	if cb.toUser == "" {
		cb.toUser = str(data["toUser"])
	}
	cb.content = nestedStr(data, "Content", "string")
	if cb.content == "" {
		cb.content = str(data["content"])
	}
	cb.newMsgID = str(data["newMsgId"])
	if cb.newMsgID == "" {
		cb.newMsgID = itoa(data["NewMsgId"])
	}
	// 微信数字 msgId（Data.MsgId），语音下载 API 需要
	if n, ok := data["MsgId"].(float64); ok && n > 0 {
		cb.msgID = itoa(data["MsgId"])
	}
	if cb.msgID == "" {
		cb.msgID = str(data["msgId"])
	}
	cb.createTime = int64(num(data["createTime"]))
	if cb.createTime == 0 {
		cb.createTime = int64(num(data["CreateTime"]))
	}

	// isSelf：v2.0 显式字段，或 fromUser == wxid
	if b, ok := data["isSelf"].(bool); ok {
		cb.isSelf = b
	} else {
		cb.isSelf = cb.fromUser != "" && cb.wxid != "" && cb.fromUser == cb.wxid
	}
	return cb, nil
}

// normalizeMsgType 数字类型归一化（v1.0）
func normalizeMsgType(t int) string {
	switch t {
	case 1:
		return "text"
	case 3:
		return "image"
	case 34:
		return "voice"
	case 43:
		return "video"
	case 49:
		// 官方要点：MsgType=49 再解析 XML（appmsg），XML 是核心。类型由 toCoreContent 解析 XML 细分。
		return "appmsg"
	default:
		return "other"
	}
}

// normalizeMsgTypeString 字符串类型归一化（v2.0）
func normalizeMsgTypeString(s string) string {
	switch strings.ToUpper(s) {
	case "TEXT":
		return "text"
	case "IMAGE", "PIC":
		return "image"
	case "VOICE", "RECORD", "AUDIO":
		return "voice"
	case "VIDEO":
		return "video"
	case "FILE":
		return "file"
	default:
		return "other"
	}
}

// toCoreContent 将消息内容转换为供核心消费的文本（多媒体转 CQ 码）
// 多媒体 CQ 码会附带 GeWe 下载参数（aeskey/size/msgId/xml），供依赖模块按需调 downloadCdn/downloadVoice。
func toCoreContent(msgType, content, msgID string) string {
	switch msgType {
	case "text":
		return strings.TrimSpace(content)
	case "image":
		u := extractURL(content)
		if u == "" {
			return ""
		}
		aeskey := extractAttr(content, "aeskey")
		size := firstNonEmpty(extractAttr(content, "total_size"), extractAttr(content, "totalsize"))
		suffix := firstNonEmpty(extractAttr(content, "suffix"), "jpg")
		return fmt.Sprintf("[CQ:image,file=%s,aeskey=%s,size=%s,suffix=%s]", u, aeskey, size, suffix)
	case "voice":
		u := extractURL(content)
		if u == "" {
			return ""
		}
		b64 := base64.StdEncoding.EncodeToString([]byte(content))
		return fmt.Sprintf("[CQ:record,file=%s,msgId=%s,xml=%s]", u, msgID, b64)
	case "video":
		u := extractURL(content)
		if u == "" {
			return ""
		}
		aeskey := extractAttr(content, "aeskey")
		size := firstNonEmpty(extractAttr(content, "total_size"), extractAttr(content, "totalsize"))
		return fmt.Sprintf("[CQ:video,file=%s,aeskey=%s,size=%s,type=4]", u, aeskey, size)
	case "file":
		u := extractURL(content)
		if u == "" {
			return ""
		}
		aeskey := extractAttr(content, "aeskey")
		size := firstNonEmpty(extractAttr(content, "total_size"), extractAttr(content, "totalsize"))
		suffix := extractAttr(content, "suffix")
		return fmt.Sprintf("[CQ:file,file=%s,aeskey=%s,size=%s,suffix=%s]", u, aeskey, size, suffix)
	case "appmsg":
		// MsgType=49：解析 XML（appmsg.type）细分类型，提取标题进核心。
		appType, title, u := parseAppMsg(content)
		switch appType {
		case "file":
			if u != "" {
				return fmt.Sprintf("[CQ:file,file=%s]", u)
			}
			if title != "" {
				return "[文件]" + title
			}
		case "link":
			if u != "" {
				return title + " " + u
			}
			if title != "" {
				return "[链接]" + title
			}
		case "mini_program", "quote", "transfer", "red_packet", "finder":
			if title != "" {
				return "[" + appType + "]" + title
			}
		default:
			if title != "" {
				return title
			}
		}
		return ""
	}
	return ""
}

// parseAppMsg 解析 MsgType=49 的 appmsg XML，返回 (类型, 标题, 链接/文件URL)。
// 参考官方判断：appmsg.type=5 链接、6 文件完成、74 文件通知、33/36 小程序、57 引用、
// 51 视频号、2000 转账、2001 红包。
func parseAppMsg(content string) (kind, title, url string) {
	kind = ""
	title = strings.TrimSpace(extractXML(content, "title"))
	u := strings.TrimSpace(extractXML(content, "url"))
	typ := strings.TrimSpace(extractXML(content, "type"))
	switch typ {
	case "6", "74":
		kind = "file"
	case "5":
		kind = "link"
	case "33", "36":
		kind = "mini_program"
	case "57":
		kind = "quote"
	case "2000":
		kind = "transfer"
	case "2001":
		kind = "red_packet"
	case "51":
		kind = "finder"
	default:
		kind = "appmsg"
	}
	return kind, title, u
}

// extractXML 从 XML 中提取指定标签的文本（兼容 CDATA）。
func extractXML(xml, tag string) string {
	re := regexp.MustCompile("<" + tag + `[^>]*>(?:<!\[CDATA\[)?([^<]*?)(?:\]\]>)?</` + tag + `>`)
	if m := re.FindStringSubmatch(xml); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// extractURL 从多媒体 XML 中提取 CDN URL
func extractURL(xml string) string {
	if xml == "" {
		return ""
	}
	// <cdnurl>...</cdnurl> 或 cdnurl="..."
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`<cdnurl[^>]*>([^<]+)</cdnurl>`),
		regexp.MustCompile(`cdnurl="([^"]+)"`),
		regexp.MustCompile(`<url[^>]*>([^<]+)</url>`),
		regexp.MustCompile(`url="([^"]+)"`),
	} {
		if m := re.FindStringSubmatch(xml); len(m) > 1 && m[1] != "" {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}

// extractAttr 从 XML 片段中提取指定属性值（不区分大小写）。
func extractAttr(xml, attr string) string {
	if xml == "" {
		return ""
	}
	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(attr) + `="([^"]*)"`)
	if m := re.FindStringSubmatch(xml); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// splitGroupContent 解析群消息 content 前缀 "成员wxid（昵称）...:\n正文"。
// Gewe 群消息 content 前缀带发送者 wxid（不含 @；含 @chatroom 结尾的是群 id，不会出现在成员前缀）。
func splitGroupContent(content string) (member, rest string) {
	content = strings.TrimPrefix(content, "\u2005")
	idx := strings.Index(content, ":\n")
	if idx <= 0 {
		return "", content
	}
	prefix := content[:idx]
	if strings.Contains(prefix, "@chatroom") {
		return "", content
	}
	if m := groupSenderPattern.FindString(prefix); m != "" {
		return m, strings.TrimSpace(content[idx+2:])
	}
	return "", content
}

var groupSenderPattern = regexp.MustCompile(`^wxid_[A-Za-z0-9_]+`)

// isDuplicate 消息去重
func isDuplicate(appid, newMsgID string) bool {
	key := appid + "\x00" + newMsgID
	now := time.Now()
	dedupMu.Lock()
	defer dedupMu.Unlock()
	if _, ok := dedupSeen[key]; ok {
		return true
	}
	// 清理过期
	for k, t := range dedupSeen {
		if now.Sub(t) > dedupTTL {
			delete(dedupSeen, k)
		}
	}
	dedupSeen[key] = now
	return false
}

// ----- 发送 -----

func doReply(msg map[string]interface{}) string {
	runtimeMu.RLock()
	base := apiBase
	tk := token
	app := appID
	runtimeMu.RUnlock()
	if base == "" || tk == "" {
		core.Logs.Warn("gewe: token/api_base 未配置，无法发送消息")
		return ""
	}
	if app == "" {
		core.Logs.Warn("gewe: app_id 未配置，无法发送消息")
		return ""
	}
	toWxid := stringVal(msg[core.CHAT_ID])
	text := strings.TrimSpace(stringVal(msg[core.CONETNT]))
	if text == "" {
		return ""
	}
	for _, seg := range splitSegments(text) {
		switch seg.kind {
		case "image":
			sendGewe(base, tk, app, "/message/postImage", map[string]interface{}{"appId": app, "toWxid": toWxid, "imgUrl": seg.value})
		case "audio", "voice", "record":
			// Gewe 语音仅支持 silk 格式（voiceUrl 需为 .silk）
			sendGewe(base, tk, app, "/message/postVoice", map[string]interface{}{"appId": app, "toWxid": toWxid, "voiceUrl": seg.value})
		case "video":
			sendGewe(base, tk, app, "/message/postVideo", map[string]interface{}{"appId": app, "toWxid": toWxid, "videoUrl": seg.value})
		case "file":
			sendGewe(base, tk, app, "/message/postFile", map[string]interface{}{"appId": app, "toWxid": toWxid, "fileUrl": seg.value})
		default:
			sendGewe(base, tk, app, "/message/postText", map[string]interface{}{"appId": app, "toWxid": toWxid, "content": seg.value})
		}
	}
	return ""
}

func sendGewe(apiBase, tk, app, path string, payload map[string]interface{}) {
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiBase+"/gewe/v2/api"+path, bytes.NewReader(b))
	if err != nil {
		core.Logs.Warn("gewe 构造请求失败 %s: %v", path, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GEWE-TOKEN", tk)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		core.Logs.Warn("gewe发送%s失败: %v", path, err)
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	core.Logs.Info("gewe发送%s to=%s resp=%s", path, stringVal(payload["toWxid"]), string(respBody))
}

type geweSegment struct {
	kind  string
	value string
}

func splitSegments(text string) []geweSegment {
	matches := cqAnyPattern.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return []geweSegment{{kind: "text", value: text}}
	}
	var segs []geweSegment
	last := 0
	for _, m := range matches {
		if m[0] > last {
			if t := strings.TrimSpace(text[last:m[0]]); t != "" {
				segs = append(segs, geweSegment{kind: "text", value: t})
			}
		}
		cqType := text[m[2]:m[3]]
		params := parseCQParams(text[m[4]:m[5]])
		url := params["file"]
		if url == "" {
			url = params["url"]
		}
		if url == "" {
			last = m[1]
			continue
		}
		kind := "text"
		switch strings.ToLower(cqType) {
		case "image", "pic":
			kind = "image"
		case "record", "audio", "voice":
			kind = "voice"
		case "video":
			kind = "video"
		case "file":
			kind = "file"
		}
		if kind != "text" {
			segs = append(segs, geweSegment{kind: kind, value: url})
		}
		last = m[1]
	}
	if last < len(text) {
		if t := strings.TrimSpace(text[last:]); t != "" {
			segs = append(segs, geweSegment{kind: "text", value: t})
		}
	}
	return segs
}

func parseCQParams(raw string) map[string]string {
	params := map[string]string{}
	for _, item := range strings.Split(raw, ",") {
		key, value, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		params[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return params
}

// ----- 工具 -----

func str(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func itoa(v interface{}) string {
	if v == nil {
		return ""
	}
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", n)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func num(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case string:
		f := float64(0)
		fmt.Sscanf(n, "%f", &f)
		return f
	}
	return 0
}

func nestedStr(m map[string]interface{}, outer, inner string) string {
	if o, ok := m[outer].(map[string]interface{}); ok {
		if s, ok := o[inner].(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func stringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}
