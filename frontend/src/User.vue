<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import AntApp from 'ant-design-vue/es/app';
import Avatar from 'ant-design-vue/es/avatar';
import Button from 'ant-design-vue/es/button';
import Card from 'ant-design-vue/es/card';
import Col from 'ant-design-vue/es/col';
import ConfigProvider from 'ant-design-vue/es/config-provider';
import Form from 'ant-design-vue/es/form';
import Input from 'ant-design-vue/es/input';
import Row from 'ant-design-vue/es/row';
import Space from 'ant-design-vue/es/space';
import Tag from 'ant-design-vue/es/tag';
import Typography from 'ant-design-vue/es/typography';
import message from 'ant-design-vue/es/message';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import { Bell, Link2, LogOut, Plug, QrCode, User as UserIcon } from 'lucide-vue-next';

type ApiEnvelope<T> = { status: boolean; message: string; data: T };
type PublicUser = { id: string; username: string; nickname: string; created_at: number };
type Bindings = { qq?: string; telegram?: string; updated_at?: number };
type UserAnnouncement = { enabled: boolean; content: string; format: string };
type OpenPlugin = { id: string; title: string; desc?: string; icon?: string; version?: string };

const userAuthTokenKey = 'sillygirl_user_jwt';
const loading = ref(true);
const currentUser = ref<PublicUser | null>(null);
const bindings = ref<Bindings>({});
const announcement = ref<UserAnnouncement>({ enabled: false, content: '', format: 'text' });
const openPlugins = ref<OpenPlugin[]>([]);

const authMode = ref<'login' | 'register'>('login');
const loginForm = reactive({ username: '', password: '' });
const registerForm = reactive({ username: '', nickname: '', password: '' });
const authLoading = ref(false);

const qqValue = ref('');
const tgValue = ref('');
const savingQQ = ref(false);
const savingTG = ref(false);

const userInitial = computed(() => {
  const name = currentUser.value?.nickname || currentUser.value?.username || 'U';
  return String(name).slice(0, 1).toUpperCase();
});

const createdDate = computed(() => {
  if (!currentUser.value?.created_at) return '-';
  return new Date(currentUser.value.created_at * 1000).toLocaleDateString('zh-CN');
});

const announceText = computed(() => {
  const c = announcement.value.content;
  return typeof c === 'string' ? c : '';
});

