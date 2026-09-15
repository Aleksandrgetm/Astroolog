import { i18n } from './index'
export function setPageMeta(title: string, description: string, appendExpert = true) {
  document.title = appendExpert ? `${title} — ${i18n.global.t('siteLayout.expertName')}` : title
  document.querySelector('meta[name="description"]')?.setAttribute('content', description)
  document.querySelector('meta[property="og:title"]')?.setAttribute('content', document.title)
  document.querySelector('meta[property="og:description"]')?.setAttribute('content', description)
  document.querySelector('meta[property="og:locale"]')?.setAttribute('content', {ru:'ru_RU',lv:'lv_LV',en:'en_GB'}[i18n.global.locale.value])
}
