import type { RouteLocationRaw, Router } from 'vue-router'
import { defaultLocale, isLocale, setLocale } from './index'
import { basePath, languagePath, localeFromPath } from './urls'
export function localeFromQuery(value:unknown){return isLocale(value)?value:defaultLocale}
export function localizedRoute(router:Router,target:RouteLocationRaw):RouteLocationRaw {
 const resolved=router.resolve(target),query={...resolved.query}
 const explicitPrefix=/^\/(lv|en)(\/|$)/.test(resolved.path)
 const locale='lang' in query?localeFromQuery(query.lang):explicitPrefix?localeFromPath(resolved.path):localeFromPath(router.currentRoute.value.path)
 delete query.lang
 return {path:languagePath(resolved.path,locale),query,hash:resolved.hash}
}
export function installLocaleRouting(router:Router){
 const push=router.push.bind(router),replace=router.replace.bind(router)
 router.push=target=>push(localizedRoute(router,target))
 router.replace=target=>replace(localizedRoute(router,target))
 router.beforeEach(to=>{
  const locale='lang' in to.query?localeFromQuery(to.query.lang):localeFromPath(to.path)
  const query={...to.query};delete query.lang
  const path=languagePath(basePath(to.path),locale)
  if(path!==to.path || 'lang' in to.query)return {path,query,hash:to.hash,replace:true}
 })
 router.afterEach((to,_from,failure)=>{if(!failure)setLocale(localeFromPath(to.path))})
}
