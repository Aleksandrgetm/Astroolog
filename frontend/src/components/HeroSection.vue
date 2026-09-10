<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <section class="home-hero" aria-labelledby="hero-title">
    <div class="home-hero__visual">
      <img
        class="home-hero__photo"
        src="/images/hero/hero-main.webp"
        :alt="t('home.heroExpertPortrait', { name: t('siteLayout.expertName') })"
        width="1672"
        height="941"
        loading="eager"
        fetchpriority="high"
      />
    </div>

    <div class="container home-hero__inner">
      <div class="home-hero__copy">
        <span class="eyebrow home-hero__eyebrow">
          <span class="home-hero__star" aria-hidden="true">✧</span>
          {{ t('home.numerologySupportiveCoaching') }}
        </span>
        <h1 id="hero-title">
          {{ t('home.understandYourself') }}<br />
          {{ t('home.discoverYourStrength') }}<br />
          <em>{{ t('home.findYourOwnPath') }}</em>
        </h1>
        <p class="home-hero__description">{{ t('home.iHelpWomenListenToThemselvesSee') }}</p>
        <router-link class="button home-hero__cta" to="/contacts?booking=1">
          {{ t('finalCta.iWantToUnderstandMyself') }}
          <span aria-hidden="true">↗</span>
        </router-link>
        <router-link class="home-hero__question" to="/contacts">
          <v-icon icon="mdi-message-outline" size="18" aria-hidden="true" />
          {{ t('home.askAQuestionFirst') }}
          <span aria-hidden="true">↗</span>
        </router-link>
        <p class="home-hero__note">
          <span aria-hidden="true"></span>
          {{ t('home.lessFixingYourselfMoreBeingYourself') }}
        </p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.home-hero {
  position: relative;
  isolation: isolate;
  background: #faf7f2;
}
.home-hero__visual {
  position: absolute;
  inset: 0;
  z-index: -1;
}
.home-hero__photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 40%;
}
.home-hero__visual::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(0deg, #faf7f2 0%, transparent 10%),
    linear-gradient(180deg, #faf7f24d 0%, transparent 12%),
    linear-gradient(90deg, #faf7f2eb 0%, #faf7f2b3 25%, #faf7f24d 44%, transparent 61%);
  pointer-events: none;
}
.home-hero__inner {
  display: flex;
  align-items: center;
  min-height: clamp(660px, 53vw, 820px);
  padding-block: 64px 80px;
}
.home-hero__copy {
  width: 48%;
  max-width: 620px;
}
.home-hero__eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 24px;
  font-size: 9px;
  letter-spacing: 1.7px;
  color: #76583f;
}
.home-hero__star { font-size: 27px; line-height: 1; }
.home-hero h1 {
  font-size: clamp(48px, 4.4vw, 70px);
  line-height: 1.06;
  letter-spacing: -1.5px;
  overflow-wrap: break-word;
}
.home-hero h1 em { color: #89603f; }
.home-hero__description {
  max-width: 385px;
  margin: 25px 0 30px;
  font-size: 14px;
  line-height: 1.85;
  color: #665648;
}
.home-hero__cta { max-width: 100%; }
.home-hero__question {
  display: flex;
  align-items: center;
  gap: 10px;
  width: fit-content;
  min-height: 44px;
  margin-top: 9px;
  font-size: 11px;
  color: #695442;
}
.home-hero__question:hover { color: var(--brown); }
.home-hero__question > span { margin-left: 7px; }
.home-hero__note {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 25px;
  font-family: var(--serif);
  font-size: 19px;
  font-style: italic;
  line-height: 1.4;
  color: #826b55;
}
.home-hero__note > span {
  flex: 0 0 32px;
  height: 1px;
  background: #ac9279;
}
@media (min-width: 1600px) {
  .home-hero__inner { min-height: 820px; }
}
@media (min-width: 601px) and (max-width: 1199px) {
  .home-hero__inner { min-height: 650px; padding-block: 55px 65px; }
  .home-hero__photo { object-position: 67% center; }
  .home-hero__copy { width: 47%; }
  .home-hero h1 { font-size: clamp(40px, 4.5vw, 53px); letter-spacing: -1px; }
  .home-hero__eyebrow { font-size: 8px; letter-spacing: 1px; gap: 7px; }
  .home-hero__description { font-size: 13px; margin-block: 22px 25px; }
  .home-hero__cta { gap: 18px; padding: 16px 20px; font-size: 11px; }
  .home-hero__note { font-size: 17px; }
  .home-hero__visual::after {
    background:
      linear-gradient(0deg, #faf7f2 0%, transparent 9%),
      linear-gradient(90deg, #faf7f2eb 0%, #faf7f2b3 29%, #faf7f24d 46%, transparent 62%);
  }
}
@media (max-width: 600px) {
  .home-hero { display: flex; flex-direction: column; }
  .home-hero__inner {
    order: 0;
    min-height: 0;
    padding-block: 28px 10px;
  }
  .home-hero__copy { width: 100%; max-width: 470px; margin-inline: auto; }
  .home-hero__eyebrow { font-size: 8px; letter-spacing: 1.1px; gap: 8px; margin-bottom: 17px; }
  .home-hero__star { font-size: 24px; }
  .home-hero h1 { font-size: clamp(36px, 10.25vw, 46px); line-height: 1.04; letter-spacing: -1px; }
  .home-hero__description { font-size: 13px; line-height: 1.75; margin: 18px 0 21px; }
  .home-hero__cta { min-height: 52px; gap: 25px; padding: 15px 20px; font-size: 11px; }
  .home-hero__question { margin-top: 5px; font-size: 10px; }
  .home-hero__note { display: none; }
  .home-hero__visual {
    position: relative;
    inset: auto;
    order: 1;
    z-index: 0;
    height: clamp(280px, 81vw, 400px);
  }
  .home-hero__photo { object-position: 90% center; }
  .home-hero__visual::after {
    background: linear-gradient(180deg, #faf7f2 0%, transparent 12%, transparent 90%, #faf7f2 100%);
  }
}
</style>
