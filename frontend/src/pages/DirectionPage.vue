<script setup lang="ts">
import { content, ctaDestination, imageSrcset } from '../stores/content'
import { computed, onMounted, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { findDirection } from '../data/directions'
import { useCatalog } from '../stores/catalog'
import { setPageMeta, setBreadcrumbs } from '../i18n/seo'
import ServiceOptions from '../components/ServiceOptions.vue'
import FinalCta from '../components/FinalCta.vue'
import NotFoundPage from './NotFoundPage.vue'

const route = useRoute()
const { t, locale } = useI18n()
const catalog = useCatalog()
const direction = computed(() => findDirection(String(route.params.slug)))
const key = computed(() => `directions.${direction.value?.contentKey}`)
const service = computed(() => catalog.services.find(item => item.id === content.directions.find(d=>d.slug===direction.value?.slug)?.service_id || item.slug === direction.value?.relatedServiceSlug))
const availableOptions = computed(() => catalog.options.filter(option => option.service_id === service.value?.id && direction.value?.relatedBookingOptionCodes.includes(option.code)))
const bookingTarget = computed(() => ctaDestination('global.cta')!=='/contacts?booking=1'?ctaDestination('global.cta'): ({ path: '/contacts', query: { booking: '1', ...(availableOptions.value[0] ? { option: availableOptions.value[0].code } : {}) } }))
onMounted(() => catalog.load())
watchEffect(() => {
  if (route.meta.seo !== 'direction') return
  void route.fullPath
  void locale.value
  const item = direction.value
  if (!item) {
    setPageMeta(t('seo.notFoundTitle'), t('seo.notFoundDescription'), true, true)
    return
  }
  setPageMeta(t(`${key.value}.seoTitle`), t(`${key.value}.seoDescription`), true, false, {
    path: item.image, width: item.width, height: item.height, alt: t(`${key.value}.imageAlt`),
  })
  setBreadcrumbs([
    { name: t('navigation.home'), path: '/' },
    { name: t('directions.shared.label'), path: '/#directions' },
    { name: t(`${key.value}.title`), path: `/directions/${item.slug}` },
  ])
})
</script>

<template>
  <template v-if="direction">
    <article class="section container direction-page">
      <nav class="breadcrumbs" :aria-label="t('navigation.breadcrumbs')">
        <LocaleLink to="/">{{ t('navigation.home') }}</LocaleLink><ArrowIcon direction="right" />
        <LocaleLink to="/#directions">{{ t('directions.shared.label') }}</LocaleLink><ArrowIcon direction="right" />
        <span aria-current="page">{{ t(`${key}.title`) }}</span>
      </nav>
      <div class="service-detail direction-hero">
        <div>
          <span class="eyebrow">{{ t(`${key}.eyebrow`) }}</span>
          <h1>{{ t(`${key}.title`) }}</h1>
          <p class="lead">{{ t(`${key}.lead`) }}</p>
          <p class="direction-context">{{ t(`${key}.context`) }}</p>
        </div>
        <img :src="direction.image" :srcset="imageSrcset(direction.image)" sizes="(max-width:600px) 100vw, 50vw" :alt="t(`${key}.imageAlt`)" :width="direction.width" :height="direction.height" :style="{ objectPosition: direction.imagePosition }" loading="eager" decoding="async" fetchpriority="high" />
      </div>
      <section class="direction-formats" aria-labelledby="formats-title">
        <h2 id="formats-title">{{ t('directions.shared.formatsTitle') }}</h2>
        <p class="direction-format-intro">{{ t(`${key}.formatNote`) }}</p>
        <p v-if="catalog.loading" role="status">{{ t('serviceCards.loadingSessionOptions') }}</p>
        <div v-else-if="catalog.error" role="alert">
          <p>{{ t(catalog.error) }}</p><button type="button" class="text-link" @click="catalog.load({ force: true })">{{ t('requestForm.tryAgain') }}</button>
        </div>
        <ServiceOptions v-else-if="service && availableOptions.length" :service-id="service.id" :option-codes="direction.relatedBookingOptionCodes" />
        <p v-else>{{ t('directions.shared.unavailable') }}</p>
        <div class="direction-related-links">
          <LocaleLink :to="`/services/${service?.slug || direction.relatedServiceSlug}`" class="text-link">{{ t('directions.shared.serviceLink') }} <ArrowIcon /></LocaleLink>
          <LocaleLink v-if="direction.id === 'purpose-money'" to="/services/full-matrix" class="text-link">{{ t('directions.shared.fullLink') }} <ArrowIcon /></LocaleLink>
        </div>
      </section>
    </article>
    <FinalCta>
      <template #actions>
        <LocaleLink :to="bookingTarget" class="button button-light">{{ t('finalCta.iWantToUnderstandMyself') }} <ArrowIcon /></LocaleLink>
        <LocaleLink :to="ctaDestination('global.cta','secondary_destination','question')" class="quiet-link">{{ t('finalCta.askAQuestionFirst') }} <ArrowIcon /></LocaleLink>
      </template>
    </FinalCta>
  </template>
  <NotFoundPage v-else />
</template>

<style scoped>
.breadcrumbs { display:flex;flex-wrap:wrap;align-items:center;gap:8px;margin-bottom:24px;font-size:14px; }
.breadcrumbs > * { min-width:0;overflow-wrap:anywhere; }
.breadcrumbs svg { width:18px;height:18px; }
.direction-hero > div { min-width:0; }
.direction-hero h1 { overflow-wrap:break-word; }
.direction-hero .lead,
.direction-context,
.direction-format-intro { white-space:pre-line; }
.direction-context { color:var(--muted);line-height:1.85; }
.direction-formats h2 { font-size:clamp(32px,3.2vw,44px); }
.direction-format-intro { margin-top:24px;line-height:1.85;color:var(--muted);max-width:850px; }
.direction-formats { margin-top:56px;padding-top:40px;border-top:1px solid var(--line); }
.direction-related-links { display:flex;flex-wrap:wrap;gap:24px;margin-top:28px; }
.direction-related-links a { max-width:100%; }
.direction-related-links svg { flex-shrink:0; }
@media(max-width:900px) {
  .direction-hero { grid-template-columns:1fr; }
  .direction-hero > img { grid-row:auto;height:auto;max-height:490px;aspect-ratio:3 / 2; }
}
@media(max-width:600px) {
  .direction-formats { margin-top:36px;padding-top:32px; }
  .direction-hero h1 { font-size:42px; }
  .direction-hero > img { aspect-ratio:3 / 2; }
  .breadcrumbs { font-size:12px; }
}
</style>
