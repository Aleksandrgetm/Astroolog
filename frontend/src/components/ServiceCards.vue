<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

import { getLocalizedField } from '../i18n/localized'
import { onMounted } from 'vue'
import { useCatalog } from '../stores/catalog'
const catalog=useCatalog();onMounted(()=>catalog.load())
</script>
<template><div v-if="catalog.loading" class="state" role="status"><v-progress-circular indeterminate size="24" /> {{ t('serviceCards.loadingSessionOptions') }}</div><div v-else-if="catalog.error" class="state" role="alert">{{t(catalog.error)}} <v-btn variant="text" @click="catalog.load">{{ t('requestForm.tryAgain') }}</v-btn></div><div v-else-if="!catalog.services.length" class="state">{{ t('serviceCards.sessionOptionsWillBeAvailableSoon') }} <router-link to="/contacts">{{ t('serviceCards.askMeAQuestion') }}</router-link>.</div><div v-else class="service-grid"><article v-for="(s,i) in catalog.services" :key="s.id" class="service-card"><router-link :to="`/services/${s.slug}`" tabindex="-1" aria-hidden="true"><img :src="s.image" :alt="getLocalizedField(s, 'title')" loading="lazy" width="520" height="360" /></router-link><div class="service-content"><span class="eyebrow">0{{i+1}} / {{getLocalizedField(s, 'duration')}}</span><h3><router-link :to="`/services/${s.slug}`">{{getLocalizedField(s, 'title')}}</router-link></h3><p>{{getLocalizedField(s, 'short_description')}}</p><router-link class="text-link" :to="`/services/${s.slug}`">{{ t('serviceCards.learnMore') }} <span>↗</span></router-link></div></article></div></template>
