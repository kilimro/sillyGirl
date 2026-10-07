<script setup lang="ts">
import { computed } from "vue";
import Alert from "ant-design-vue/es/alert";
import Avatar from "ant-design-vue/es/avatar";
import Button from "ant-design-vue/es/button";
import Card from "ant-design-vue/es/card";
import { CloudDownload, MessageSquare, Plug, Settings, Clock, ArrowRight } from "lucide-vue-next";
import Col from "ant-design-vue/es/col";
import Modal from "ant-design-vue/es/modal";
import Progress from "ant-design-vue/es/progress";
import Row from "ant-design-vue/es/row";
import Space from "ant-design-vue/es/space";
import Statistic from "ant-design-vue/es/statistic";
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
</script>

<template>
  <section v-if="page === 'welcome'" class="welcome-page">
    <Card :bordered="false" class="hero-card">
      <div class="hero-inner">
        <div class="hero-left">
          <svg width="56" height="56" viewBox="-125 -125 250 250" class="hero-avatar-svg">
            <defs>
              <mask id="bot-mask-welcome" maskUnits="userSpaceOnUse" x="-158" y="-158" width="316" height="316">
                <path d="M100.39 0.25C100.39 3.53 100.23 6.82 99.91 10.09C99.59 13.35 99.1 16.61 98.47 19.82C97.83 23.04 97.03 26.24 96.08 29.37C95.13 32.51 94.03 35.61 92.78 38.64C91.53 41.67 90.12 44.65 88.58 47.54C87.04 50.43 85.35 53.26 83.53 55.98C81.72 58.71 79.76 61.36 77.69 63.89C75.62 66.42 73.41 68.87 71.1 71.19C68.79 73.5 66.35 75.72 63.83 77.8C61.3 79.88 58.66 81.84 55.95 83.66C53.23 85.48 50.41 87.18 47.53 88.72C44.65 90.27 41.67 91.68 38.66 92.93C35.64 94.18 32.54 95.29 29.42 96.25C26.29 97.2 23.1 98 19.9 98.64C16.69 99.28 13.44 99.76 10.19 100.08C6.94 100.4 3.66 100.57 0.39 100.57C-2.88 100.57 -6.16 100.4 -9.41 100.08C-12.67 99.76 -15.92 99.28 -19.12 98.64C-22.33 98 -25.51 97.2 -28.64 96.25C-31.77 95.29 -34.86 94.18 -37.88 92.93C-40.9 91.68 -43.87 90.27 -46.75 88.72C-49.63 87.18 -52.45 85.48 -55.17 83.66C-57.89 81.84 -60.53 79.88 -63.05 77.8C-65.58 75.72 -68.01 73.5 -70.32 71.19C-72.63 68.87 -74.84 66.42 -76.91 63.89C-78.99 61.36 -80.94 58.71 -82.76 55.98C-84.57 53.26 -86.26 50.43 -87.8 47.54C-89.34 44.65 -90.75 41.67 -92 38.64C-93.25 35.61 -94.36 32.51 -95.31 29.37C-96.25 26.24 -97.05 23.04 -97.69 19.82C-98.33 16.61 -98.81 13.35 -99.13 10.09C-99.45 6.82 -99.61 3.53 -99.61 0.25C-99.61 -3.02 -99.45 -6.32 -99.13 -9.58C-98.81 -12.84 -98.33 -16.1 -97.69 -19.32C-97.05 -22.53 -96.25 -25.73 -95.31 -28.86C-94.36 -32 -93.25 -35.11 -92 -38.13C-90.75 -41.16 -89.34 -44.14 -87.8 -47.03C-86.26 -49.92 -84.57 -52.75 -82.76 -55.48C-80.94 -58.2 -78.99 -60.85 -76.91 -63.38C-74.84 -65.92 -72.63 -68.36 -70.32 -70.68C-68.01 -72.99 -65.58 -75.21 -63.05 -77.29C-60.53 -79.37 -57.89 -81.33 -55.17 -83.15C-52.45 -84.97 -49.63 -86.67 -46.75 -88.21C-43.87 -89.76 -40.9 -91.17 -37.88 -92.42C-34.86 -93.68 -31.77 -94.79 -28.64 -95.74C-25.51 -96.69 -22.33 -97.49 -19.12 -98.13C-15.92 -98.77 -12.67 -99.25 -9.41 -99.57C-6.16 -99.9 -2.88 -100.06 0.39 -100.06C3.66 -100.06 6.94 -99.9 10.19 -99.57C13.44 -99.25 16.69 -98.77 19.9 -98.13C23.1 -97.49 26.29 -96.69 29.42 -95.74C32.54 -94.79 35.64 -93.68 38.66 -92.42C41.67 -91.17 44.65 -89.76 47.53 -88.21C50.41 -86.67 53.23 -84.97 55.95 -83.15C58.66 -81.33 61.3 -79.37 63.83 -77.29C66.35 -75.21 68.79 -72.99 71.1 -70.68C73.41 -68.36 75.62 -65.92 77.69 -63.38C79.76 -60.85 81.72 -58.2 83.53 -55.48C85.35 -52.75 87.04 -49.92 88.58 -47.03C90.12 -44.14 91.53 -41.16 92.78 -38.13C94.03 -35.11 95.13 -32 96.08 -28.86C97.03 -25.73 97.83 -22.53 98.47 -19.32C99.1 -16.1 99.59 -12.84 99.91 -9.58C100.23 -6.32 100.39 -3.02 100.39 0.25Z" fill="#fff"/>
              </mask>
            </defs>
            <path d="M100.39 0.25C100.39 3.53 100.23 6.82 99.91 10.09C99.59 13.35 99.1 16.61 98.47 19.82C97.83 23.04 97.03 26.24 96.08 29.37C95.13 32.51 94.03 35.61 92.78 38.64C91.53 41.67 90.12 44.65 88.58 47.54C87.04 50.43 85.35 53.26 83.53 55.98C81.72 58.71 79.76 61.36 77.69 63.89C75.62 66.42 73.41 68.87 71.1 71.19C68.79 73.5 66.35 75.72 63.83 77.8C61.3 79.88 58.66 81.84 55.95 83.66C53.23 85.48 50.41 87.18 47.53 88.72C44.65 90.27 41.67 91.68 38.66 92.93C35.64 94.18 32.54 95.29 29.42 96.25C26.29 97.2 23.1 98 19.9 98.64C16.69 99.28 13.44 99.76 10.19 100.08C6.94 100.4 3.66 100.57 0.39 100.57C-2.88 100.57 -6.16 100.4 -9.41 100.08C-12.67 99.76 -15.92 99.28 -19.12 98.64C-22.33 98 -25.51 97.2 -28.64 96.25C-31.77 95.29 -34.86 94.18 -37.88 92.93C-40.9 91.68 -43.87 90.27 -46.75 88.72C-49.63 87.18 -52.45 85.48 -55.17 83.66C-57.89 81.84 -60.53 79.88 -63.05 77.8C-65.58 75.72 -68.01 73.5 -70.32 71.19C-72.63 68.87 -74.84 66.42 -76.91 63.89C-78.99 61.36 -80.94 58.71 -82.76 55.98C-84.57 53.26 -86.26 50.43 -87.8 47.54C-89.34 44.65 -90.75 41.67 -92 38.64C-93.25 35.61 -94.36 32.51 -95.31 29.37C-96.25 26.24 -97.05 23.04 -97.69 19.82C-98.33 16.61 -98.81 13.35 -99.13 10.09C-99.45 6.82 -99.61 3.53 -99.61 0.25C-99.61 -3.02 -99.45 -6.32 -99.13 -9.58C-98.81 -12.84 -98.33 -16.1 -97.69 -19.32C-97.05 -22.53 -96.25 -25.73 -95.31 -28.86C-94.36 -32 -93.25 -35.11 -92 -38.13C-90.75 -41.16 -89.34 -44.14 -87.8 -47.03C-86.26 -49.92 -84.57 -52.75 -82.76 -55.48C-80.94 -58.2 -78.99 -60.85 -76.91 -63.38C-74.84 -65.92 -72.63 -68.36 -70.32 -70.68C-68.01 -72.99 -65.58 -75.21 -63.05 -77.29C-60.53 -79.37 -57.89 -81.33 -55.17 -83.15C-52.45 -84.97 -49.63 -86.67 -46.75 -88.21C-43.87 -89.76 -40.9 -91.17 -37.88 -92.42C-34.86 -93.68 -31.77 -94.79 -28.64 -95.74C-25.51 -96.69 -22.33 -97.49 -19.12 -98.13C-15.92 -98.77 -12.67 -99.25 -9.41 -99.57C-6.16 -99.9 -2.88 -100.06 0.39 -100.06C3.66 -100.06 6.94 -99.9 10.19 -99.57C13.44 -99.25 16.69 -98.77 19.9 -98.13C23.1 -97.49 26.29 -96.69 29.42 -95.74C32.54 -94.79 35.64 -93.68 38.66 -92.42C41.67 -91.17 44.65 -89.76 47.53 -88.21C50.41 -86.67 53.23 -84.97 55.95 -83.15C58.66 -81.33 61.3 -79.37 63.83 -77.29C66.35 -75.21 68.79 -72.99 71.1 -70.68C73.41 -68.36 75.62 -65.92 77.69 -63.38C79.76 -60.85 81.72 -58.2 83.53 -55.48C85.35 -52.75 87.04 -49.92 88.58 -47.03C90.12 -44.14 91.53 -41.16 92.78 -38.13C94.03 -35.11 95.13 -32 96.08 -28.86C97.03 -25.73 97.83 -22.53 98.47 -19.32C99.1 -16.1 99.59 -12.84 99.91 -9.58C100.23 -6.32 100.39 -3.02 100.39 0.25Z" fill="#f9f9f9"/>
            <g mask="url(#bot-mask-welcome)">
              <rect x="-158" y="-158" width="316" height="316" fill="#0a0a0c"/>
            </g>
            <path d="M-9.3 -11.3A9.3 9.3 0 0 1 0 -20.6L0 -20.6A9.3 9.3 0 0 1 9.3 -11.3L9.3 11.3A9.3 9.3 0 0 1 0 20.6L0 20.6A9.3 9.3 0 0 1 -9.3 11.3Z" transform="matrix(0.88,-0.31,0.42,0.85,20.5,-41.91)" fill="#fff"/>
            <path d="M-9.3 -11.3A9.3 9.3 0 0 1 0 -20.6L0 -20.6A9.3 9.3 0 0 1 9.3 -11.3L9.3 11.3A9.3 9.3 0 0 1 0 20.6L0 20.6A9.3 9.3 0 0 1 -9.3 11.3Z" transform="matrix(0.66,-0.05,0.42,0.85,63.09,-51.83)" fill="#fff"/>
          </svg>
          <div>
            <Typography.Title :level="3" style="margin: 0 0 4px">{{ greeting() }}，{{ user?.name || "管理员" }}</Typography.Title>
            <Typography.Text type="secondary">MianPro 机器人管理平台 · 当前版本 {{ overviewVersion.local }}</Typography.Text>
          </div>
        </div>
        <div class="hero-right">
          <Statistic title="已接入 Bot" :value="enabledBots" />
          <Statistic title="在线 Bot" :value="onlineBots" />
          <Statistic title="脚本数量" :value="realScripts.length" />
          <Statistic title="总用户数" :value="overviewUserStats.total" />
        </div>
      </div>
    </Card>

    <Row :gutter="[16, 16]" style="margin-top: 16px">
      <Col :xs="24" :md="17">
        <Card title="已接入的 Bot" :bordered="false">
          <template #extra>
            <Button type="link" size="small" @click="page = 'bots'">查看全部 <ArrowRight :size="12" /></Button>
          </template>
          <div class="bot-grid">
            <div v-for="bot in botStatusRows" :key="bot.platform" class="bot-item" @click="page = 'bots'">
              <div class="bot-icon" :class="bot.online ? 'online' : ''">
                <span v-html="platformIcon(bot.platform)"></span>
              </div>
              <div class="bot-item-info">
                <Typography.Text strong>{{ bot.label }}</Typography.Text>
                <Tag v-if="bot.online" color="success" size="small" style="margin-top:4px">在线</Tag>
                <Tag v-else-if="bot.enabled" color="warning" size="small" style="margin-top:4px">未连接</Tag>
                <Tag v-else size="small" style="margin-top:4px">未启用</Tag>
              </div>
            </div>
          </div>
        </Card>

        <Card title="快捷操作" :bordered="false" style="margin-top: 16px">
          <Row :gutter="[12, 12]">
            <Col :xs="12" :sm="8" :md="6">
              <div class="quick-action" @click="page = 'bots'">
                <MessageSquare :size="20" class="quick-icon" />
                <Typography.Text>BOT 对接管理</Typography.Text>
              </div>
            </Col>
            <Col :xs="12" :sm="8" :md="6">
              <div class="quick-action" @click="page = 'plugins'">
                <Plug :size="20" class="quick-icon" />
                <Typography.Text>插件市场</Typography.Text>
              </div>
            </Col>
            <Col :xs="12" :sm="8" :md="6">
              <div class="quick-action" @click="page = 'tasks'">
                <Clock :size="20" class="quick-icon" />
                <Typography.Text>定时任务</Typography.Text>
              </div>
            </Col>
            <Col :xs="12" :sm="8" :md="6">
              <div class="quick-action" @click="page = 'settings'">
                <Settings :size="20" class="quick-icon" />
                <Typography.Text>基础设置</Typography.Text>
              </div>
            </Col>
          </Row>
        </Card>
      </Col>

      <Col :xs="24" :md="7" class="right-col">
        <Card title="系统信息" :bordered="false" class="side-card">
          <Space direction="vertical" style="width: 100%">
            <div class="info-row">
              <Typography.Text type="secondary">当前版本</Typography.Text>
              <Tag color="blue">{{ overviewVersion.local }}</Tag>
            </div>
            <div class="info-row">
              <Typography.Text type="secondary">最新版本</Typography.Text>
              <Tag color="green">{{ overviewVersion.remote }}</Tag>
            </div>
            <div class="info-row">
              <Typography.Text type="secondary">今日新增用户</Typography.Text>
              <Typography.Text strong>{{ overviewUserStats.today }}</Typography.Text>
            </div>
            <Button type="primary" block :loading="systemUpdate.running" @click="startOnlineUpdate">
              <template #icon><CloudDownload :size="15" /></template>
              在线更新
            </Button>
          </Space>
        </Card>

        <Card title="容器接入" :bordered="false" style="margin-top: 16px" class="side-card">
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
        <Progress :percent="systemUpdate.percent" :status="systemUpdate.status === 'error' ? 'exception' : systemUpdate.status === 'done' ? 'success' : 'active'" />
        <Alert :type="systemUpdate.status === 'error' ? 'error' : systemUpdate.status === 'done' ? 'success' : 'info'" :message="systemUpdate.message || '准备更新'" show-icon />
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
  background: #fff;
  border: none;
}
.hero-card :deep(.ant-card-body) {
  padding: 0;
}
.hero-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 32px;
  flex-wrap: wrap;
  gap: 16px;
}
.hero-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.hero-right {
  display: flex;
  gap: 40px;
  flex-wrap: wrap;
}
.hero-right :deep(.ant-statistic-title) {
  color: #999;
}
.hero-right :deep(.ant-statistic-content) {
  color: #333;
}
.bot-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}
.bot-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.bot-item:hover {
  border-color: #91caff;
  background: #e6f4ff;
}
.bot-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  flex-shrink: 0;
}
.bot-item-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.quick-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px 8px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.quick-action:hover {
  border-color: #91caff;
  background: #e6f4ff;
}
.quick-icon {
  color: #1677ff;
}
.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}
.right-col {
  display: flex;
  flex-direction: column;
}
.right-col :deep(.ant-card) {
  flex: 1;
}
</style>
