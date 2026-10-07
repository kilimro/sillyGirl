<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { User, ArrowRight, BookOpen, MessageSquare, Boxes, Zap, Radio, Shield } from 'lucide-vue-next';
import AppBrand from './components/common/AppBrand.vue';

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
          <div class="preview-card">
            <div class="preview-body">
              <div class="pv-header">
                <div class="pv-title">MianPro 控制台</div>
                <div class="pv-sub">你的机器人，在这里统一管理</div>
              </div>
              <div class="stat-grid">
                <div class="stat">
                  <div class="stat-num">9</div>
                  <div class="stat-label">已接入平台</div>
                </div>
                <div class="stat">
                  <div class="stat-num">12</div>
                  <div class="stat-label">可用插件</div>
                </div>
                <div class="stat">
                  <div class="stat-num">2.4k</div>
                  <div class="stat-label">服务用户</div>
                </div>
              </div>
              <div class="pv-feature">
                <div class="pv-feat-item">
                  <div class="pv-feat-dot" style="background:#22c55e"></div>
                  <div>
                    <div class="pv-feat-title">多平台统一接入</div>
                    <div class="pv-feat-desc">QQ / 微信 / 飞书 / Telegram 一个面板搞定</div>
                  </div>
                </div>
                <div class="pv-feat-item">
                  <div class="pv-feat-dot" style="background:#3b82f6"></div>
                  <div>
                    <div class="pv-feat-title">插件化扩展</div>
                    <div class="pv-feat-desc">按需开启，灵活组合你的 AI 助手</div>
                  </div>
                </div>
                <div class="pv-feat-item">
                  <div class="pv-feat-dot" style="background:#8b5cf6"></div>
                  <div>
                    <div class="pv-feat-title">开源免费</div>
                    <div class="pv-feat-desc">Docker 一键部署，数据自己掌控</div>
                  </div>
                </div>
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

.blob { position: absolute; border-radius: 50%; filter: blur(80px); opacity: 0.5; }
.blob1 { width: 500px; height: 500px; background: #93c5fd; top: -100px; right: 10%; }
.blob2 { width: 400px; height: 400px; background: #c4b5fd; bottom: 10%; left: -100px; }

.hero-wrap {
  max-width: 1100px; margin: 0 auto; padding: 80px 32px 40px;
  display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 60px; align-items: center;
  position: relative; z-index: 1;
}

.badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 6px 14px; border-radius: 20px;
  background: #eff6ff; color: #2563eb; font-size: 13px; font-weight: 500;
  margin-bottom: 24px;
}
.home-brand { margin-bottom: 16px; }
.badge .dot { width: 6px; height: 6px; border-radius: 50%; background: #2563eb; }

.hero-title { font-size: 48px; font-weight: 800; line-height: 1.15; color: #111827; margin: 0 0 20px; }
.gradient-text {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}
.hero-desc { font-size: 16px; color: #6b7280; line-height: 1.7; margin: 0 0 32px; max-width: 480px; }

.hero-actions { display: flex; gap: 12px; margin-bottom: 28px; }
.btn-primary {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 12px 24px; border: 0; border-radius: 10px;
  background: #2563eb; color: #fff; font-size: 15px; font-weight: 500; cursor: pointer;
  box-shadow: 0 4px 14px rgba(37,99,235,0.3);
}
.btn-primary:hover { background: #1d4ed8; }
.btn-ghost {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 12px 24px; border: 1px solid #e5e7eb; border-radius: 10px;
  background: #fff; color: #374151; font-size: 15px; cursor: pointer;
}
.btn-ghost:hover { background: #f9fafb; }

.hero-features { display: flex; gap: 16px; flex-wrap: wrap; }
.feat { display: flex; align-items: center; gap: 4px; font-size: 13px; color: #6b7280; }
.feat svg { color: #2563eb; }

/* 右侧预览卡片 */
.hero-right { position: relative; }
.preview-card {
  background: #fff; border-radius: 16px; overflow: hidden;
  box-shadow: 0 20px 60px rgba(0,0,0,0.1);
  border: 1px solid #f0f0f0;
}
.preview-body { padding: 24px; }
.pv-header { margin-bottom: 20px; }
.pv-title { font-size: 16px; font-weight: 700; color: #111827; }
.pv-sub { font-size: 13px; color: #6b7280; margin-top: 4px; }
.stat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-bottom: 20px; }
.stat { text-align: center; padding: 14px 8px; background: #f8fafc; border-radius: 10px; }
.stat-num { font-size: 22px; font-weight: 700; color: #111827; }
.stat-label { font-size: 12px; color: #6b7280; margin-top: 2px; }
.pv-feature { display: grid; gap: 14px; }
.pv-feat-item { display: flex; align-items: flex-start; gap: 10px; }
.pv-feat-dot { width: 8px; height: 8px; border-radius: 50%; margin-top: 6px; flex-shrink: 0; }
.pv-feat-title { font-size: 13px; font-weight: 600; color: #1f2937; }
.pv-feat-desc { font-size: 12px; color: #6b7280; margin-top: 2px; }

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
  background: #2563eb; color: #fff; font-size: 13px; font-weight: 500;
  cursor: pointer; box-shadow: 0 4px 16px rgba(37,99,235,0.35);
}

@media (max-width: 900px) {
  .hero-wrap { grid-template-columns: 1fr; padding-top: 40px; }
  .hero-title { font-size: 32px; }
  .stats-row { grid-template-columns: repeat(2, 1fr); }
}
</style>
