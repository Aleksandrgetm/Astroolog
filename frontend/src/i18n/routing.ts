import type { RouteLocationRaw, Router } from 'vue-router'
import { defaultLocale, isLocale, setLocale } from './index'

export function localeFromQuery(value: unknown) {
  return isLocale(value) ? value : defaultLocale
}

export function localizedRoute(router: Router, target: RouteLocationRaw): RouteLocationRaw {
  const resolved = router.resolve(target)
  const query = { ...resolved.query }
  // An explicit language (including RU) overrides the current language.
  const locale = localeFromQuery('lang' in query ? query.lang : router.currentRoute.value.query.lang)
  if (locale === 'ru') delete query.lang
  else query.lang = locale
  return { ...(typeof target === 'object' ? target : {}), path: resolved.path, query, hash: resolved.hash }
}

export function installLocaleRouting(router: Router) {
  // Only programmatic navigation inherits language. Initial loads and history
  // navigation read their own URL, so a clean URL always means Russian.
  const push = router.push.bind(router)
  const replace = router.replace.bind(router)
  router.push = target => push(localizedRoute(router, target))
  router.replace = target => replace(localizedRoute(router, target))
  router.afterEach((to, _from, failure) => {
    if (!failure) setLocale(localeFromQuery(to.query.lang))
  })
}
