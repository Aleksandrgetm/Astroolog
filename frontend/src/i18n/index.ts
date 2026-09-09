import { watch } from 'vue'
import { createI18n } from 'vue-i18n'
import { ru as vuetifyRu, lv as vuetifyLv, en as vuetifyEn } from 'vuetify/locale'
import ru from './locales/ru'
import lv from './locales/lv'
import en from './locales/en'

export const supportedLocales = ['ru', 'lv', 'en'] as const
export type Locale = typeof supportedLocales[number]
export const defaultLocale: Locale = 'ru'
export const localeStorageKey = 'astroolog_locale'
export function isLocale(value: unknown): value is Locale {
  return supportedLocales.some(locale => locale === value)
}
function initialLocale(): Locale {
  try {
    const saved = localStorage.getItem(localeStorageKey)
    return isLocale(saved) ? saved : defaultLocale
  } catch {
    return defaultLocale
  }
}
export const i18n = createI18n({
  legacy: false,
  locale: initialLocale(),
  fallbackLocale: defaultLocale,
  messages: {
    ru: { ...ru, $vuetify: vuetifyRu },
    lv: { ...lv, $vuetify: vuetifyLv },
    en: { ...en, $vuetify: vuetifyEn },
  },
})
export function setLocale(locale: Locale) {
  i18n.global.locale.value = locale
}
watch(i18n.global.locale, locale => {
  document.documentElement.lang = locale
  try { localStorage.setItem(localeStorageKey, locale) } catch { /* Storage may be disabled. */ }
}, { immediate: true, flush: 'sync' })