const announceHtml = computed(() => {
  const text = announceText.value;
  if (announcement.value.format === 'html') return text;
  if (announcement.value.format === 'markdown' || announcement.value.format === 'md') {
    return text
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
      .replace(/^### (.+)$/gm, '<h4>$1</h4>')
      .replace(/^## (.+)$/gm, '<h3>$1</h3>')
      .replace(/^# (.+)$/gm, '<h2>$1</h2>')
      .replace(/\*\*(.+?)\*\*/g, '<b>$1</b>')
      .replace(/\*(.+?)\*/g, '<i>$1</i>')
      .replace(/\[(.+?)\]\((.+?)\)/g, '<a href="$2" target="_blank">$1</a>')
      .replace(/`(.+?)`/g, '<code>$1</code>')
      .replace(/\n/g, '<br>');
  }
  return '';
});

async function requestJSON<T>(url: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers);
  if (!headers.has('Content-Type') && options.body) headers.set('Content-Type', 'application/json');
  const token = localStorage.getItem(userAuthTokenKey)?.trim();
  if (token && !headers.has('token')) headers.set('token', token);
  const res = await fetch(url, { ...options, headers });
  const payload = await res.json().catch(() => ({ status: false, message: '错误', data: null }));
  if (!res.ok || payload.status === false) throw new Error(payload.message || '请求失败');
  return payload.data;
}

async function login() {
  authLoading.value = true;
  try {
    const data = await requestJSON<{ token: string; user: PublicUser }>('/api/user/sessions', {
      method: 'POST', body: JSON.stringify({ username: loginForm.username.trim(), password: loginForm.password }),
    });
    localStorage.setItem(userAuthTokenKey, data.token.trim());
    currentUser.value = data.user;
    message.success('登录成功');
    await loadProfile();
  } catch (e) { message.error(e instanceof Error ? e.message : '登录失败'); }
  finally { authLoading.value = false; }
}

async function register() {
  authLoading.value = true;
  try {
    const data = await requestJSON<{ token: string; user: PublicUser }>('/api/user/accounts', {
      method: 'POST', body: JSON.stringify({ username: registerForm.username.trim(), nickname: registerForm.nickname.trim(), password: registerForm.password }),
    });
    localStorage.setItem(userAuthTokenKey, data.token.trim());
    currentUser.value = data.user;
    message.success('注册成功');
    await loadProfile();
  } catch (e) { message.error(e instanceof Error ? e.message : '注册失败'); }
  finally { authLoading.value = false; }
}

async function loadProfile() {
  try {
    const data = await requestJSON<{ user: PublicUser; bindings: Bindings; announcement: UserAnnouncement }>('/api/user/profile');
    currentUser.value = data.user;
    bindings.value = data.bindings || {};
    announcement.value = data.announcement || { enabled: false, content: '', format: 'text' };
    qqValue.value = String(bindings.value.qq || '');
    tgValue.value = String(bindings.value.telegram || '');
  } catch (_) { currentUser.value = null; }
  finally { loading.value = false; }
}

async function loadPlugins() {
  try { openPlugins.value = await requestJSON<OpenPlugin[]>('/api/public/plugins'); }
  catch (_) { openPlugins.value = []; }
}

async function saveQQ() {
  savingQQ.value = true;
  try {
    bindings.value = await requestJSON<Bindings>('/api/user/bindings/qq', {
      method: 'POST', body: JSON.stringify({ value: String(qqValue.value || '').trim() }),
    });
    message.success('QQ 已保存');
  } catch (e) { message.error(e instanceof Error ? e.message : '保存失败'); }
  finally { savingQQ.value = false; }
}

async function saveTG() {
  savingTG.value = true;
  try {
    bindings.value = await requestJSON<Bindings>('/api/user/bindings/telegram', {
      method: 'POST', body: JSON.stringify({ value: String(tgValue.value || '').trim() }),
    });
    message.success('Telegram 已保存');
  } catch (e) { message.error(e instanceof Error ? e.message : '保存失败'); }
  finally { savingTG.value = false; }
}

async function logout() {
  try { await requestJSON<null>('/api/user/sessions/current/deletions', { method: 'POST' }); } catch (_) {}
  localStorage.removeItem(userAuthTokenKey);
  currentUser.value = null;
  message.success('已退出');
}

onMounted(() => { loadProfile(); loadPlugins(); });
</script>

<template>
  <ConfigProvider :locale="zhCN">
    <AntApp>
      <div class="user-page">
        <header class="topbar">
          <div class="brand">MianPro 用户中心</div>
          <Button v-if="currentUser" @click="logout">
            <template #icon><LogOut :size="16" /></template>退出
          </Button>
        </header>

        <main class="content">
          <!-- 未登录 -->
          <div v-if="!loading && !currentUser" class="auth-box">
            <Card :bordered="false">
              <div class="tabs">
                <button :class="{ active: authMode === 'login' }" @click="authMode = 'login'">登录</button>
                <button :class="{ active: authMode === 'register' }" @click="authMode = 'register'">注册</button>
              </div>
              <Form v-if="authMode === 'login'" layout="vertical">
                <Form.Item label="账号" required>
                  <Input v-model:value="loginForm.username" placeholder="账号" />
                </Form.Item>
                <Form.Item label="密码" required>
                  <Input.Password v-model:value="loginForm.password" placeholder="密码" />
                </Form.Item>
                <Button type="primary" block :loading="authLoading" @click="login">登录</Button>
              </Form>
              <Form v-else layout="vertical">
                <Form.Item label="账号" required>
                  <Input v-model:value="registerForm.username" placeholder="3-32 位字母数字" />
                </Form.Item>
                <Form.Item label="昵称">
                  <Input v-model:value="registerForm.nickname" placeholder="选填" />
                </Form.Item>
                <Form.Item label="密码" required>
                  <Input.Password v-model:value="registerForm.password" placeholder="至少 6 位" />
                </Form.Item>
                <Button type="primary" block :loading="authLoading" @click="register">注册</Button>
              </Form>
            </Card>
          </div>

          <!-- 已登录 -->
          <template v-if="!loading && currentUser">
            <Row :gutter="[16, 16]">
              <!-- 用户卡片 -->
              <Col :xs="24" :md="8">
                <Card :bordered="false" class="profile-card" style="height: 100%">
                  <div class="profile-head">
                    <Avatar :size="56" class="avatar">{{ userInitial }}</Avatar>
                    <div>
                      <Typography.Title :level="4" style="margin: 0">{{ currentUser.nickname || currentUser.username }}</Typography.Title>
                      <Typography.Text type="secondary">@{{ currentUser.username }}</Typography.Text>
                    </div>
                  </div>
                  <div class="profile-meta">
                    <div><span>注册时间</span><b>{{ createdDate }}</b></div>
                    <div><span>用户 ID</span><b>{{ currentUser.id.slice(0, 10) }}</b></div>
                  </div>
                </Card>
              </Col>

              <!-- 绑定 -->
              <Col :xs="24" :md="16">
                <Card :bordered="false" title="账号绑定">
                  <div class="bind-row">
                    <div class="bind-icon qq"><Link2 :size="18" /></div>
                    <div class="bind-body">
                      <Typography.Text strong>QQ 号</Typography.Text>
                      <div class="bind-input">
                        <Input v-model:value="qqValue" placeholder="未绑定" style="width: 140px" />
                        <Button size="small" type="primary" :loading="savingQQ" @click="saveQQ">保存</Button>
                      </div>
                    </div>
                  </div>
                  <div class="divider"></div>
                  <div class="bind-row">
                    <div class="bind-icon tg"><QrCode :size="18" /></div>
                    <div class="bind-body">
                      <Typography.Text strong>Telegram</Typography.Text>
                      <div class="bind-input">
                        <Input v-model:value="tgValue" placeholder="未绑定" style="width: 140px" />
                        <Button size="small" type="primary" :loading="savingTG" @click="saveTG">保存</Button>
                      </div>
                    </div>
                  </div>
                </Card>
              </Col>
            </Row>

            <!-- 开放服务 -->
            <Card :bordered="false" style="margin-top: 16px" :title="`开放服务 (${openPlugins.length})`">
              <div v-if="openPlugins.length" class="plugin-list">
                <div v-for="p in openPlugins" :key="p.id" class="plugin-item">
                  <Plug :size="16" class="p-icon" />
                  <div class="p-body">
                    <Typography.Text strong>{{ p.title || p.id }}</Typography.Text>
                    <Typography.Text type="secondary" class="p-desc">{{ p.desc || '暂无介绍' }}</Typography.Text>
                  </div>
                  <Tag v-if="p.version" color="blue">{{ p.version }}</Tag>
                </div>
              </div>
              <Typography.Text v-else type="secondary">暂无开放服务</Typography.Text>
            </Card>

            <!-- 公告 -->
            <Card v-if="announcement.enabled && announceText" :bordered="false" style="margin-top: 16px">
              <template #title><Space><Bell :size="16" />系统公告</Space></template>
              <div v-if="announceHtml" v-html="announceHtml"></div>
              <div v-else>{{ announceText }}</div>
            </Card>
          </template>
        </main>
      </div>
    </AntApp>
  </ConfigProvider>
</template>

<style scoped>
.user-page { min-height: 100vh; background: #f5f7fb; }
.topbar { display: flex; align-items: center; justify-content: space-between; padding: 14px 24px; background: #fff; border-bottom: 1px solid #edf0f5; }
.brand { font-size: 17px; font-weight: 700; }
.content { max-width: 880px; margin: 0 auto; padding: 24px 16px; }

.auth-box { max-width: 380px; margin: 60px auto; }
.tabs { display: grid; grid-template-columns: 1fr 1fr; gap: 4px; padding: 4px; margin-bottom: 20px; background: #f0f2f5; border-radius: 8px; }
.tabs button { height: 34px; border: 0; border-radius: 6px; background: transparent; cursor: pointer; color: #666; }
.tabs button.active { background: #fff; font-weight: 600; box-shadow: 0 2px 6px rgba(0,0,0,0.08); }

.profile-card { border: none; }
.profile-head { display: flex; align-items: center; gap: 14px; margin-bottom: 16px; }
.avatar { background: #1677ff; color: #fff; font-size: 22px; font-weight: 700; }
.profile-meta { display: grid; gap: 8px; padding-top: 14px; border-top: 1px solid #f0f0f0; }
.profile-meta > div { display: flex; justify-content: space-between; font-size: 13px; }
.profile-meta span { color: #8c8c8c; }

.bind-row { display: flex; align-items: center; gap: 12px; padding: 6px 0; }
.bind-icon { display: grid; place-items: center; width: 38px; height: 38px; border-radius: 8px; flex-shrink: 0; }
.bind-icon.qq { background: #e6f4ff; color: #1677ff; }
.bind-icon.tg { background: #e6fffb; color: #13c2c2; }
.bind-body { flex: 1; }
.bind-input { display: flex; gap: 6px; margin-top: 6px; }
.divider { height: 1px; margin: 8px 0; background: #f0f0f0; }

.plugin-list { display: grid; gap: 10px; }
.plugin-item { display: flex; align-items: center; gap: 10px; padding: 10px; background: #f8fafc; border-radius: 8px; }
.p-icon { color: #1677ff; flex-shrink: 0; }
.p-body { flex: 1; min-width: 0; }
.p-desc { font-size: 12px; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
