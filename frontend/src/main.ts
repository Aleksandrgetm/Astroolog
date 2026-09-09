import { createApp } from 'vue'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import router from './router'
import App from './App.vue'
import { i18n } from './i18n'
import './style.css'
createApp(App).use(i18n).use(createPinia()).use(router).use(vuetify).mount('#app')
