package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kilimro/sillyGirl/core"
	"github.com/kilimro/sillyGirl/core/storage"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/larksuite/oapi-sdk-go/v3/scene/registration"
)

var cqAnyPattern = regexp.MustCompile(`\[CQ:(\w+),([^\]]+)\]`)

var feishu = core.MakeBucket("feishu")

var runtime = struct {
	sync.Mutex
	cancel context.CancelFunc
	inst   *feishuBot
}{}

var restartTimer *time.Timer

type feishuBot struct {
	appID     string
	appSecret string
	apiCli    *lark.Client
	adapter   *core.Factory
	botID     string
	pendingRx sync.Map
}

func restart() {
	// debounce：快速连续改配置时只在最后一次触发后 1 秒才真正重启
	if restartTimer != nil {
		restartTimer.Stop()
	}
	restartTimer = time.AfterFunc(1*time.Second, doRestart)
}

func doRestart() {
	runtime.Lock()
	if runtime.cancel != nil {
		runtime.cancel()
		runtime.cancel = nil
		if runtime.inst != nil && runtime.inst.adapter != nil {
			runtime.inst.adapter.Destroy()
		}
		runtime.inst = nil
	}
	runtime.Unlock()

	if !enabled() {
		core.Logs.Info("feishu机器人未启动：feishu.enable=false")
		return
	}
	appID := strings.TrimSpace(feishu.GetString("app_id"))
	appSecret := strings.TrimSpace(feishu.GetString("app_secret"))
	core.Logs.Info("feishu restart: app_id_len=%d secret_len=%d", len(appID), len(appSecret))
	if appID == "" || appSecret == "" {
		core.Logs.Info("feishu机器人未启动：未配置 feishu.app_id / feishu.app_secret")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	b := &feishuBot{
		appID:     appID,
		appSecret: appSecret,
	}
	if err := b.start(ctx); err != nil {
		core.Logs.Warn("feishu机器人启动失败：%v", err)
		cancel()
		return
	}
	runtime.Lock()
	runtime.cancel = cancel
	runtime.inst = b
	runtime.Unlock()
}

func init() {
	core.FeishuRegisterHandler = StartRegister
	core.FeishuRegisterStatusHandler = GetRegisterStatus
	storage.Watch(feishu, "app_id", func(old, new, key string) *storage.Final {
		go restart()
		return nil
	})
	storage.Watch(feishu, "app_secret", func(old, new, key string) *storage.Final {
		go restart()
		return nil
	})
	storage.Watch(feishu, "enable", func(old, new, key string) *storage.Final {
		go restart()
		return nil
	})
	go func() {
		time.Sleep(2 * time.Second)
		restart()
	}()
}

// 一键注册状态
var registerState = struct {
	sync.Mutex
	running  bool
	qrURL    string
	qrExpire int
	done     bool
	err      string
	appID    string
	appSecret string
}{}

func StartRegister(ctx *gin.Context) {
	core.Logs.Info("feishu一键接入：收到请求, token=%s", ctx.GetHeader("token"))
	registerState.Lock()
	if registerState.running {
		url := registerState.qrURL
		expire := registerState.qrExpire
		registerState.Unlock()
		core.ApiOK(ctx, gin.H{"ok": true, "qr_url": url, "expire_in": expire, "running": true})
		return
	}
	registerState.running = true
	registerState.done = false
	registerState.err = ""
	registerState.appID = ""
	registerState.appSecret = ""
	registerState.Unlock()

	go func() {
		core.Logs.Info("feishu一键接入：启动 registration.RegisterApp")
		falseVal := false
		result, err := registration.RegisterApp(context.Background(), &registration.Options{
			Source:     "sdk",
			AppPreset:  &registration.AppPreset{Name: "MianPro Bot", Desc: "MianPro 机器人"},
			CreateOnly: false,
			Addons: &registration.AppAddons{
				Preset: &falseVal,
				Scopes: registration.AppAddonsScopes{
					Tenant: []string{
						"im:message:send_as_bot",
						"im:message.p2p_msg:readonly",
						"im:message:readonly",
						"im:resource",
						"im:message.reactions:write_only",
						"im:chat:read",
						"application:bot.basic_info:read",
					},
				},
				Events: registration.AppAddonsEvents{
					Items: registration.AppAddonsEventItems{
						Tenant: []string{"im.message.receive_v1"},
					},
				},
			},
			OnQRCode: func(info *registration.QRCodeInfo) {
				core.Logs.Info("feishu一键接入：收到QRCode URL=%s expire=%d", info.URL, info.ExpireIn)
				registerState.Lock()
				registerState.qrURL = info.URL
				registerState.qrExpire = info.ExpireIn
				registerState.Unlock()
			},
			OnStatusChange: func(info *registration.StatusChangeInfo) {
				core.Logs.Info("feishu一键接入：状态变更=%s", info.Status)
			},
		})
		core.Logs.Info("feishu一键接入：RegisterApp返回 err=%v", err)
		registerState.Lock()
		defer registerState.Unlock()
		registerState.running = false
		if err != nil {
			registerState.err = err.Error()
			core.Logs.Warn("feishu一键接入失败: %v", err)
			return
		}
		registerState.done = true
		registerState.appID = result.ClientID
		registerState.appSecret = result.ClientSecret
		// 自动写入配置
		feishu.Set("app_id", result.ClientID)
		feishu.Set("app_secret", result.ClientSecret)
		feishu.Set("enable", "true")
		core.Logs.Info("feishu一键接入成功 app_id=%s", result.ClientID)
	}()

	core.ApiOK(ctx, gin.H{"ok": true, "started": true})
}

func GetRegisterStatus(ctx *gin.Context) {
	registerState.Lock()
	defer registerState.Unlock()
	core.ApiOK(ctx, gin.H{
		"running":    registerState.running,
		"qr_url":     registerState.qrURL,
		"expire_in":  registerState.qrExpire,
		"done":       registerState.done,
		"error":      registerState.err,
		"app_id":     registerState.appID,
	})
}

func (b *feishuBot) start(ctx context.Context) error {
	b.apiCli = lark.NewClient(b.appID, b.appSecret)

	b.botID = b.appID

	b.adapter = &core.Factory{}
	b.adapter.Init("feishu", b.botID, nil)
	b.adapter.SetReplyHandler(func(msg map[string]interface{}) string {
		return b.reply(ctx, msg)
	})

	handler := larkevent.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			b.handleMessage(event)
			return nil
		}).
		OnP2MessageReactionCreatedV1(func(ctx context.Context, event *larkim.P2MessageReactionCreatedV1) error {
			// 忽略用户表情反应事件
			return nil
		}).
		OnP2MessageReactionDeletedV1(func(ctx context.Context, event *larkim.P2MessageReactionDeletedV1) error {
			// 忽略用户取消表情事件
			return nil
		})

	wsCli := larkws.NewClient(b.appID, b.appSecret,
		larkws.WithEventHandler(handler),
	)

	core.Logs.Info("feishu机器人(%s)长连接已启动", b.botID)
	go func() {
		for {
			if err := wsCli.Start(ctx); err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				core.Logs.Warn("feishu长连接断开，5秒后重连：%v", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}()
	return nil
}

func (b *feishuBot) handleMessage(event *larkim.P2MessageReceiveV1) {
	if event == nil || event.Event == nil || event.Event.Message == nil {
		return
	}
	msg := event.Event.Message
	if msg.MessageType == nil || *msg.MessageType != "text" {
		return
	}
	if msg.Content == nil {
		return
	}
	var content struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*msg.Content), &content); err != nil {
		return
	}
	text := strings.TrimSpace(content.Text)
	if text == "" {
		return
	}

	userID := ""
	if event.Event.Sender != nil && event.Event.Sender.SenderId != nil {
		if event.Event.Sender.SenderId.OpenId != nil {
			userID = *event.Event.Sender.SenderId.OpenId
		}
		if userID == "" && event.Event.Sender.SenderId.UserId != nil {
			userID = *event.Event.Sender.SenderId.UserId
		}
	}
	chatID := ""
	if msg.ChatId != nil {
		chatID = *msg.ChatId
	}
	messageID := ""
	if msg.MessageId != nil {
		messageID = *msg.MessageId
	}
	chatType := ""
	if msg.ChatType != nil {
		chatType = *msg.ChatType
	}

	params := map[string]interface{}{
		core.USER_ID:     userID,
		core.CHAT_ID:     chatID,
		core.CONETNT:     text,
		core.MESSAGE_ID:  messageID,
		"user_name":      userID,
		"chat_name":      chatID,
		"feishu_chat_id": chatID,
		"feishu_open_id": userID,
		"feishu_chat_type": chatType,
	}
	core.Logs.Info("feishu收到消息 from=%s chat=%s text=%s", userID, chatID, text)
	// 给用户消息加个"稍等"表情，回复完删掉
	if messageID != "" {
		go b.addWaitReaction(messageID)
	}
	if feishu.GetBool("debug") {
		core.Logs.Debug("feishu 消息详情：%s", larkcore.Prettify(params))
	}
	b.adapter.Receive(params)
}

