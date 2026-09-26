import type { RouteLocation } from 'vue-router'
import { watch } from 'vue'
import { anchorElement, restoreReady } from './scroll'
import { basePath, languagePath, localeFromPath } from '../i18n/urls'
import { createRouter, createWebHistory } from 'vue-router'
import { i18n } from '../i18n'
import { installLocaleRouting } from '../i18n/routing'
import { setPageMeta } from '../i18n/seo'

let scrollRequest = 0
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/services/growth-point', redirect: (to: RouteLocation) => ({ path: languagePath('/services',localeFromPath(to.path)), query: to.query, hash:to.hash }) },
    { path: '/', component: () => import('../pages/HomePage.vue'), meta: { seo: 'home' } },
    { path: '/reviews', component: () => import('../pages/ReviewsPage.vue'), meta: { seo: 'reviews' } },
    { path: '/about', component: () => import('../pages/AboutPage.vue'), meta: { seo: 'about' } },
    { path: '/services', component: () => import('../pages/ServicesPage.vue'), meta: { seo: 'services' } },
    { path: '/directions/:slug', component: () => import('../pages/DirectionPage.vue'), meta: { seo: 'direction' } },
    { path: '/services/:slug', component: () => import('../pages/ServicePage.vue'), meta: { seo: 'service' } },
    { path: '/contacts', component: () => import('../pages/ContactsPage.vue'), meta: { seo: 'contacts' } },
    { path: '/privacy', component: () => import('../pages/PrivacyPage.vue'), meta: { seo: 'privacy' } },
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFoundPage.vue'), meta: { seo: 'notFound' } },
  ].flatMap(route => ['','/lv','/en'].map(prefix=>({...route,path:prefix+route.path}))),
  async scrollBehavior(to, from, savedPosition) {
    const request = ++scrollRequest
    if(savedPosition){await restoreReady();return request === scrollRequest ? savedPosition : false}
    if(basePath(to.path)===basePath(from.path) && to.hash===from.hash)return false
    if(to.hash){const el=await anchorElement(to.hash);if(request !== scrollRequest)return false;if(el)return {el,top:Math.ceil(document.querySelector('header')?.getBoundingClientRect().height||90)+20,behavior:'instant'}}
    return {top:0}

  },
})
installLocaleRouting(router)
function updateRouteMeta() {
  const key = router.currentRoute.value.meta.seo
  if (key && key !== 'direction') setPageMeta(i18n.global.t(`seo.${key}Title`), i18n.global.t(`seo.${key}Description`), key !== 'reviews', key === 'notFound')
}
router.afterEach(updateRouteMeta)
// ServicePage owns its metadata after the API response arrives.
watch(i18n.global.locale, () => {
  if (router.currentRoute.value.meta.seo !== 'service') updateRouteMeta()
})
export default router
