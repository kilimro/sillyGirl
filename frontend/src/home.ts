import { createApp } from 'vue';
import 'ant-design-vue/dist/reset.css';
import Home from './Home.vue';
import './styles.css';
import { bootWatchdog } from './watchdog';

console.log('欢迎使用 MianPro！');
bootWatchdog('home');
createApp(Home).mount('#home-root');
