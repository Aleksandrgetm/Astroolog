<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { certificates, certificateImage } from '../data/certificates'
const { t, locale } = useI18n()
</script>

<template>
  <section class="certificates" aria-labelledby="certificates-title">
    <span class="eyebrow">{{ t('certificates.eyebrow') }}</span>
    <h2 id="certificates-title">{{ t('certificates.title') }}</h2>
    <p class="certificates-intro">{{ t('certificates.intro') }}</p>
    <div class="certificate-grid">
      <article v-for="item in certificates" :key="item.id" class="certificate-card">
        <div class="certificate-preview">
          <img :src="certificateImage(item, locale)" :alt="t(`certificates.${item.alt}`, { name: t('siteLayout.expertName') })" width="1600" height="1131" loading="lazy" />
        </div>
        <div class="certificate-caption"><span class="eyebrow">{{ item.year }}</span><h3>{{ t(`certificates.${item.title}`) }}</h3><p>{{ t(`certificates.${item.subtitle}`) }}</p></div>
      </article>
    </div>
  </section>
</template>

<style scoped>
.certificates { margin-top: 64px; padding-top: 48px; border-top: 1px solid var(--line); }
.certificates h2 { max-width: 760px; }
.certificates-intro { max-width: 690px; margin: 22px 0 32px; color: var(--muted); line-height: 1.8; }
.certificate-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 28px; }
.certificate-card { padding: 18px; border: 1px solid var(--line); background: #fffdfa; min-width: 0; }
.certificate-preview { display: block; width: 100%; background: transparent; border: 0; }
.certificate-preview img { display: block; width: 100%; height: auto; aspect-ratio: 1600 / 1131; object-fit: contain; }
.certificate-caption { padding: 22px 4px 6px; }
.certificate-caption .eyebrow { margin-bottom: 10px; }
.certificate-caption h3 { font-size: 28px; line-height: 1.2; }
.certificate-caption p { margin: 10px 0 0; font-size: 13px; line-height: 1.7; color: var(--muted); }
@media (max-width: 600px) {
  .certificates { margin-top: 40px; padding-top: 32px; }
  .certificate-grid { grid-template-columns: 1fr; gap: 24px; }
  .certificate-card { padding: 12px; }
  .certificates h2 { font-size: 36px; }
}
</style>
