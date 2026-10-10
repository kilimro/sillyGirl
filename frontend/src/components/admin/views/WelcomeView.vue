<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import Alert from "ant-design-vue/es/alert";
import Button from "ant-design-vue/es/button";
import Card from "ant-design-vue/es/card";
import { CloudDownload, Plug, ArrowRight } from "lucide-vue-next";
import Col from "ant-design-vue/es/col";
import Modal from "ant-design-vue/es/modal";
import Progress from "ant-design-vue/es/progress";
import Row from "ant-design-vue/es/row";
import Space from "ant-design-vue/es/space";
import Tag from "ant-design-vue/es/tag";
import Typography from "ant-design-vue/es/typography";
import { useAdminViewContext } from "../adminViewContext";
import { platformIcon } from "./platformIcons";

const {
  botStatusRows,
  overviewIntegrations,
  overviewUserStats,
  overviewVersion,
  page,
  realScripts,
  restartAfterUpdate,
  startOnlineUpdate,
  systemUpdate,
  user,
} = useAdminViewContext();

function greeting() {
  const h = new Date().getHours();
  if (h < 6) return "凌晨好";
  if (h < 9) return "早上好";
  if (h < 12) return "上午好";
  if (h < 14) return "中午好";
  if (h < 18) return "下午好";
  return "晚上好";
}

const onlineBots = computed(() => botStatusRows.value.filter((b: any) => b.online).length);
const enabledBots = computed(() => botStatusRows.value.filter((b: any) => b.enabled).length);

const quickEntries = [
  { key: "bots", label: "BOT 对接管理", desc: "统一接入各平台机器人", img: "/ill-bot.jpg" },
  { key: "plugins", label: "插件市场", desc: "按需扩展机器人能力", img: "/ill-book.jpg" },
  { key: "tasks", label: "定时任务", desc: "创建自动化定时任务", img: "/ill-clock.jpg" },
  { key: "settings", label: "基础设置", desc: "配置通知与存储偏好", img: "/ill-gear.jpg" },
];

// 数字滚动动画
const display = ref({ bots: 0, online: 0, scripts: 0, users: 0 });
onMounted(() => {
  const targets = {
    bots: enabledBots.value,
    online: onlineBots.value,
    scripts: realScripts.value.length,
    users: overviewUserStats.value.total,
  };
  const dur = 900;
  const start = performance.now();
  const tick = (now: number) => {
    const p = Math.min((now - start) / dur, 1);
    const ease = 1 - Math.pow(1 - p, 3);
    display.value = {
      bots: Math.round(targets.bots * ease),
      online: Math.round(targets.online * ease),
      scripts: Math.round(targets.scripts * ease),
      users: Math.round(targets.users * ease),
    };
    if (p < 1) requestAnimationFrame(tick);
  };
  requestAnimationFrame(tick);
});
</script>

