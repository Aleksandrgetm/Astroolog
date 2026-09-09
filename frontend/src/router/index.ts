import { watch } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { i18n } from '../i18n'
import { setPageMeta } from '../i18n/seo'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: () => import('../pages/HomePage.vue'), meta: { seo: 'home' } },
    { path: '/about', component: () => import('../pages/AboutPage.vue'), meta: { seo: 'about' } },
    { path: '/services', component: () => import('../pages/ServicesPage.vue'), meta: { seo: 'services' } },
    { path: '/services/:slug', component: () => import('../pages/ServicePage.vue'), meta: { seo: 'service' } },
    { path: '/contacts', component: () => import('../pages/ContactsPage.vue'), meta: { seo: 'contacts' } },
    { path: '/privacy', component: () => import('../pages/PrivacyPage.vue'), meta: { seo: 'privacy' } },
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFoundPage.vue'), meta: { seo: 'notFound' } },
  ],
  scrollBehavior(to) {
    return to.hash ? { el: to.hash, top: 110, behavior: 'smooth' } : { top: 0 }
  },
})
function updateRouteMeta() {
  const key = router.currentRoute.value.meta.seo
  if (key) setPageMeta(i18n.global.t(`seo.${key}Title`), i18n.global.t(`seo.${key}Description`))
}
router.afterEach(updateRouteMeta)
// ServicePage owns its metadata after the API response arrives.
watch(i18n.global.locale, () => {
  if (router.currentRoute.value.meta.seo !== 'service') updateRouteMeta()
})
export default router
