import type { RouteLocation } from 'vue-router'
import { watch } from 'vue'
import { anchorElement, restoreReady } from './scroll'
import { basePath, languagePath, localeFromPath } from '../i18n/urls'
import { createRouter, createWebHistory } from 'vue-router'
import { i18n } from '../i18n'
import { installLocaleRouting } from '../i18n/routing'
import { loadSession, adminSession } from '../admin/api'
import { setPageMeta } from '../i18n/seo'

let scrollRequest = 0
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {path:'/admin/login',component:()=>import('../admin/LoginPage.vue'),meta:{admin:true}},
    ...['','requests','audit'].map(area=>({path:'/admin'+(area?'/'+area:''),component:()=>import('../admin/AdminPage.vue'),meta:{admin:true,adminArea:area||'dashboard'}})),
    {path:'/admin/content/:page?',component:()=>import('../admin/ContentOverview.vue'),meta:{admin:true}},
    {path:'/admin/content/:page/:section',component:()=>import('../admin/SectionEditor.vue'),meta:{admin:true}},
    ...['services','directions','reviews','certificates','seo'].map(area=>({path:'/admin/'+area,component:()=>import('../admin/EntityList.vue'),meta:{admin:true,adminArea:area}})),
    ...['services','options','directions','reviews','certificates','seo'].map(area=>({path:'/admin/'+area+'/:id',component:()=>import('../admin/EntityEditorPage.vue'),meta:{admin:true,adminArea:area}})),
    {path:'/admin/settings',component:()=>import('../admin/EntityEditorPage.vue'),meta:{admin:true,adminArea:'settings'}},
    {path:'/admin/media',component:()=>import('../admin/MediaLibrary.vue'),meta:{admin:true}},
    ...[
    { path: '/services/growth-point', redirect: (to: RouteLocation) => ({ path: languagePath('/services',localeFromPath(to.path)), query: to.query, hash:to.hash }) },
    { path: '/', component: () => import('../pages/HomePage.vue'), meta: { seo: 'home' } },
    { path: '/reviews', component: () => import('../pages/ReviewsPage.vue'), meta: { seo: 'reviews' } },
    { path: '/about', component: () => import('../pages/AboutPage.vue'), meta: { seo: 'about' } },
    { path: '/services', component: () => import('../pages/ServicesPage.vue'), meta: { seo: 'services' } },
    { path: '/directions/:slug', component: () => import('../pages/DirectionPage.vue'), meta: { seo: 'direction' } },
    { path: '/services/:slug', component: () => import('../pages/ServicePage.vue'), meta: { seo: 'service' } },
    { path: '/contacts', component: () => import('../pages/ContactsPage.vue'), meta: { seo: 'contacts' } },
    ...['privacy','cookies','terms'].map(page=>({ path: '/'+page, component: () => import('../pages/LegalPage.vue'), meta: { seo: page } })),
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFoundPage.vue'), meta: { seo: 'notFound' } },
  ].flatMap(route => ['','/lv','/en'].map(prefix=>({...route,path:prefix+route.path}))),
  ],
  async scrollBehavior(to, from, savedPosition) {
    const request = ++scrollRequest
    if(savedPosition){await restoreReady();return request === scrollRequest ? savedPosition : false}
    if(basePath(to.path)===basePath(from.path) && to.hash===from.hash)return false
    if(to.hash){const el=await anchorElement(to.hash);if(request !== scrollRequest)return false;if(el)return {el,top:Math.ceil(document.querySelector('header')?.getBoundingClientRect().height||90)+20,behavior:'instant'}}
    return {top:0}

  },
})
installLocaleRouting(router)
router.beforeEach(async to=>{
 if(to.meta.admin){
  document.title='Astroolog — управление сайтом'
  document.querySelectorAll('link[rel=canonical],link[hreflang],[data-page-seo],#breadcrumbs-schema').forEach(el=>el.remove())
  document.querySelector('meta[name="robots"]')?.setAttribute('content','noindex, nofollow')
  if(to.path!=='/admin/login' && !adminSession.user && !await loadSession())return '/admin/login'
 }
})
function updateRouteMeta() {
  const key = router.currentRoute.value.meta.seo
  if (key && key !== 'direction' && !['privacy','cookies','terms'].includes(String(key))) setPageMeta(i18n.global.t(`seo.${key}Title`), i18n.global.t(`seo.${key}Description`), key !== 'reviews', key === 'notFound')
}
router.afterEach(updateRouteMeta)
// ServicePage owns its metadata after the API response arrives.
watch(i18n.global.locale, () => {
  if (router.currentRoute.value.meta.seo !== 'service') updateRouteMeta()
})
export default router