<template>
  <section v-if="page === 'welcome'" class="welcome-page">
    <!-- hero 横幅 -->
    <Card :bordered="false" class="hero-card">
      <div class="hero-inner">
        <div class="hero-left">
          <img class="hero-blob" src="/blob.gif" alt="bot" />
        </div>
        <div class="hero-text">
          <div class="hero-eyebrow">WELCOME BACK</div>
          <Typography.Title :level="3" class="hero-title" style="margin: 0 0 4px">
            {{ greeting() }}，{{ user?.name || "管理员" }}
          </Typography.Title>
          <Typography.Text type="secondary" class="hero-sub">
            MianPro 机器人管理平台 · 当前版本 {{ overviewVersion.local }}
          </Typography.Text>
          <div class="hero-actions">
            <Button type="primary" :loading="systemUpdate.running" @click="startOnlineUpdate">
              <template #icon><CloudDownload :size="15" /></template>
              在线更新
            </Button>
            <Button class="hero-ghost" @click="page = 'plugins'">
              <template #icon><Plug :size="15" /></template>
              插件市场
            </Button>
          </div>
        </div>
      </div>
    </Card>

    <Row :gutter="[14, 12]" style="margin-top: 10px">
      <Col :xs="24" :md="17">
        <!-- 快捷入口 -->
        <div class="sec-head"><span class="sec-dot"></span>快捷入口</div>
        <div class="vision-grid">
          <div v-for="q in quickEntries" :key="q.key" class="vision-card" @click="page = q.key">
            <div class="vision-text">
              <div class="vision-title">{{ q.label }}</div>
              <div class="vision-desc">{{ q.desc }}</div>
            </div>
            <div class="vision-pic"><img :src="q.img" :alt="q.label" /></div>
          </div>
        </div>

        <!-- 已接入的 Bot -->
        <Card title="已接入的 Bot" :bordered="false" style="margin-top: 16px">
          <template #extra>
            <Button type="link" size="small" @click="page = 'bots'">
              查看全部 <ArrowRight :size="12" />
            </Button>
          </template>
          <div class="bot-list">
            <div v-for="bot in botStatusRows" :key="bot.platform" class="bot-row" @click="page = 'bots'">
              <div class="bot-row-icon"><span v-html="platformIcon(bot.platform)"></span></div>
              <span class="bot-row-name">{{ bot.label }}</span>
              <span class="bot-row-status" :class="bot.online ? 'on' : bot.enabled ? 'idle' : 'off'">
                <i class="dot"></i>{{ bot.online ? '在线' : bot.enabled ? '未连接' : '未启用' }}
              </span>
            </div>
          </div>
        </Card>
      </Col>

      <Col :xs="24" :md="7" class="right-col">
        <div class="sec-head"><span class="sec-dot"></span>Quick Stats</div>
        <div class="quick-grid">
          <div class="q"><div class="q-v num">{{ display.bots }}</div><div class="q-l">已接入 Bot</div></div>
          <div class="q"><div class="q-v num">{{ display.online }}</div><div class="q-l">在线 Bot</div></div>
          <div class="q"><div class="q-v num">{{ display.scripts }}</div><div class="q-l">脚本数量</div></div>
          <div class="q"><div class="q-v num">{{ display.users }}</div><div class="q-l">总用户数</div></div>
        </div>

        <Card title="系统信息" :bordered="false" style="margin-top: 14px" class="side-card">
          <div class="info-row">
            <Typography.Text type="secondary">当前版本</Typography.Text>
            <Tag class="ver-tag">{{ overviewVersion.local }}</Tag>
          </div>
          <div class="info-row">
            <Typography.Text type="secondary">最新版本</Typography.Text>
            <Tag class="ver-tag">{{ overviewVersion.remote }}</Tag>
          </div>
          <div class="side-divider"></div>
          <div v-for="item in overviewIntegrations" :key="item.key" class="info-row">
            <Typography.Text type="secondary">{{ item.label }}</Typography.Text>
            <Tag :color="item.online ? 'success' : 'default'" size="small">
              {{ item.online ? `${item.count} 个在线` : '未接入' }}
            </Tag>
          </div>
        </Card>
      </Col>
    </Row>

    <Modal
      v-model:open="systemUpdate.open"
      title="在线更新"
      :footer="null"
      :closable="!systemUpdate.running && !systemUpdate.restartChecking"
      :mask-closable="!systemUpdate.running && !systemUpdate.restartChecking"
    >
      <Space direction="vertical" style="width: 100%" size="middle">
        <Progress
          :percent="systemUpdate.percent"
          :status="systemUpdate.status === 'error' ? 'exception' : systemUpdate.status === 'done' ? 'success' : 'active'"
        />
        <Alert
          :type="systemUpdate.status === 'error' ? 'error' : systemUpdate.status === 'done' ? 'success' : 'info'"
          :message="systemUpdate.message || '准备更新'"
          show-icon
        />
        <div v-if="systemUpdate.result" class="update-result">
          <Typography.Text class="block">版本：{{ systemUpdate.result.before || "-" }} -> {{ systemUpdate.result.after || "-" }}</Typography.Text>
          <Typography.Text class="block">文件：{{ systemUpdate.result.asset || "-" }}</Typography.Text>
          <Typography.Text class="block muted">{{ systemUpdate.result.output || "" }}</Typography.Text>
        </div>
        <Space v-if="systemUpdate.status === 'done' && !systemUpdate.restartChecking" style="justify-content: flex-end; width: 100%">
          <Button @click="systemUpdate.open = false">关闭</Button>
          <Button v-if="systemUpdate.result" type="primary" :loading="systemUpdate.restarting" @click="restartAfterUpdate">立即重启</Button>
        </Space>
        <Space v-if="systemUpdate.status === 'error'" style="justify-content: flex-end; width: 100%">
          <Button @click="systemUpdate.open = false">关闭</Button>
        </Space>
      </Space>
    </Modal>
  </section>
</template>

