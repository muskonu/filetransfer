import './assets/main.css';
import 'webrtc-adapter';

import { createApp } from 'vue';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import App from './App.vue';

// 导入路由文件
import router from './router/index';
import 'element-plus/theme-chalk/el-message.css';

import pinia from '@/stores';

import VueClipboard from "vue3-clipboard";

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(pinia)
app.use(router)
app.use(VueClipboard, {
  autoSetContainer: true,
  appendToBody: true,
})
app.mount('#app')