<script setup lang="ts">
import { basePath } from '../i18n/urls'
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()
import CookieConsent from '../components/CookieConsent.vue'
import { legalDocuments, legalPages, type LegalLocale } from '../legal/content'
import { legalUI } from '../legal/ui'
import { useCookieConsent } from '../legal/useCookieConsent'
const {openSettings}=useCookieConsent()

import { computed } from 'vue'
import { useRoute } from 'vue-router'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
import MobileNavigation from '../components/MobileNavigation.vue'
const route=useRoute()
const links=computed(()=>[[t('siteLayout.aboutMe'),'/about'],[t('siteLayout.areasOfFocus'),'/#directions'],[t('siteLayout.sessions'),'/services'],[t('siteLayout.testimonials'),'/reviews'],[t('siteLayout.contact'),'/contacts']])
function isNavigationActive(target: string) {
  const [path, hash] = target.split('#')
  if (hash === 'directions' && basePath(route.path).startsWith('/directions/')) return true
  if (hash) return basePath(route.path) === path && route.hash === `#${hash}`
  return basePath(route.path) === path || basePath(route.path).startsWith(`${path}/`)
}
</script>
<template>
<a class="skip-link" href="#main">{{ t('siteLayout.skipToContent') }}</a>
<header class="site-header"><div class="container header-inner"><LocaleLink to="/" class="brand" :aria-label="t('siteLayout.homeLabel')"><span class="monogram">{{ t('siteLayout.monogramFirst') }}<span>{{ t('siteLayout.monogramSecond') }}</span></span><span class="brand-name">{{ t('siteLayout.expertName') }}<small>{{ t('siteLayout.numerologistCoach') }}</small></span></LocaleLink><nav class="desktop-nav" :aria-label="t('siteLayout.mainNavigation')"><LocaleLink v-for="l in links" :key="l[0]" :to="l[1]!" active-class="" exact-active-class="" :class="{ 'is-nav-active': isNavigationActive(l[1]!) }" :aria-current="isNavigationActive(l[1]!) ? 'location' : undefined">{{l[0]}}</LocaleLink></nav><LanguageSwitcher class="desktop-language"/><LocaleLink to="/contacts?booking=1" class="button header-cta">{{ t('siteLayout.bookASession') }} <span><ArrowIcon /></span></LocaleLink><MobileNavigation :links="links" /></div></header>

<main id="main" tabindex="-1"><slot /></main>
<footer class="site-footer"><div class="container footer-top"><LocaleLink to="/" class="brand"><span class="monogram">{{ t('siteLayout.monogramFirst') }}<span>{{ t('siteLayout.monogramSecond') }}</span></span><span class="brand-name">{{ t('siteLayout.expertName') }}<small>{{ t('siteLayout.numerologistCoach') }}</small></span></LocaleLink><p>{{ t('siteLayout.closerToYourself') }}<br>{{ t('siteLayout.closerToALifeYouChoose') }}</p><nav :aria-label="t('siteLayout.footerNavigation')"><LocaleLink to="/about">{{ t('siteLayout.aboutMe') }}</LocaleLink><LocaleLink to="/services">{{ t('siteLayout.sessions') }}</LocaleLink><LocaleLink to="/contacts">{{ t('siteLayout.getInTouch') }} <ArrowIcon /></LocaleLink></nav></div><div class="container footer-bottom"><span>© {{new Date().getFullYear()}} {{ t('siteLayout.expertName') }}</span><span>{{ t('siteLayout.onlineInYourOwnSpace') }}</span><nav class="footer-legal" :aria-label="legalUI[locale as LegalLocale].legalNavigation"><LocaleLink v-for="page in legalPages" :key="page" :to="'/'+page">{{legalDocuments[locale as LegalLocale][page].title}}</LocaleLink><button type="button" @click="openSettings">{{legalUI[locale as LegalLocale].settings}}</button></nav></div></footer>
<CookieConsent />
</template>

<style scoped>
.footer-legal{display:flex;gap:12px 24px;flex-wrap:wrap;align-items:center;}.footer-legal button{font:inherit;color:inherit;background:none;border:0;padding:0;cursor:pointer;text-align:left;}.footer-legal>*{overflow-wrap:anywhere;}.footer-legal button:focus-visible{outline:2px solid currentColor;outline-offset:4px;}
</style>
