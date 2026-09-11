import { watch } from 'vue'
import { createI18n } from 'vue-i18n'
import { ru as vuetifyRu, lv as vuetifyLv, en as vuetifyEn } from 'vuetify/locale'
import ru from './locales/ru'
import lv from './locales/lv'
import en from './locales/en'

export const supportedLocales = ['ru', 'lv', 'en'] as const
export type Locale = typeof supportedLocales[number]
export const defaultLocale: Locale = 'ru'
export function isLocale(value: unknown): value is Locale {
  return supportedLocales.some(locale => locale === value)
}
function initialLocale(): Locale {
  const values = new URLSearchParams(window.location.search).getAll('lang')
  return values.length === 1 && isLocale(values[0]) ? values[0] : defaultLocale
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
}, { immediate: true, flush: 'sync' })
