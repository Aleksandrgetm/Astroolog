import { createApp } from 'vue'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import router from './router'
import App from './App.vue'
import { loadContent } from './stores/content'
import LocaleLink from './components/LocaleLink.vue'
import ArrowIcon from './components/ArrowIcon.vue'
import { i18n } from './i18n'
import './style.css'
const app = createApp(App).use(i18n).use(createPinia()).use(router).use(vuetify).component('LocaleLink', LocaleLink).component('ArrowIcon', ArrowIcon)
Promise.all([router.isReady(), loadContent()]).then(() => { app.mount('#app'); router.replace(router.currentRoute.value.fullPath) })
