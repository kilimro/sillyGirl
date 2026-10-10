<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { User, ArrowRight, BookOpen, MessageSquare, Boxes, Zap, Radio, Shield } from 'lucide-vue-next';
import AppBrand from './components/common/AppBrand.vue';
import { platformIcon } from './components/admin/views/platformIcons';

const platforms = ['wechat', 'qq', 'feishu', 'telegram', 'dingtalk', 'web'];

const iframeSrc = ref('');
const loading = ref(true);
const showDefault = ref(false);

onMounted(async () => {
  try {
    const res = await fetch('/api/public/home-config');
    const data = await res.json();
    const url = data?.data?.iframe_url;
    if (url && String(url).trim()) {
      iframeSrc.value = String(url).trim();
    } else {
      showDefault.value = true;
    }
  } catch {
    showDefault.value = true;
  } finally {
    loading.value = false;
  }
});

function openUserCenter() { window.location.href = '/user'; }
function openDashboard() { window.location.href = '/admin'; }
</script>

<template>
  <div class="home-page">
    <!-- iframe 模式 -->
    <iframe v-if="!loading && iframeSrc" :src="iframeSrc" class="main-iframe" />

    <!-- 默认首页 -->
    <div v-if="!loading && showDefault" class="default-home">
      <!-- 背景渐变光斑 -->
      <div class="blob blob1"></div>
      <div class="blob blob2"></div>

      <div class="hero-wrap">
        <div class="hero-left">
          <AppBrand class="home-brand" />
          <div class="badge">
            <span class="dot"></span>
            多平台机器人管理
          </div>
          <h1 class="hero-title">
            一个平台，管理你的<br />
            <span class="gradient-text">所有 AI 机器人</span>
          </h1>
          <p class="hero-desc">
            MianPro 是一个开源的多平台机器人管理面板，支持 QQ、微信、飞书、Telegram 等平台统一接入，插件化扩展，开箱即用。
          </p>
          <div class="hero-actions">
            <button class="btn-primary" @click="openDashboard">
              前往控制台 <ArrowRight :size="16" />
            </button>
            <button class="btn-ghost" @click="openUserCenter">
              <User :size="16" /> 用户中心
            </button>
          </div>
          <div class="hero-features">
            <div class="feat"><MessageSquare :size="14" /> 多平台同步</div>
            <div class="feat"><Boxes :size="14" /> 插件生态</div>
            <div class="feat"><Zap :size="14" /> Docker 部署</div>
            <div class="feat"><Shield :size="14" /> 开源免费</div>
          </div>
        </div>

        <div class="hero-right">
          <div class="glass-card">
            <div class="logo">
              <span class="circle c1"></span>
              <span class="circle c2"></span>
              <span class="circle c3"></span>
              <span class="circle c4"></span>
              <span class="circle c5">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="4" y="8" width="16" height="10" rx="2"></rect>
                  <path d="M12 8V4"></path>
                  <circle cx="12" cy="3" r="1"></circle>
                  <path d="M9 13h.01M15 13h.01"></path>
                </svg>
              </span>
            </div>
            <div class="glass"></div>
            <div class="g-content">
              <span class="g-title">MianPro 控制台</span>
              <span class="g-text">微信 · QQ · 飞书 · Telegram · 钉钉，一个面板统一接入</span>
            </div>
            <div class="g-platforms">
              <span class="p-icon" v-for="p in platforms" :key="p" v-html="platformIcon(p)"></span>
            </div>
            <div class="g-bottom">
              <div class="g-feats">
                <span class="g-feat"><MessageSquare :size="13" /> 多平台</span>
                <span class="g-feat"><Boxes :size="13" /> 插件生态</span>
                <span class="g-feat"><Shield :size="13" /> 开源免费</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="stats-row">
        <div class="stat-item"><div class="num">9+</div><div class="label">平台接入</div></div>
        <div class="stat-item"><div class="num">50+</div><div class="label">插件扩展</div></div>
        <div class="stat-item"><div class="num">100k+</div><div class="label">消息处理</div></div>
        <div class="stat-item"><div class="num">10+</div><div class="label">部署方式</div></div>
      </div>
    </div>

    <!-- 悬浮用户中心按钮 -->
    <button class="float-btn" @click="openUserCenter">
      <User :size="18" />
      <span>用户中心</span>
    </button>
  </div>
</template>

<style scoped>
.home-page { width: 100vw; height: 100vh; height: 100dvh; overflow: hidden; position: relative; }
.main-iframe { width: 100%; height: 100%; border: 0; }

.default-home {
  width: 100%; height: 100%; overflow-y: auto;
  background: #fafbff;
  position: relative;
}

