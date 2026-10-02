package yyw

import (
	"bytes"
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

var cqAnyPattern = regexp.MustCompile(`\[CQ:(\w+),([^\]]+)\]`)

var yyw = core.MakeBucket("yyw")

var (
	yywAdapter     *core.Factory
	yywAdapterOnce sync.Once
	runtimeMu      sync.RWMutex
	apiBase        string
	selfid         string
	authorization  string
)

type yywIncoming struct {
	Event         string      `json:"event"`
	RobotWxid     string      `json:"robot_wxid"`
	Type          int         `json:"type"`
	FromWxid      string      `json:"from_wxid"`
	FromName      string      `json:"from_name"`
	FinalFromWxid string      `json:"final_from_wxid"`
	FinalFromName string      `json:"final_from_name"`
	ToWxid        string      `json:"to_wxid"`
	Msg           interface{} `json:"msg"`
}

func init() {
	core.GinApi(core.POST, "/api/yyw/webhook", receiveWebhook)
	core.GinApi(core.GET, "/api/yyw/webhook", receiveWebhook)
	refreshConfig()
	storage.Watch(yyw, "api_base", func(old, new, key string) *storage.Final {
		go refreshConfig()
		return nil
	})
	storage.Watch(yyw, "selfid", func(old, new, key string) *storage.Final {
		go refreshConfig()
		return nil
	})
	storage.Watch(yyw, "authorization", func(old, new, key string) *storage.Final {
		go refreshConfig()
		return nil
	})
	storage.Watch(yyw, "enable", func(old, new, key string) *storage.Final {
		go refreshConfig()
		return nil
	})
	go func() {
		time.Sleep(3 * time.Second)
		refreshConfig()
	}()
}

func enabled() bool {
	return yyw.GetBool("enable", false)
}

func refreshConfig() {
	runtimeMu.Lock()
	apiBase = strings.TrimRight(yyw.GetString("api_base"), "/")
	selfid = yyw.GetString("selfid")
	authorization = yyw.GetString("authorization")
	runtimeMu.Unlock()
	if enabled() && apiBase != "" {
		initAdapter()
	}
}

func initAdapter() {
	yywAdapterOnce.Do(func() {
		yywAdapter = &core.Factory{}
		yywAdapter.Init("yyw", "yyw-bot", nil)
		yywAdapter.SetReplyHandler(func(msg map[string]interface{}) string {
			return doReply(msg)
		})
	})
}

func receiveWebhook(c *gin.Context) {
	if !enabled() {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "yyw disabled"})
		return
	}
	var incoming yywIncoming
	if err := c.ShouldBindJSON(&incoming); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "bad json"})
		return
	}
	if incoming.Event != "EventGroupMsg" && incoming.Event != "EventFriendMsg" {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	text := extractText(incoming.Msg)
	if text == "" {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	chatID := incoming.FromWxid
	isGroup := incoming.Event == "EventGroupMsg"
	if !isGroup {
		chatID = incoming.FinalFromWxid
	}
	userID := incoming.FinalFromWxid
	if userID == "" {
		userID = incoming.FromWxid
	}

	// 自动记录机器人 wxid
	if incoming.RobotWxid != "" {
		runtimeMu.Lock()
		if selfid == "" {
			selfid = incoming.RobotWxid
		}
		runtimeMu.Unlock()
	}

	core.Logs.Info("yyw收到消息 from=%s chat=%s group=%v text=%s", userID, chatID, isGroup, text)

	initAdapter()
	yywAdapter.Receive(map[string]interface{}{
		core.USER_ID:   userID,
		core.CHAT_ID:   chatID,
		core.CONETNT:   text,
		"user_name":    incoming.FinalFromName,
		"chat_name":    incoming.FromName,
		"yyw_is_group": isGroup,
	})

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok"})
}

func extractText(msg interface{}) string {
	switch v := msg.(type) {
	case string:
		return strings.TrimSpace(v)
	case []interface{}:
		var parts []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.TrimSpace(strings.Join(parts, ""))
	}
	return ""
}

func doReply(msg map[string]interface{}) string {
	runtimeMu.RLock()
	base := apiBase
	botSelfid := selfid
	auth := authorization
	runtimeMu.RUnlock()
	if base == "" {
		core.Logs.Warn("yyw: api_base 未配置，无法发送消息")
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
			sendMsg(base, auth, "SendImageMsg", botSelfid, toWxid, seg.value)
		case "audio", "voice", "record":
			sendMsg(base, auth, "SendVoiceMsg", botSelfid, toWxid, seg.value)
		case "video":
			sendMsg(base, auth, "SendVideoMsg", botSelfid, toWxid, seg.value)
		case "file":
			sendMsg(base, auth, "SendFileMsg", botSelfid, toWxid, seg.value)
		default:
			sendMsg(base, auth, "SendTextMsg", botSelfid, toWxid, seg.value)
		}
	}
	return ""
}

func sendMsg(apiBase, auth, event, selfid, toWxid, msg string) {
	body := map[string]string{
		"event":   event,
		"selfid":  selfid,
		"to_wxid": toWxid,
		"msg":     msg,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiBase+"/"+event, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		core.Logs.Warn("yyw发送%s失败: %v", event, err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	core.Logs.Info("yyw发送%s成功 to=%s", event, toWxid)
}

type yywSegment struct {
	kind  string
	value string
}

func splitSegments(text string) []yywSegment {
	matches := cqAnyPattern.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return []yywSegment{{kind: "text", value: text}}
	}
	var segs []yywSegment
	last := 0
	for _, m := range matches {
		if m[0] > last {
			if t := strings.TrimSpace(text[last:m[0]]); t != "" {
				segs = append(segs, yywSegment{kind: "text", value: t})
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
			kind = "audio"
		case "video":
			kind = "video"
		case "file":
			kind = "file"
		}
		if kind != "text" {
			segs = append(segs, yywSegment{kind: kind, value: url})
		}
		last = m[1]
	}
	if last < len(text) {
		if t := strings.TrimSpace(text[last:]); t != "" {
			segs = append(segs, yywSegment{kind: "text", value: t})
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

func stringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}