// addWaitReaction 给消息加"稍等"表情，reaction_id 存起来等回复后删
func (b *feishuBot) addWaitReaction(messageID string) {
	req := larkim.NewCreateMessageReactionReqBuilder().
		MessageId(messageID).
		Body(&larkim.CreateMessageReactionReqBody{
			ReactionType: &larkim.Emoji{EmojiType: larkcore.StringPtr("StatusInFlight")},
		}).
		Build()
	resp, err := b.apiCli.Im.MessageReaction.Create(context.Background(), req)
	if err != nil {
		core.Logs.Warn("feishu加稍等表情失败 msg=%s err=%v", messageID, err)
		return
	}
	if resp != nil && resp.Data != nil && resp.Data.ReactionId != nil {
		b.pendingRx.Store(messageID, *resp.Data.ReactionId)
	}
}

// removeWaitReaction 删掉之前加的"稍等"表情
func (b *feishuBot) removeWaitReaction(messageID string) {
	v, ok := b.pendingRx.LoadAndDelete(messageID)
	if !ok {
		return
	}
	reactionID := v.(string)
	req := larkim.NewDeleteMessageReactionReqBuilder().
		MessageId(messageID).
		ReactionId(reactionID).
		Build()
	_, err := b.apiCli.Im.MessageReaction.Delete(context.Background(), req)
	if err != nil {
		core.Logs.Warn("feishu删稍等表情失败 msg=%s err=%v", messageID, err)
	}
}

