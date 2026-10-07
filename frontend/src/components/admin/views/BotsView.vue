<script setup lang="ts">
import { ref, watch, onUnmounted } from "vue";
import QRCode from "qrcode";
import { get, post } from "../../../api";
import Alert from "ant-design-vue/es/alert";
import {
  Bot,
  MessageSquare,
  QrCode,
  RefreshCw,
} from "lucide-vue-next";
import Button from "ant-design-vue/es/button";
import Empty from "ant-design-vue/es/empty";
import Form from "ant-design-vue/es/form";
import Input from "ant-design-vue/es/input";
import Modal from "ant-design-vue/es/modal";
import Segmented from "ant-design-vue/es/segmented";
import Space from "ant-design-vue/es/space";
import Spin from "ant-design-vue/es/spin";
import Switch from "ant-design-vue/es/switch";
import Typography from "ant-design-vue/es/typography";
import message from "ant-design-vue/es/message";
import { useAdminViewContext } from "../adminViewContext";
import BotCard from "./BotCard.vue";

const {
  botEnabled,
  botSettings,
  botSettingsModal,
  botStatusRows,
  cancelCurrentBotSettings,
  clawbotLogin,
  closeClawbotLogin,
  login,
  oneBotReceiveURL,
  openBotSettings,
  page,
  pagermaidBridgeURL,
  qqGuildWebhookURL,
  refreshBots,
  saveCurrentBotSettings,
  setBotEnabled,
  settings,
  startClawbotLogin,
  submitClawbotVerifyCode,
  toggleWebChat,
  webChat,
  webChatEndpointURL,
  feishuWebhookURL,
} = useAdminViewContext();

// 飞书一键接入
const feishuReg = ref<{ running: boolean; qr_url: string; expire_in: number; done: boolean; error: string; app_id: string }>({
  running: false, qr_url: "", expire_in: 0, done: false, error: "", app_id: "",
});
const feishuQrDataUrl = ref("");
watch(() => feishuReg.value.qr_url, async (url) => {
  if (url && !feishuReg.value.done) {
    try {
      feishuQrDataUrl.value = await QRCode.toDataURL(url, { width: 180, margin: 1 });
    } catch { feishuQrDataUrl.value = ""; }
  } else {
    feishuQrDataUrl.value = "";
  }
});
let feishuRegTimer: any = null;
const webhookUrl = (typeof window !== "undefined" ? window.location.origin : "") + "/api/yyw/webhook";
const geweWebhookUrl = (typeof window !== "undefined" ? window.location.origin : "") + "/api/gewe/webhook";


async function startFeishuRegister() {
  feishuReg.value = { running: true, qr_url: "", expire_in: 0, done: false, error: "", app_id: "" };
  try {
    const r = await post("/api/admin/feishu/register");
    console.log("feishu register POST:", r);
  } catch (e: any) {
    feishuReg.value.error = e.message || String(e);
    feishuReg.value.running = false;
    return;
  }
  // 轮询状态
  feishuRegTimer = setInterval(async () => {
    try {
      const r = await get<any>("/api/admin/feishu/register");
      const d = r.data || r;
      feishuReg.value = {
        running: d.running, qr_url: d.qr_url || "", expire_in: d.expire_in || 0,
        done: d.done, error: d.error || "", app_id: d.app_id || "",
      };
      if (d.done) {
        clearInterval(feishuRegTimer);
        feishuRegTimer = null;
        message.success("飞书接入成功！");
        refreshBots();
      }
    } catch {}
  }, 2000);
}
onUnmounted(() => { if (feishuRegTimer) clearInterval(feishuRegTimer); });
</script>

