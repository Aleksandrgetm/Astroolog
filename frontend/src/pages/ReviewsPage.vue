<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import FinalCta from '../components/FinalCta.vue'
import { reviews } from '../data/reviews'
const { t } = useI18n()
</script>

<template>
  <section class="section container reviews-page">
    <div class="page-heading reviews-intro">
      <span class="eyebrow">{{ t('reviews.eyebrow') }}</span>
      <h1>{{ t('reviews.title') }}</h1>
      <p>{{ t('reviews.intro') }}</p>
    </div>
    <div class="full-reviews">
      <article v-for="review in reviews" :key="review.id" :id="review.id" class="full-review">
        <span class="quote-mark" aria-hidden="true">“</span>
        <blockquote><p v-for="(paragraph, index) in t(review.textKey).split('\n\n')" :key="index">{{ paragraph }}</p></blockquote>
        <div class="full-review-author"><span class="anonymous-avatar" aria-hidden="true"></span>{{ t('home.testimonials.author') }}</div>
      </article>
    </div>
  </section>
  <FinalCta class="reviews-cta">
    <template #eyebrow>{{ t('reviews.ctaEyebrow') }}</template>
    <template #heading>{{ t('reviews.ctaTitle') }}</template>
    <template #description>{{ t('reviews.ctaText') }}</template>
    <template #actions><LocaleLink class="button button-light" to="/contacts?booking=1">{{ t('siteLayout.bookASession') }} <span><ArrowIcon /></span></LocaleLink></template>
  </FinalCta>
</template>

<style scoped>
.reviews-intro { max-width:760px; }
.full-reviews { columns:2; column-gap:28px; }
.full-review { display:inline-block; width:100%; break-inside:avoid; margin:0 0 28px; padding:32px; background:#f4efe8; border:1px solid #e9e1d5; vertical-align:top; overflow-wrap:anywhere; }
.quote-mark { display:block; font-size:52px; }
blockquote { margin:12px 0 28px; }
blockquote p { font-size:16px; line-height:1.85; color:var(--muted); margin:0 0 18px; }
blockquote p:last-child { margin-bottom:0; }
.full-review-author { display:flex; align-items:center; gap:12px; padding-top:20px; border-top:1px solid #e1d5c7; font-size:12px; }
.anonymous-avatar { width:28px; height:28px; border-radius:50%; background:#e5d9c9; flex-shrink:0; }
.reviews-cta :deep(h2) { max-width:760px; margin-inline:auto; text-wrap:balance; }
@media(max-width:600px) { .full-reviews { columns:1; } .full-review { padding:24px; margin-bottom:20px; } }
</style>
