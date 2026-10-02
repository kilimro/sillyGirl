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
                <!-- 微信系 -->
                <svg v-if="['clawbot','flowbot','yyw'].includes(bot.platform)" viewBox="0 0 1024 1024" width="28" height="28">
                  <path d="M337.39 341.83c-17.76 0-35.51 11.84-35.51 29.6s17.76 29.59 35.51 29.59 29.6-11.84 29.6-29.59c0-18.5-11.84-29.6-29.6-29.6zM577.85 513.48c-11.84 0-22.94 12.58-22.94 23.68 0 12.58 11.84 23.68 22.94 23.68 17.76 0 29.6-11.84 29.6-23.68s-11.84-23.68-29.6-23.68zM501.64 401.02c17.76 0 29.6-12.58 29.6-29.6 0-17.76-11.84-29.6-29.6-29.6s-35.51 11.84-35.51 29.6 17.76 29.6 35.51 29.6zM706.59 513.48c-11.84 0-22.94 12.58-22.94 23.68 0 12.58 11.84 23.68 22.94 23.68 17.76 0 29.6-11.84 29.6-23.68s-11.84-23.68-29.6-23.68z" fill="#28C445"/>
                  <path d="M510.52 2.96C228.62 2.96 0 231.58 0 513.48s228.62 510.52 510.52 510.52 510.52-228.62 510.52-510.52S792.42 2.96 510.52 2.96zM413.6 644.44c-29.6 0-53.27-5.92-81.39-12.58l-81.39 41.43 22.94-71.77c-58.45-41.43-93.97-95.45-93.97-159.82 0-113.2 105.8-201.99 233.8-201.99 114.68 0 216.05 71.03 236.02 166.47-7.4-.74-14.8-1.48-22.2-1.48-110.98 1.48-198.29 85.09-198.29 188.67 0 17.02 2.96 33.29 7.4 49.57-7.4.74-15.54 1.48-22.94 1.48zM759.86 727.31l17.76 59.19-63.63-35.51c-22.94 5.92-46.61 11.84-70.29 11.84-111.72 0-199.77-76.95-199.77-172.39-.74-94.7 87.31-171.65 198.29-171.65 105.8 0 199.03 77.69 199.03 172.39 0 53.27-34.77 100.62-81.39 136.14z" fill="#28C445"/>
                </svg>
                <!-- QQ -->
                <svg v-else-if="['qq','qqguild'].includes(bot.platform)" viewBox="0 0 1024 1024" width="28" height="28">
                  <path d="M512 0C229.22 0 0 229.22 0 512s229.22 512 512 512 512-229.22 512-512S794.78 0 512 0zM801.26 668.86c-21.73 18.64-49.95-61.35-54.01-49.04-9.88 29.92-14.51 49.93-43.63 82.51-1.55 1.74 33.66 14.47 43.63 41.64 9.56 26.02 28.14 67.26-93.49 80.21-71.36 7.59-122.94-38.02-128.08-37.58-9.54.83-5.29 0-15.53 0-8.38 0-8.93.61-16.82 0-2.16-.17-25.88 37.58-131.96 37.58-82.22 0-103.52-51.75-86.98-80.21 16.54-28.46 44.13-36.75 40.24-41.26-19.15-22.19-32.35-45.91-40.24-67.36-1.95-5.35-3.58-10.54-4.87-15.53-2.99-11.45-25.88 67.2-50.46 49.04-24.58-18.17-22.38-64.42-6.47-108.68 16.04-44.62 56.47-87.59 56.93-97.07 1.61-35.29-3.49-41.15 0-50.42 7.76-20.76 17.2-12.8 17.2-23.57 0-135.74 100.86-245.76 225.28-245.76s225.28 110.04 225.28 245.76c0 5.2 13.52 0 19.99 23.57 1.33 4.87 2.24 23.67.66 50.42-.74 12.86 34.27 28.5 52.38 97.07C828.43 628.75 810.31 661.11 801.26 668.86z" fill="#68A5E1"/>
                </svg>
                <!-- Telegram -->
                <svg v-else-if="bot.platform === 'telegram'" viewBox="0 0 1024 1024" width="28" height="28">
                  <path d="M512 66.4C264.6 66.4 64 265.9 64 512c0 246.1 200.6 445.6 448 445.6S960 758.1 960 512c0-246.1-200.6-445.6-448-445.6z m232.1 243.3l-80.7 420.5c-5.6 29.9-22 37.1-44.5 23.3L496.1 653.1l-59 63.3c-6.8 7.6-12.5 13.8-24.6 13.8l8.3-138.9 227.8-227.2c10.1-10.2-2.1-15.1-15.4-6.2L351.9 554.6l-121.6-42.7c-26.1-8.2-26.4-28.2 5.9-42.7L709.9 267c21.7-10.8 42.4 5.9 34.2 42.7z" fill="#4DA8EF"/>
                </svg>
                <!-- 钉钉 -->
                <svg v-else-if="bot.platform === 'dingtalk'" viewBox="0 0 1024 1024" width="28" height="28">
                  <path d="M504.37 513.13m-450.82 0a450.82 450.82 0 1 0 901.63 0 450.82 450.82 0 1 0-901.63 0Z" fill="#59ADF8"/>
                  <path d="M665.65 688.28h89.14c-55.45 77.52-108.8 152.06-162.15 226.61l-3.84-1.38c11.62-48.95 23.3-97.89 35.43-148.94H554.5c7.94-36.15 15.36-70.09 23.5-107.21-19.81 5.99-37.12 9.42-52.94 16.28-34.15 14.9-64.87 8.55-92.62-13.47-18.33-14.59-36.3-30.41-51.2-48.33-15.56-18.74-10.7-29.39 13-33.43 50.23-8.45 100.61-15.92 151.09-25.29-10.55 0-21.09.1-31.69 0-43.42-.41-86.89-1.13-130.3-1.28-26.68-.1-45.47-13.93-61.54-33.23-20.79-24.93-34.51-53.4-42.24-84.79-4.92-19.86.51-25.34 21.04-20.58 51.87 12.08 103.68 24.58 155.55 36.71 26.52 6.2 53.15 11.98 80.44 15.41-31.49-10.39-63.28-19.92-94.36-31.44-47.36-17.56-94.11-36.66-141.31-54.68-14.34-5.48-23.76-15-29.8-28.77-18.43-42.04-32.2-85.3-32.56-131.74-.1-16.95 5.53-21.04 20.07-14.03 148.12 71.63 300.9 131.94 454.35 190.82 14.75 5.63 28.72 13.98 41.93 22.73 25.5 16.95 33.79 38.35 20.58 65.18-27.08 55.14-56.88 108.9-85.81 163.12-7.07 13.47-15.1 26.37-24.01 41.73z" fill="#FFFFFF"/>
                </svg>
                <!-- 飞书 -->
                <svg v-else-if="bot.platform === 'feishu'" viewBox="0 0 24 24" width="28" height="28">
                  <path d="M5.5 3h13l-3 6.5 3 6.5h-13l3-6.5z" fill="#3370FF"/>
                  <path d="M5.5 3L11 9.5 5.5 16h6l5.5-6.5L11 3z" fill="#00D6B9"/>
                </svg>
                <!-- Web/Pagermaid -->
                <svg v-else viewBox="0 0 1024 1024" width="28" height="28">
                  <path d="M170.67 0h682.67a170.67 170.67 0 0 1 170.67 170.67v682.67a170.67 170.67 0 0 1-170.67 170.67H170.67A170.67 170.67 0 0 1 0 853.33V170.67A170.67 170.67 0 0 1 162.13.21L170.67 0z" fill="#8F3FFF"/>
                  <path d="M514.52 170.67l301.18 170.67v341.33l-301.23 170.67L213.33 682.67V341.33l301.19-170.67zM276.74 377.26v269.65l237.78 134.7V512L276.74 377.26z" fill="#FFFFFF"/>
                </svg>
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
