import { mountSillyGirlVersion, SILLYGIRL_VERSION } from './version';

declare global {
  interface Window {
    __SILLYGIRL_WATCHDOG__?: Promise<void>;
  }
}

// 遥测已禁用：不再向任何第三方端点上报浏览行为。
export function bootWatchdog(_userType: 'home' | 'user' | 'admin') {
  if (typeof window === 'undefined') return;
  mountSillyGirlVersion();
  window.__SILLYGIRL_WATCHDOG__ = Promise.resolve();
}