func (b *feishuBot) reply(ctx context.Context, msg map[string]interface{}) string {
	chatID := stringValue(msg[core.CHAT_ID])
	if chatID == "" {
		chatID = stringValue(msg["feishu_chat_id"])
	}
	if chatID == "" {
		return ""
	}
	text := strings.TrimSpace(stringValue(msg[core.CONETNT]))
	if text == "" {
		return ""
	}

	// 解析 CQ 码，把 [CQ:image,file=url] 拆成文本段和图片段
	segments := splitFeishuReplySegments(text)
	lastMsgID := ""
	for _, seg := range segments {
		var (
			mid string
			err error
		)
		switch seg.kind {
		case "image":
			mid, err = b.sendImage(ctx, chatID, seg.value)
		case "audio":
			mid, err = b.sendFileMsg(ctx, chatID, seg.value, "audio", "opus", "voice.opus")
		case "media":
			mid, err = b.sendFileMsg(ctx, chatID, seg.value, "media", "mp4", "video.mp4")
		case "file":
			mid, err = b.sendFileMsg(ctx, chatID, seg.value, "file", "stream", "file")
		default:
			mid, err = b.sendText(ctx, chatID, seg.value)
		}
		if err != nil {
			core.Logs.Warn("feishu发送消息失败 chat=%s kind=%s err=%v", chatID, seg.kind, err)
			continue
		}
		if mid != "" {
			lastMsgID = mid
		}
	}
	// 删掉之前给用户消息加的"稍等"表情
	if mid := stringValue(msg[core.MESSAGE_ID]); mid != "" {
		go b.removeWaitReaction(mid)
	}
	return lastMsgID
}

type feishuSegment struct {
	kind  string // "text" | "image" | "file" | "audio" | "media"
	value string
}

