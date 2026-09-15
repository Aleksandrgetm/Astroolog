<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCatalog } from '../stores/catalog'
import { getLocalizedField } from '../i18n/localized'
import ServiceOptions from './ServiceOptions.vue'
const { t } = useI18n()
const catalog = useCatalog()
onMounted(() => catalog.load())
</script>
<template>
 <p v-if="catalog.loading" role="status">{{ t('serviceCards.loadingSessionOptions') }}</p>
 <div v-else-if="catalog.error" role="alert">{{ t(catalog.error) }} <button @click="catalog.load">{{ t('requestForm.tryAgain') }}</button></div>
 <section v-for="service in catalog.services" v-else :id="service.slug" :key="service.id" class="catalog-section">
  <span class="eyebrow">{{ getLocalizedField(service, 'duration') }}</span>
  <h2>{{ getLocalizedField(service, 'title') }}</h2>
  <p class="catalog-description">{{ getLocalizedField(service, 'description') }}</p>
  <ServiceOptions :service-id="service.id" />
 </section>
</template>
<style scoped>
.catalog-section { padding: 44px 0; scroll-margin-top: 110px; }
.catalog-section + .catalog-section { margin-top: 24px; }
.catalog-description { white-space: pre-line; max-width: 800px; color: var(--muted); line-height: 1.85; margin: 24px 0; }
@media(max-width:600px) { .catalog-section { padding:28px 0; } }
</style>
