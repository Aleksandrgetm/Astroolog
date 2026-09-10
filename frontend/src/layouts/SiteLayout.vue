<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
const menu=ref(false),route=useRoute()
const links=computed(()=>[[t('siteLayout.aboutMe'),'/about'],[t('siteLayout.areasOfFocus'),'/#directions'],[t('siteLayout.sessions'),'/services'],[t('siteLayout.testimonials'),'/#testimonials'],[t('siteLayout.contact'),'/contacts']])
watch(()=>route.fullPath,()=>menu.value=false)
</script>
<template>
<a class="skip-link" href="#main">{{ t('siteLayout.skipToContent') }}</a>
<header class="site-header"><div class="container header-inner"><router-link to="/" class="brand" :aria-label="t('siteLayout.homeLabel')"><span class="monogram">{{ t('siteLayout.monogramFirst') }}<span>{{ t('siteLayout.monogramSecond') }}</span></span><span class="brand-name">{{ t('siteLayout.expertName') }}<small>{{ t('siteLayout.numerologistCoach') }}</small></span></router-link><nav class="desktop-nav" :aria-label="t('siteLayout.mainNavigation')"><router-link v-for="l in links" :key="l[0]" :to="l[1]!">{{l[0]}}</router-link></nav><LanguageSwitcher class="desktop-language"/><router-link to="/contacts?booking=1" class="button header-cta">{{ t('siteLayout.bookASession') }} <span>↗</span></router-link><v-btn class="mobile-toggle" icon="mdi-menu" :aria-label="t('siteLayout.openMenu')" @click="menu=true" /></div></header>
<v-dialog v-model="menu" fullscreen transition="dialog-bottom-transition"><div class="mobile-menu"><div class="menu-top"><span class="eyebrow">{{ t('siteLayout.expertNameUppercase') }}</span><v-btn icon="mdi-close" variant="text" :aria-label="t('siteLayout.closeMenu')" @click="menu=false" /></div><LanguageSwitcher class="mobile-language"/><nav :aria-label="t('siteLayout.mobileNavigation')"><router-link v-for="l in links" :key="l[0]" :to="l[1]!" @click="menu=false">{{l[0]}} <span aria-hidden="true">↗</span></router-link></nav><router-link class="button" to="/contacts?booking=1" @click="menu=false">{{ t('siteLayout.bookASession13') }}</router-link></div></v-dialog>
<main id="main" tabindex="-1"><slot /></main>
<footer class="site-footer"><div class="container footer-top"><router-link to="/" class="brand"><span class="monogram">{{ t('siteLayout.monogramFirst') }}<span>{{ t('siteLayout.monogramSecond') }}</span></span><span class="brand-name">{{ t('siteLayout.expertName') }}<small>{{ t('siteLayout.numerologistCoach') }}</small></span></router-link><p>{{ t('siteLayout.closerToYourself') }}<br>{{ t('siteLayout.closerToALifeYouChoose') }}</p><nav :aria-label="t('siteLayout.footerNavigation')"><router-link to="/about">{{ t('siteLayout.aboutMe') }}</router-link><router-link to="/services">{{ t('siteLayout.sessions') }}</router-link><router-link to="/contacts">{{ t('siteLayout.getInTouch') }}</router-link></nav></div><div class="container footer-bottom"><span>© {{new Date().getFullYear()}} {{ t('siteLayout.expertName') }}</span><span>{{ t('siteLayout.onlineInYourOwnSpace') }}</span><router-link to="/privacy">{{ t('siteLayout.privacy') }}</router-link></div></footer>
</template>
