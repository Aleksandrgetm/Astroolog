<script setup lang="ts">
import { formatPrice } from '../utils/prices'
import { expertImage } from '../utils/expertImage'
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

import { getLocalizedField } from '../i18n/localized'
import { onMounted } from 'vue'
import { useCatalog } from '../stores/catalog'
const catalog=useCatalog();onMounted(()=>catalog.load())
</script>
<template><div v-if="catalog.loading" class="state" role="status"><v-progress-circular indeterminate size="24" /> {{ t('serviceCards.loadingSessionOptions') }}</div><div v-else-if="catalog.error" class="state" role="alert">{{t(catalog.error)}} <v-btn variant="text" @click="catalog.load">{{ t('requestForm.tryAgain') }}</v-btn></div><div v-else-if="!catalog.services.length" class="state">{{ t('serviceCards.sessionOptionsWillBeAvailableSoon') }} <LocaleLink to="/contacts">{{ t('serviceCards.askMeAQuestion') }}</LocaleLink>.</div><div v-else class="service-grid"><article v-for="(s,i) in catalog.services" :key="s.id" class="service-card"><LocaleLink :to="`/services/${s.slug}`" tabindex="-1" aria-hidden="true"><img :src="expertImage(s.image)" :class="{ 'expert-photo': expertImage(s.image) === '/images/Hero-new.png' }" :alt="getLocalizedField(s, 'title')" loading="lazy" width="520" height="360" /></LocaleLink><div class="service-content"><span class="eyebrow">0{{i+1}} / {{getLocalizedField(s, 'duration')}}</span><h3><LocaleLink :to="`/services/${s.slug}`">{{getLocalizedField(s, 'title')}}</LocaleLink></h3><p>{{getLocalizedField(s, 'short_description')}}</p><div class="service-base-price">{{ s.price === 0 ? t('catalog.free') : t('catalog.fromPrice', { price: formatPrice(s.price, locale, t('catalog.free'), t('bookingChoice.pricePending')) }) }}</div><LocaleLink class="text-link" :to="`/services/${s.slug}`">{{ t('serviceCards.learnMore') }} <span><ArrowIcon /></span></LocaleLink></div></article></div></template>

<style scoped>
.service-base-price { margin: 12px 0; color: var(--brown); font-size: 14px; }
</style>
