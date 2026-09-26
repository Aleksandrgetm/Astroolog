import { i18n } from './index'
import { basePath, languagePath, localeFromPath, locales } from './urls'
const configured = import.meta.env.VITE_SITE_URL?.replace(/\/$/, '')
function origin() { return configured || window.location.origin }
function meta(key: string, content: string, property = false) {
  const attr = property ? 'property' : 'name'
  let el = document.head.querySelector<HTMLMetaElement>(`meta[${attr}="${key}"]`)
  if (!el) { el = document.createElement('meta'); el.setAttribute(attr, key); document.head.append(el) }
  el.content = content
}
function link(rel: string, href: string, lang?: string) {
  const el = document.createElement('link'); el.rel = rel; el.href = href
  if (lang) el.hreflang = lang
  el.dataset.pageSeo = ''; document.head.append(el)
}
export function setPageMeta(title: string, description: string, appendExpert = true, noindex = false, socialImage?: {path: string; width: number; height: number; alt: string}) {
  const path = window.location.pathname, base = basePath(path), locale = localeFromPath(path)
  const url = origin() + languagePath(base, locale), image = origin() + (socialImage?.path || '/images/optimized/expert-1672.webp')
  document.documentElement.lang = locale
  document.title = appendExpert ? `${title} — ${i18n.global.t('siteLayout.expertName')}` : title
  meta('description', description)
  meta('robots', noindex || base === '/privacy' ? 'noindex, follow' : 'index, follow')
  document.head.querySelectorAll('[data-page-seo], link[rel="canonical"], link[rel="alternate"][hreflang]').forEach(el => el.remove())
  link('canonical', url)
  for (const lang of locales) link('alternate', origin() + languagePath(base, lang), lang)
  link('alternate', origin() + languagePath(base, 'ru'), 'x-default')
  const ogLocales = { ru: 'ru_RU', lv: 'lv_LV', en: 'en_GB' }
  for (const [key, value] of Object.entries({title:document.title, description, type:'website', url, image, 'image:width':String(socialImage?.width || 1672), 'image:height':String(socialImage?.height || 941), 'image:alt':socialImage?.alt || i18n.global.t('siteLayout.expertName'), locale:ogLocales[locale]})) meta(`og:${key}`, value, true)
  document.head.querySelectorAll('meta[property="og:locale:alternate"]').forEach(el=>el.remove())
  for (const lang of locales.filter(l=>l!==locale)) { const el=document.createElement('meta');el.setAttribute('property','og:locale:alternate');el.content=ogLocales[lang];document.head.append(el) }
  for (const [key,value] of Object.entries({card:'summary_large_image',title:document.title,description,image})) meta(`twitter:${key}`,value)
  if(import.meta.env.VITE_GOOGLE_SITE_VERIFICATION) meta('google-site-verification',import.meta.env.VITE_GOOGLE_SITE_VERIFICATION)
  if(base==='/' || base==='/about') {
    const el=document.createElement('script');el.type='application/ld+json';el.dataset.pageSeo=''
    el.textContent=JSON.stringify({'@context':'https://schema.org','@type':'Person',name:i18n.global.t('siteLayout.expertName'),url:origin()+languagePath('/about',locale),image,jobTitle:{ru:'Нумеролог и женский коуч',lv:'Numeroloģe un sieviešu koučs',en:'Numerologist and women’s coach'}[locale]}).replace(/</g, '\\u003c');document.head.append(el)
  }
}
export function setBreadcrumbs(items: {name:string;path:string}[]) {
 document.head.querySelector('#breadcrumbs-schema')?.remove()
 const el=document.createElement('script');el.id='breadcrumbs-schema';el.type='application/ld+json';el.dataset.pageSeo=''
 el.textContent=JSON.stringify({'@context':'https://schema.org','@type':'BreadcrumbList',itemListElement:items.map((item,index)=>({'@type':'ListItem',position:index+1,name:item.name,item:origin()+languagePath(item.path,localeFromPath(window.location.pathname))}))}).replace(/</g, '\\u003c');document.head.append(el)
}