<template>
  <section v-if="page === 'bots'" class="panel">
    <div class="toolbar">
      <div class="toolbar-left">
        <Bot :size="16" />
        <Typography.Text strong>BOT 对接管理</Typography.Text>
      </div>
      <div class="toolbar-right">
        <Button @click="refreshBots"
          ><template #icon><RefreshCw :size="16" /></template>刷新状态</Button
        >
      </div>
    </div>
    <Spin :spinning="botSettings.loading">
      <div class="bot-card-grid">
        <BotCard
          v-for="record in botStatusRows"
          :key="record.platform"
          :record="record"
          :enabled="botEnabled(record)"
          @toggle="setBotEnabled(record, $event)"
          @open="openBotSettings(record)"
        />
      </div>
    </Spin>
  </section>

  <Modal
    v-model:open="botSettingsModal.open"
    :title="`${botSettingsModal.label} 设置`"
    width="680px"
    ok-text="保存"
    cancel-text="取消"
    :confirm-loading="botSettings.saving"
    @cancel="cancelCurrentBotSettings"
    @ok="saveCurrentBotSettings"
  >
    <Spin :spinning="botSettings.loading">
      <Form
        v-if="botSettingsModal.platform === 'clawbot'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用 ClawBot" html-for="bot-clawbot-enable">
          <Switch
            id="bot-clawbot-enable"
            v-model:checked="botSettings.form.clawbot_enable"
          />
        </Form.Item>
        <Form.Item
          label="Token"
          html-for="bot-clawbot-token"
          extra="ClawBot / OpenClaw 微信通道的 iLink bot token。保存后适配器会自动重启。"
        >
          <Space.Compact style="width: 100%">
            <Input.Password
              id="bot-clawbot-token"
              v-model:value="botSettings.form.clawbot_token"
              name="clawbot-token"
              placeholder="请输入 ClawBot Token"
            />
            <Button :loading="clawbotLogin.starting" @click="startClawbotLogin">
              <template #icon><QrCode :size="16" /></template>扫码获取
            </Button>
          </Space.Compact>
        </Form.Item>
        <Form.Item
          label="API 地址"
          html-for="bot-clawbot-api"
          extra="默认使用腾讯 iLink API；如果你有兼容反代地址可以填写在这里。"
        >
          <Input
            id="bot-clawbot-api"
            v-model:value="botSettings.form.clawbot_api_base"
            name="clawbot-api"
            placeholder="https://ilinkai.weixin.qq.com"
          />
        </Form.Item>
        <Form.Item label="ClawBot 调试日志" html-for="bot-clawbot-debug">
          <Switch
            id="bot-clawbot-debug"
            v-model:checked="botSettings.form.clawbot_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'qq'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用 QQ" html-for="bot-qq-enable">
          <Switch
            id="bot-qq-enable"
            v-model:checked="botSettings.form.qq_enable"
          />
        </Form.Item>
        <Form.Item label="反向 WebSocket 地址" html-for="bot-qq-receive">
          <Input id="bot-qq-receive" :value="oneBotReceiveURL" readonly>
            <template #suffix
              ><Typography.Text class="muted"
                >NapCat 填这个 URL</Typography.Text
              ></template
            >
          </Input>
        </Form.Item>
        <Form.Item
          label="连接密钥"
          html-for="bot-qq-token"
          extra="需要和 NapCat / OneBot 客户端配置里的 accessToken 保持一致；公网部署建议必须填写。"
        >
          <Input.Password
            id="bot-qq-token"
            v-model:value="botSettings.form.qq_token"
            name="qq-token"
            placeholder="请输入 QQ 连接密钥"
          />
        </Form.Item>
        <Form.Item label="QQ 调试日志" html-for="bot-qq-debug">
          <Switch
            id="bot-qq-debug"
            v-model:checked="botSettings.form.qq_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'telegram'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用 Telegram" html-for="bot-telegram-enable">
          <Switch
            id="bot-telegram-enable"
            v-model:checked="botSettings.form.telegram_enable"
          />
        </Form.Item>
        <Form.Item
          label="Token"
          html-for="bot-telegram-token"
          extra="BotFather 提供的 Bot Token，保存后 Telegram 适配器会自动重启。"
        >
          <Input.Password
            id="bot-telegram-token"
            v-model:value="botSettings.form.telegram_token"
            name="telegram-token"
            placeholder="123456:ABC-DEF..."
          />
        </Form.Item>
        <Form.Item
          label="代理 API"
          html-for="bot-telegram-api"
          extra="默认使用 https://api.telegram.org；网络不通时填写兼容反代地址。"
        >
          <Input
            id="bot-telegram-api"
            v-model:value="botSettings.form.telegram_api_base"
            name="telegram-api"
            placeholder="https://api.telegram.org"
          />
        </Form.Item>
        <Form.Item label="Telegram 调试日志" html-for="bot-telegram-debug">
          <Switch
            id="bot-telegram-debug"
            v-model:checked="botSettings.form.telegram_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'feishu'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用飞书机器人" html-for="bot-feishu-enable">
          <Switch
            id="bot-feishu-enable"
            v-model:checked="botSettings.form.feishu_enable"
          />
        </Form.Item>
        <Form.Item label="快捷接入" extra="无需手动创建应用，点按钮后扫码确认即可自动配置。">
          <Space>
            <Button type="primary" :loading="feishuReg.running && !feishuReg.done" @click="startFeishuRegister">
              一键接入飞书
            </Button>
          </Space>
        </Form.Item>
        <div v-if="feishuReg.qr_url && !feishuReg.done" style="margin-bottom: 12px; padding: 12px; background: #e6f4ff; border: 1px solid #91caff; border-radius: 6px;">
          <div style="margin-bottom: 8px; color: #0958d9;">请在 10 分钟内完成创建：手机飞书扫码 或 打开链接（两种方式均可）：</div>
          <div style="display: flex; align-items: flex-start; gap: 16px;">
            <img v-if="feishuQrDataUrl" :src="feishuQrDataUrl" alt="飞书扫码接入二维码" style="width: 180px; height: 180px; border: 1px solid #d9d9d9; border-radius: 4px;" />
            <a :href="feishuReg.qr_url" target="_blank" style="word-break: break-all; color: #1677ff;">{{ feishuReg.qr_url }}</a>
          </div>
        </div>
        <div v-if="feishuReg.done" style="margin-bottom: 12px; padding: 12px; background: #f6ffed; border: 1px solid #b7eb8f; border-radius: 6px; color: #389e0d;">
          接入成功！App ID: {{ feishuReg.app_id }}，配置已自动保存。
        </div>
        <div v-if="feishuReg.error" style="margin-bottom: 12px; padding: 12px; background: #fff2f0; border: 1px solid #ffccc7; border-radius: 6px; color: #cf1322;">
          {{ feishuReg.error }}
        </div>
        <div style="margin-bottom: 16px; padding: 12px; background: #fafafa; border: 1px solid #e8e8e8; border-radius: 6px; font-size: 13px; line-height: 1.8;">
          <div style="font-weight: 600; margin-bottom: 8px; color: #333;">手动配置步骤（无需一键接入时）：</div>
          <div style="color: #666;">
            <div>1. 打开 <a href="https://open.feishu.cn/app" target="_blank" style="color: #1677ff;">飞书开放平台</a>，创建「企业自建应用」</div>
            <div>2. 在「凭证与基础信息」复制 <b>App ID</b> 和 <b>App Secret</b> 填到下面</div>
            <div>3. 在「权限管理」开通以下权限：</div>
            <div style="padding-left: 20px; color: #888;">
              · im:message:send_as_bot（以机器人身份发消息）<br>
              · im:message.p2p_msg:readonly（接收单聊消息）<br>
              · im:resource（上传图片/文件）<br>
              · im:message.reactions:write_only（表情回应）
            </div>
            <div>4. 在「事件与回调」→「事件配置」→ 选择<b>使用长连接接收事件</b>，添加事件：<b>接收消息 im.message.receive_v1</b></div>
            <div>5. 在「应用能力」→「机器人」启用机器人能力</div>
            <div>6. 点「保存」即可，无需公网地址</div>
          </div>
        </div>
        <Form.Item
          label="App ID"
          html-for="bot-feishu-app-id"
          extra="飞书开放平台 → 企业自建应用 → 凭证与基础信息 里的 App ID。"
        >
          <Input
            id="bot-feishu-app-id"
            v-model:value="botSettings.form.feishu_app_id"
            name="feishu-app-id"
            placeholder="cli_xxxxxxxxxxxxxxxx"
          />
        </Form.Item>
        <Form.Item label="App Secret" html-for="bot-feishu-app-secret">
          <Input.Password
            id="bot-feishu-app-secret"
            v-model:value="botSettings.form.feishu_app_secret"
            name="feishu-app-secret"
            placeholder="应用 App Secret"
          />
        </Form.Item>
        <Form.Item
          label="Verification Token"
          html-for="bot-feishu-token"
          extra="长连接模式下无需填写，保留兼容。"
        >
          <Input.Password
            id="bot-feishu-token"
            v-model:value="botSettings.form.feishu_verification_token"
            name="feishu-token"
            placeholder="（可选，长连接不需要）"
          />
        </Form.Item>
        <Form.Item
          label="接入方式"
          html-for="bot-feishu-mode"
          extra="长连接模式：飞书开放平台 → 事件与回调 → 回调配置 → 选择「使用长连接接收事件」，无需公网地址。"
        >
          <Input value="长连接（WebSocket）" disabled />
        </Form.Item>
        <Form.Item label="飞书调试日志" html-for="bot-feishu-debug">
          <Switch
            id="bot-feishu-debug"
            v-model:checked="botSettings.form.feishu_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'yyw'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <div style="margin-bottom: 16px; padding: 12px; background: #fafafa; border: 1px solid #e8e8e8; border-radius: 6px; font-size: 13px; line-height: 1.8;">
          <div style="font-weight: 600; margin-bottom: 8px; color: #333;">YYW 微信插件配置步骤：</div>
          <div style="color: #666;">
            <div>1. 在 YYW 插件里，把「远程处理接口」填成：</div>
            <div style="padding: 8px 12px; background: #f0f0f0; border-radius: 4px; font-family: monospace; margin: 8px 0; color: #333;">{{ webhookUrl }}</div>
            <div>2. 把插件管理页显示的「内网 IP + 端口」填到下面的 API 地址</div>
            <div>3. 把插件管理页的「鉴权密钥」填到下面（如果插件开了通信鉴权）</div>
            <div>4. 保存即可</div>
          </div>
        </div>
        <Form.Item label="启用 YYW 微信" html-for="bot-yyw-enable">
          <Switch id="bot-yyw-enable" v-model:checked="botSettings.form.yyw_enable" />
        </Form.Item>
        <Form.Item label="插件 API 地址" html-for="bot-yyw-api-base" extra="插件管理页显示的内网 IP + 端口，如 http://192.168.1.100:2025">
          <Input id="bot-yyw-api-base" v-model:value="botSettings.form.yyw_api_base" placeholder="http://IP:端口" />
        </Form.Item>
        <Form.Item label="鉴权密钥" html-for="bot-yyw-auth" extra="插件管理页的「鉴权密钥」，开启通信鉴权时必填">
          <Input.Password id="bot-yyw-auth" v-model:value="botSettings.form.yyw_authorization" placeholder="插件的鉴权密钥" />
        </Form.Item>
        <Form.Item label="机器人 wxid（可选）" html-for="bot-yyw-selfid" extra="留空则自动从收到的消息中识别">
          <Input id="bot-yyw-selfid" v-model:value="botSettings.form.yyw_selfid" placeholder="wxid_xxx（可选）" />
        </Form.Item>
      </Form>

      <Form v-else-if="botSettingsModal.platform === 'gewe'" layout="vertical" class="bot-settings-modal-form">
        <Form.Item label="启用 GeWe 微信" html-for="bot-gewe-enable">
          <Switch id="bot-gewe-enable" v-model:checked="botSettings.form.gewe_enable" />
        </Form.Item>
        <Form.Item label="Token" html-for="bot-gewe-token" extra="GeWe 后台「Token 中心」生成的 Token，即 X-GEWE-TOKEN。">
          <Input.Password id="bot-gewe-token" v-model:value="botSettings.form.gewe_token" name="gewe-token" placeholder="请输入 GeWe Token" />
        </Form.Item>
        <Form.Item label="设备 ID（AppID）" html-for="bot-gewe-app-id" extra="GeWe 后台扫码登录微信后复制到的设备 ID。">
          <Input id="bot-gewe-app-id" v-model:value="botSettings.form.gewe_app_id" name="gewe-app-id" placeholder="wx_xxxxxxxxxxxxxxxx" />
        </Form.Item>
        <Form.Item label="API 地址" html-for="bot-gewe-api" extra="默认官方 https://api.geweapi.com；使用自建节点时可填节点地址。">
          <Input id="bot-gewe-api" v-model:value="botSettings.form.gewe_api_base" name="gewe-api" placeholder="https://api.geweapi.com" />
        </Form.Item>
        <Form.Item label="回调地址" html-for="bot-gewe-webhook" extra="默认当前访问地址+路径；部署公网后改为公网域名回调。保存后自动调用 setCallback，失败可后台手动填。">
          <Input id="bot-gewe-webhook" v-model:value="botSettings.form.gewe_webhook" name="gewe-webhook" :placeholder="geweWebhookUrl" />
        </Form.Item>
        <Form.Item label="机器人 wxid（可选）" html-for="bot-gewe-self-wxid" extra="留空则自动从收到的消息中识别">
          <Input id="bot-gewe-self-wxid" v-model:value="botSettings.form.gewe_self_wxid" placeholder="wxid_xxx（可选）" />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'dingtalk'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用钉钉" html-for="bot-dingtalk-enable">
          <Switch
            id="bot-dingtalk-enable"
            v-model:checked="botSettings.form.dingtalk_enable"
          />
        </Form.Item>
        <Form.Item
          label="Client ID"
          html-for="bot-dingtalk-client-id"
          extra="钉钉开放平台应用的 Client ID（原 AppKey）。适配器使用 Stream 模式，不需要公网回调地址。"
        >
          <Input
            id="bot-dingtalk-client-id"
            v-model:value="botSettings.form.dingtalk_client_id"
            name="dingtalk-client-id"
            placeholder="dingxxxxxxxx"
          />
        </Form.Item>
        <Form.Item label="Client Secret" html-for="bot-dingtalk-secret">
          <Input.Password
            id="bot-dingtalk-secret"
            v-model:value="botSettings.form.dingtalk_client_secret"
            name="dingtalk-client-secret"
            placeholder="请输入 Client Secret"
          />
        </Form.Item>
        <Form.Item label="钉钉调试日志" html-for="bot-dingtalk-debug">
          <Switch
            id="bot-dingtalk-debug"
            v-model:checked="botSettings.form.dingtalk_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'qqguild'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用 QQ 官方频道" html-for="bot-qqguild-enable">
          <Switch
            id="bot-qqguild-enable"
            v-model:checked="botSettings.form.qqguild_enable"
          />
        </Form.Item>
        <fieldset class="bot-mode-fieldset">
          <legend id="bot-qqguild-mode-label">接入模式</legend>
          <Segmented
            v-model:value="botSettings.form.qqguild_mode"
            block
            role="radiogroup"
            aria-labelledby="bot-qqguild-mode-label"
            :options="[
              { label: 'Webhook', value: 'webhook' },
              { label: 'WebSocket', value: 'websocket' },
            ]"
          />
        </fieldset>
        <Form.Item
          v-if="botSettings.form.qqguild_mode === 'webhook'"
          label="Webhook 回调地址"
          html-for="bot-qqguild-webhook"
          extra="把该地址填入 QQ 开放平台机器人事件回调；HTTPS 由反向代理提供。"
        >
          <Input id="bot-qqguild-webhook" :value="qqGuildWebhookURL" readonly />
        </Form.Item>
        <Alert
          v-else
          type="info"
          show-icon
          message="WebSocket 主动连接 QQ Gateway，不需要公网回调地址。"
          style="margin-bottom: 18px"
        />
        <Form.Item label="AppID" html-for="bot-qqguild-app-id">
          <Input
            id="bot-qqguild-app-id"
            v-model:value="botSettings.form.qqguild_app_id"
            name="qqguild-app-id"
            placeholder="请输入机器人 AppID"
          />
        </Form.Item>
        <Form.Item label="AppSecret" html-for="bot-qqguild-secret">
          <Input.Password
            id="bot-qqguild-secret"
            v-model:value="botSettings.form.qqguild_app_secret"
            name="qqguild-app-secret"
            placeholder="请输入机器人 AppSecret"
          />
        </Form.Item>
        <Form.Item label="沙箱环境" html-for="bot-qqguild-sandbox">
          <Switch
            id="bot-qqguild-sandbox"
            v-model:checked="botSettings.form.qqguild_sandbox"
          />
        </Form.Item>
        <Form.Item label="QQ 频道调试日志" html-for="bot-qqguild-debug">
          <Switch
            id="bot-qqguild-debug"
            v-model:checked="botSettings.form.qqguild_debug"
          />
        </Form.Item>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'web'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item
          label="运行状态"
          html-for="bot-web-status"
          extra="Web Bot 是内置适配器，随 MianPro 自动启动。"
        >
          <Switch id="bot-web-status" :checked="true" disabled />
        </Form.Item>
        <Form.Item
          label="允许匿名聊天"
          html-for="bot-web-public"
          extra="关闭时仅已登录的后台管理员可以发送 Web Bot 消息。"
        >
          <Switch
            id="bot-web-public"
            v-model:checked="botSettings.form.web_chat_public"
          />
        </Form.Item>
        <Form.Item label="聊天接口" html-for="bot-web-endpoint">
          <Input id="bot-web-endpoint" :value="webChatEndpointURL" readonly />
        </Form.Item>
        <Button type="primary" @click="toggleWebChat">
          <template #icon><MessageSquare :size="16" /></template>
          {{ webChat.open ? "关闭聊天窗口" : "打开聊天窗口" }}
        </Button>
      </Form>

      <Form
        v-else-if="botSettingsModal.platform === 'flowbot'"
        layout="vertical"
        class="bot-settings-modal-form"
      >
        <Form.Item label="启用 FlowBot" html-for="bot-flowbot-enable">
          <Switch
            id="bot-flowbot-enable"
            v-model:checked="botSettings.form.flowbot_enable"
          />
        </Form.Item>
        <Form.Item
          label="API Key"
          html-for="bot-flowbot-token"
          extra="FlowBot 插件 API 的鉴权密钥（api_key）。保存后适配器会自动重启并重新连接。"
        >
          <Input.Password
            id="bot-flowbot-token"
            v-model:value="botSettings.form.flowbot_token"
            name="flowbot-token"
            placeholder="请输入 FlowBot API Key"
          />
        </Form.Item>
        <Form.Item
          label="主机地址"
          html-for="bot-flowbot-host"
          extra="FlowBot 插件进程监听的地址，默认 127.0.0.1（本机）。"
        >
          <Input
            id="bot-flowbot-host"
            v-model:value="botSettings.form.flowbot_host"
            name="flowbot-host"
            placeholder="127.0.0.1"
          />
        </Form.Item>
        <Form.Item
          label="端口"
          html-for="bot-flowbot-port"
          extra="FlowBot 插件监听的端口，默认 7400。"
        >
          <Input
            id="bot-flowbot-port"
            v-model:value="botSettings.form.flowbot_port"
            name="flowbot-port"
            placeholder="7400"
          />
        </Form.Item>
        <Form.Item label="启用 TLS" html-for="bot-flowbot-tls">
          <Switch
            id="bot-flowbot-tls"
            v-model:checked="botSettings.form.flowbot_tls"
          />
        </Form.Item>
        <Form.Item label="FlowBot 调试日志" html-for="bot-flowbot-debug">
          <Switch
            id="bot-flowbot-debug"
            v-model:checked="botSettings.form.flowbot_debug"
          />
        </Form.Item>
      </Form>

      <Form v-else layout="vertical" class="bot-settings-modal-form">
        <Form.Item label="启用 Pagermaid" html-for="bot-pagermaid-enable">
          <Switch
            id="bot-pagermaid-enable"
            v-model:checked="botSettings.form.pagermaid_enable"
          />
        </Form.Item>
        <Form.Item label="连接密钥" html-for="bot-pagermaid-token">
          <Input.Password
            id="bot-pagermaid-token"
            v-model:value="botSettings.form.pagermaid_token"
            name="pagermaid-token"
            placeholder="可选，建议填写"
          />
        </Form.Item>
        <Form.Item label="Pagermaid 调试日志" html-for="bot-pagermaid-debug">
          <Switch
            id="bot-pagermaid-debug"
            v-model:checked="botSettings.form.pagermaid_debug"
          />
        </Form.Item>
        <div class="bot-settings-readonly">
          <Typography.Text class="muted block">桥接脚本</Typography.Text>
          <Typography.Text class="mono block"
            >adapters/pagermaid/sillyplus.py</Typography.Text
          >
        </div>
        <div class="bot-settings-readonly">
          <Typography.Text class="muted block">WebSocket 地址</Typography.Text>
          <Typography.Text class="mono block">{{
            pagermaidBridgeURL
          }}</Typography.Text>
        </div>
      </Form>
    </Spin>
  </Modal>

  <Modal
    v-model:open="clawbotLogin.open"
    title="扫码获取 ClawBot Token"
    :footer="null"
    @cancel="closeClawbotLogin"
  >
    <div class="clawbot-login-modal">
      <div class="clawbot-qr-frame">
        <Spin :spinning="clawbotLogin.starting">
          <img
            v-if="clawbotLogin.qrcodeImg"
            :src="clawbotLogin.qrcodeImg"
            alt="ClawBot 登录二维码"
          />
          <Empty
            v-else
            :description="clawbotLogin.message || '等待生成二维码'"
          />
        </Spin>
      </div>
      <Typography.Text>{{
        clawbotLogin.message || "请使用微信扫码"
      }}</Typography.Text>
      <div v-if="clawbotLogin.needVerify" class="clawbot-verify-row">
        <Input
          v-model:value="clawbotLogin.verifyCode"
          placeholder="输入手机微信显示的数字"
          @pressEnter="submitClawbotVerifyCode"
        />
        <Button
          type="primary"
          :loading="clawbotLogin.polling"
          @click="submitClawbotVerifyCode"
          >确认</Button
        >
      </div>
      <Space>
        <Button :loading="clawbotLogin.starting" @click="startClawbotLogin">
          <template #icon><RefreshCw :size="16" /></template>
          重新生成
        </Button>
        <Button @click="closeClawbotLogin">关闭</Button>
      </Space>
    </div>
  </Modal>
</template>
