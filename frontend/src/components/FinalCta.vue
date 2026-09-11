<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
<template><section class="final-cta"><div class="container"><span class="eyebrow">{{ t('finalCta.yourNextChapterStartsWithYou') }}</span><h2>{{ t('finalCta.whatIfNowIs') }}<br>{{ t('finalCta.theRightTimeTo') }} <em>{{ t('finalCta.chooseYourself') }}</em></h2><p>{{ t('finalCta.youDontNeedAllTheAnswers') }}<br>{{ t('finalCta.oneHonestConversationIsEnoughToBegin') }}</p><LocaleLink to="/contacts?booking=1" class="button button-light">{{ t('finalCta.iWantToUnderstandMyself') }} <span><ArrowIcon /></span></LocaleLink><LocaleLink class="quiet-link" to="/contacts">{{ t('finalCta.askAQuestionFirst') }} <ArrowIcon /></LocaleLink></div><span class="cta-orbit" aria-hidden="true"></span></section></template>