func splitFeishuReplySegments(text string) []feishuSegment {
	matches := cqAnyPattern.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return []feishuSegment{{kind: "text", value: text}}
	}
	var segs []feishuSegment
	last := 0
	for _, m := range matches {
		if m[0] > last {
			if t := strings.TrimSpace(text[last:m[0]]); t != "" {
				segs = append(segs, feishuSegment{kind: "text", value: t})
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
			kind = "media"
		case "file":
			kind = "file"
		}
		if kind != "text" {
			segs = append(segs, feishuSegment{kind: kind, value: url})
		}
		last = m[1]
	}
	if last < len(text) {
		if t := strings.TrimSpace(text[last:]); t != "" {
			segs = append(segs, feishuSegment{kind: "text", value: t})
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

func (b *feishuBot) sendText(ctx context.Context, chatID, text string) (string, error) {
	content, _ := json.Marshal(map[string]string{"text": text})
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(&larkim.CreateMessageReqBody{
			ReceiveId: larkcore.StringPtr(chatID),
			MsgType:   larkcore.StringPtr("text"),
			Content:   larkcore.StringPtr(string(content)),
		}).
		Build()
	resp, err := b.apiCli.Im.Message.Create(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("code=%d msg=%s", resp.Code, resp.Msg)
	}
	if resp.Data != nil && resp.Data.MessageId != nil {
		core.Logs.Info("feishu发送文本成功 chat=%s msg_id=%s", chatID, *resp.Data.MessageId)
		return *resp.Data.MessageId, nil
	}
	return "", nil
}

func (b *feishuBot) sendImage(ctx context.Context, chatID, imageURL string) (string, error) {
	// 下载图片
	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("下载图片HTTP %d", resp.StatusCode)
	}
	// 上传到飞书
	uploadReq := larkim.NewCreateImageReqBuilder().
		Body(&larkim.CreateImageReqBody{
			ImageType: larkcore.StringPtr("message"),
			Image:     resp.Body,
		}).
		Build()
	uploadResp, err := b.apiCli.Im.Image.Create(ctx, uploadReq)
	if err != nil {
		return "", fmt.Errorf("上传图片失败: %w", err)
	}
	if uploadResp.Code != 0 {
		return "", fmt.Errorf("上传图片code=%d msg=%s", uploadResp.Code, uploadResp.Msg)
	}
	if uploadResp.Data == nil || uploadResp.Data.ImageKey == nil {
		return "", fmt.Errorf("上传图片无image_key")
	}
	imageKey := *uploadResp.Data.ImageKey
	// 发图片消息
	content, _ := json.Marshal(map[string]string{"image_key": imageKey})
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(&larkim.CreateMessageReqBody{
			ReceiveId: larkcore.StringPtr(chatID),
			MsgType:   larkcore.StringPtr("image"),
			Content:   larkcore.StringPtr(string(content)),
		}).
		Build()
	sendResp, err := b.apiCli.Im.Message.Create(ctx, req)
	if err != nil {
		return "", err
	}
	if sendResp.Code != 0 {
		return "", fmt.Errorf("发送图片code=%d msg=%s", sendResp.Code, sendResp.Msg)
	}
	if sendResp.Data != nil && sendResp.Data.MessageId != nil {
		core.Logs.Info("feishu发送图片成功 chat=%s msg_id=%s", chatID, *sendResp.Data.MessageId)
		return *sendResp.Data.MessageId, nil
	}
	return "", nil
}

// sendFileMsg 下载URL → 上传飞书得 file_key → 以 file/audio/media 类型发送
// msgType: file | audio | media
// fileType: 飞书文件类型枚举 opus/mp4/pdf/stream/mp3 等
func (b *feishuBot) sendFileMsg(ctx context.Context, chatID, fileURL, msgType, fileType, fileName string) (string, error) {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("下载文件HTTP %d", resp.StatusCode)
	}
	// 上传到飞书
	uploadReq := larkim.NewCreateFileReqBuilder().
		Body(larkim.NewCreateFileReqBodyBuilder().
			FileType(fileType).
			FileName(fileName).
			File(resp.Body).
			Build()).
		Build()
	uploadResp, err := b.apiCli.Im.File.Create(ctx, uploadReq)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}
	if uploadResp.Code != 0 {
		return "", fmt.Errorf("上传文件code=%d msg=%s", uploadResp.Code, uploadResp.Msg)
	}
	if uploadResp.Data == nil || uploadResp.Data.FileKey == nil {
		return "", fmt.Errorf("上传文件无file_key")
	}
	fileKey := *uploadResp.Data.FileKey
	// 发文件消息
	content, _ := json.Marshal(map[string]string{"file_key": fileKey})
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(&larkim.CreateMessageReqBody{
			ReceiveId: larkcore.StringPtr(chatID),
			MsgType:   larkcore.StringPtr(msgType),
			Content:   larkcore.StringPtr(string(content)),
		}).
		Build()
	sendResp, err := b.apiCli.Im.Message.Create(ctx, req)
	if err != nil {
		return "", err
	}
	if sendResp.Code != 0 {
		return "", fmt.Errorf("发送文件code=%d msg=%s", sendResp.Code, sendResp.Msg)
	}
	if sendResp.Data != nil && sendResp.Data.MessageId != nil {
		core.Logs.Info("feishu发送%s成功 chat=%s msg_id=%s", msgType, chatID, *sendResp.Data.MessageId)
		return *sendResp.Data.MessageId, nil
	}
	return "", nil
}

func enabled() bool {
	v := strings.TrimSpace(strings.ToLower(feishu.GetString("enable")))
	return v != "false" && v != "0" && v != "no" && v != "off"
}

func stringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}