.blob { position: absolute; border-radius: 50%; filter: blur(90px); opacity: 0.45; }
.blob1 { width: 520px; height: 520px; background: #99f6e4; top: -120px; right: 8%; }
.blob2 { width: 420px; height: 420px; background: #cbd5e1; bottom: 8%; left: -120px; }

.hero-wrap {
  max-width: 1100px; margin: 0 auto; padding: 80px 32px 40px;
  display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 60px; align-items: center;
  position: relative; z-index: 1;
}

.badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 6px 14px; border-radius: 20px;
  background: #f1f5f4; color: #134e4a; font-size: 13px; font-weight: 600;
  margin-bottom: 24px;
}
.home-brand { margin-bottom: 16px; }
.badge .dot { width: 6px; height: 6px; border-radius: 50%; background: #0f766e; }

.hero-title { font-size: 48px; font-weight: 800; line-height: 1.15; color: #111827; margin: 0 0 20px; }
.gradient-text {
  background: linear-gradient(135deg, #0f766e, #134e4a);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}
.hero-desc { font-size: 16px; color: #6b7280; line-height: 1.7; margin: 0 0 32px; max-width: 480px; }

.hero-actions { display: flex; gap: 12px; margin-bottom: 28px; }
.btn-primary {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 12px 26px; border: 0; border-radius: 10px;
  background: #0a0a0c; color: #fff; font-size: 15px; font-weight: 600; cursor: pointer;
  box-shadow: 0 6px 20px rgba(10,10,12,0.25);
}
.btn-primary:hover { background: #232327; }
.btn-ghost {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 12px 24px; border: 1px solid #e5e7eb; border-radius: 10px;
  background: #fff; color: #374151; font-size: 15px; cursor: pointer;
}
.btn-ghost:hover { background: #f9fafb; }

.hero-features { display: flex; gap: 16px; flex-wrap: wrap; }
.feat { display: flex; align-items: center; gap: 4px; font-size: 13px; color: #6b7280; }
.feat svg { color: #0f766e; }

/* 右侧玻璃卡 */
.hero-right { position: relative; perspective: 1000px; }
.glass-card {
  position: relative;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(18px);
  border: 1px solid rgba(255, 255, 255, 0.9);
  border-radius: 22px;
  padding: 26px 26px 22px;
  box-shadow: 0 24px 60px rgba(15, 118, 110, 0.12), 0 2px 8px rgba(0,0,0,0.04);
  overflow: hidden;
  transform-style: preserve-3d;
  transition: all 0.5s ease-in-out;
}
.glass-card:hover {
  transform: rotate3d(1, 1, 0, 10deg) translateY(-6px);
  box-shadow: rgba(15, 118, 110, 0.25) 30px 40px 40px -30px, rgba(15,118,110,0.15) 0 25px 35px 0;
}
.glass {
  position: absolute;
  top: -40%; right: -20%;
  width: 260px; height: 260px;
  background: radial-gradient(circle, rgba(153,246,228,0.55), transparent 70%);
  filter: blur(30px);
  pointer-events: none;
  transform: translate3d(0, 0, 25px);
}
.logo { display: flex; align-items: center; gap: 7px; margin-bottom: 22px; position: relative; z-index: 1; }
.circle {
  width: 16px; height: 16px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  transform: translate3d(0, 0, 26px);
  transition: all 0.4s ease;
}
.c1 { background: #d1fae5; }
.c2 { background: #a7f3d0; }
.c3 { background: #6ee7b7; }
.c4 { background: #34d399; }
.glass-card:hover .c1 { transform: translate3d(0, -4px, 30px); }
.glass-card:hover .c2 { transform: translate3d(0, -6px, 34px); }
.glass-card:hover .c3 { transform: translate3d(0, -8px, 38px); }
.glass-card:hover .c4 { transform: translate3d(0, -10px, 42px); }
.c5 {
  width: 34px; height: 34px; border-radius: 10px;
  background: #0a0a0c; color: #fff;
  transform: translate3d(0, 0, 40px);
}
.c5 svg { width: 18px; height: 18px; }
.g-content { position: relative; z-index: 1; margin-bottom: 20px; transform: translate3d(0,0,26px); }
.g-title { display: block; font-size: 18px; font-weight: 700; color: #111827; margin-bottom: 6px; }
.g-text { display: block; font-size: 13.5px; color: #6b7280; line-height: 1.6; }
.g-platforms {
  position: relative; z-index: 1;
  display: flex; gap: 12px; align-items: center;
  margin-bottom: 20px; transform: translate3d(0,0,26px);
}
.p-icon {
  width: 40px; height: 40px; border-radius: 12px;
  background: #fff; display: grid; place-content: center;
  box-shadow: 0 3px 10px rgba(0,0,0,0.08);
}
.p-icon :deep(svg) { width: 22px; height: 22px; }
.g-bottom {
  position: relative; z-index: 1;
  padding-top: 16px; border-top: 1px solid rgba(0,0,0,0.06);
  transform: translate3d(0,0,26px);
}
.g-feats { display: flex; gap: 18px; flex-wrap: wrap; }
.g-feat { display: inline-flex; align-items: center; gap: 5px; font-size: 12.5px; color: #4b5563; }
.g-feat svg { color: #0f766e; }

/* 底部统计 */
.stats-row {
  max-width: 900px; margin: 60px auto 0; padding: 40px 32px;
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px;
  border-top: 1px solid #e5e7eb; position: relative; z-index: 1;
}
.stat-item { text-align: center; }
.stat-item .num { font-size: 28px; font-weight: 700; color: #111827; }
.stat-item .label { font-size: 13px; color: #6b7280; margin-top: 4px; }

.float-btn {
  position: fixed; right: 24px; bottom: 24px; z-index: 999;
  display: flex; align-items: center; gap: 6px;
  padding: 10px 18px; border: 0; border-radius: 24px;
  background: #0a0a0c; color: #fff; font-size: 13px; font-weight: 600;
  cursor: pointer; box-shadow: 0 6px 18px rgba(10,10,12,0.3);
}

@media (max-width: 900px) {
  .hero-wrap { grid-template-columns: 1fr; padding-top: 40px; }
  .hero-title { font-size: 32px; }
  .stats-row { grid-template-columns: repeat(2, 1fr); }
}
</style>
