<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCatalog } from '../stores/catalog'
import { getLocalizedField } from '../i18n/localized'
import { formatPrice } from '../utils/prices'
const optionImages: Record<string, string> = {
 personality: '/images/Hero-new.png',
 finances: '/images/journal-placeholder.webp',
 relationships: '/images/relationships-placeholder.webp',
 'child-matrix': '/images/child-placeholder.webp',
 'reading-call-40': '/images/Hero-new.png',
}
const props = defineProps<{ serviceId: number }>()
const { t, locale } = useI18n()
const catalog = useCatalog()
onMounted(() => catalog.load())
const options = computed(() => catalog.options.filter(o => o.service_id === props.serviceId))
</script>
<template>
 <div class="service-options">
  <div v-for="option in options" :id="option.code" :key="option.code" class="service-option-row">
   <img v-if="optionImages[option.code]" :src="optionImages[option.code]" alt="" width="240" height="160" loading="lazy" class="option-image" />
   <div class="option-description"><h3>{{ getLocalizedField(option, 'title') }}</h3><p>{{ getLocalizedField(option, 'format') }}</p></div>
   <div class="service-option-action"><span class="option-price">{{ formatPrice(option.price, locale, t('catalog.free'), t('bookingChoice.pricePending')) }}</span><LocaleLink :to="{path:'/contacts',query:{booking:'1',option:option.code}}" class="text-link">{{ t('siteLayout.bookASession') }} <ArrowIcon /></LocaleLink></div>
  </div>
 </div>
</template>
<style scoped>
.option-image { width:160px; height:110px; object-fit:cover; flex-shrink:0; }
.option-image[src$="Hero-new.png"] { object-position:75% center; }
.option-description { flex:1; min-width:0; }
.service-options { margin-top: 28px; }
.service-option-row { display: flex; align-items: center; justify-content: space-between; gap: 24px; padding: 24px 0; border-top: 1px solid var(--line); scroll-margin-top: 120px; }
.service-option-row h3 { font-size: 27px; line-height: 1.2; }
.service-option-row p { font-size: 12px; color: var(--muted); margin: 10px 0 0; }
.service-option-action { flex-shrink: 0; text-align: right; }
.option-price { display: block; color: var(--brown); font-size: 20px; margin-bottom: 8px; }
@media(max-width:600px) { .option-image { width:100%; height:auto; aspect-ratio:3 / 2; } .service-option-row { align-items: flex-start; flex-direction: column; gap: 14px; } .service-option-action { display:flex; align-items:center; justify-content:space-between; width:100%; gap:12px; } .option-price { margin:0; } .service-option-row h3 { font-size:25px; } }
</style>
