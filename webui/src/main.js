import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import LoginView from './views/LoginView.vue'
import Header from './components/Header.vue'

import './assets/css/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("LoginView", LoginView)
app.component("Header", Header)
app.use(router)
app.mount('#app')