<style scoped>
.hero-card {
  position: relative;
  overflow: hidden;
  background: #fff;
  border: none;
}
.hero-card::after {
  content: "";
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 46%;
  background: url("/hero.jpg") center / cover no-repeat;
  -webkit-mask: linear-gradient(90deg, transparent 0, #000 40%);
  mask: linear-gradient(90deg, transparent 0, #000 40%);
  pointer-events: none;
}
.hero-card :deep(.ant-card-body) {
  padding: 0;
}
.hero-inner {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 28px;
}
.hero-blob {
  width: 56px;
  height: 56px;
  flex-shrink: 0;
}
.hero-text {
  min-width: 0;
}
.hero-eyebrow {
  font-size: 12px;
  letter-spacing: 0.14em;
  color: #8a8f99;
  font-weight: 600;
  margin-bottom: 6px;
}
.hero-title {
  font-weight: 800 !important;
}
.hero-sub {
  font-size: 13.5px;
}
.hero-actions {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}
.hero-ghost {
  background: #f1f2f4;
  border: none;
}
.hero-ghost:hover {
  background: #e6e7ea;
}

.sec-head {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  font-size: 14.5px;
  color: #1a1c20;
  margin-bottom: 10px;
}
.sec-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid #1a1c20;
  display: inline-block;
}

/* 快捷入口：紧凑、插画 cover 不拉伸 */
.vision-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}
@media (max-width: 640px) {
  .vision-grid {
    grid-template-columns: 1fr;
  }
}
.vision-card {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 76px;
  padding: 0 14px;
  background: #fff;
  border: 1px solid #ececee;
  border-radius: 14px;
  box-shadow: 0 1px 2px rgba(16, 18, 22, 0.04), 0 6px 20px rgba(16, 18, 22, 0.05);
  cursor: pointer;
  transition: transform 0.18s, box-shadow 0.18s;
}
.vision-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(16, 18, 22, 0.06), 0 16px 40px rgba(16, 18, 22, 0.1);
}
.vision-text {
  flex: 1;
  min-width: 0;
}
.vision-pic {
  width: 60px;
  height: 60px;
  flex-shrink: 0;
  background: #f5f6f8;
  border-radius: 10px;
  overflow: hidden;
}
.vision-pic img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.vision-title {
  font-size: 13.5px;
  font-weight: 700;
}
.vision-desc {
  font-size: 12px;
  color: #8a8f99;
  margin-top: 3px;
}

/* 已接入 Bot：前沿列表（行式 + 状态圆点） */
.bot-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 2px 20px;
}
.bot-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.18s;
}
.bot-row:hover {
  background: #f4f5f7;
}
.bot-row-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  flex-shrink: 0;
}
.bot-row-name {
  flex: 1;
  font-size: 13.5px;
  font-weight: 600;
  color: #2a2c30;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.bot-row-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #9aa0a8;
  flex-shrink: 0;
}
.bot-row-status .dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #c9cdd4;
}
.bot-row-status.on {
  color: #16a34a;
}
.bot-row-status.on .dot {
  background: #16a34a;
}
.bot-row-status.idle {
  color: #9aa0a8;
}
.bot-row-status.idle .dot {
  background: #d4d7dc;
}
.bot-row-status.off {
  color: #b6bac0;
}
.bot-row-status.off .dot {
  background: #e2e5e9;
}

/* Quick Stats */
.quick-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.q {
  background: #fff;
  border: 1px solid #ececee;
  border-radius: 12px;
  height: 76px;
  padding: 0 14px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  box-shadow: 0 1px 2px rgba(16, 18, 22, 0.04), 0 6px 20px rgba(16, 18, 22, 0.05);
}
.q-v {
  font-size: 22px;
  font-weight: 800;
  color: #1a1c20;
}
.q-l {
  font-size: 12px;
  color: #8a8f99;
  margin-top: 2px;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 0;
}
.side-divider {
  height: 1px;
  background: #f0f0f2;
  margin: 10px 0;
}
.info-row > *:last-child {
  margin-left: auto;
}
.ver-tag {
  background: #f1f2f4 !important;
  border: none !important;
  color: #2a2c30 !important;
  font-weight: 600;
}
.right-col {
  display: flex;
  flex-direction: column;
}
.side-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.side-card :deep(.ant-card-body) {
  flex: 1;
}
.update-result .block {
  display: block;
}
.update-result .muted {
  color: #8a8f99;
}
</style>
