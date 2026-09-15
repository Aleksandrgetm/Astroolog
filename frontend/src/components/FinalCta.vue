<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
<template><section class="final-cta"><div class="container"><span class="eyebrow"><slot name="eyebrow">{{ t('finalCta.yourNextChapterStartsWithYou') }}</slot></span><h2><slot name="heading">{{ t('finalCta.whatIfNowIs') }}<br>{{ t('finalCta.theRightTimeTo') }} <em>{{ t('finalCta.chooseYourself') }}</em></slot></h2><p><slot name="description">{{ t('finalCta.youDontNeedAllTheAnswers') }}<br>{{ t('finalCta.oneHonestConversationIsEnoughToBegin') }}</slot></p><slot name="actions"><LocaleLink to="/contacts?booking=1" class="button button-light">{{ t('finalCta.iWantToUnderstandMyself') }} <span><ArrowIcon /></span></LocaleLink><LocaleLink class="quiet-link" to="/contacts">{{ t('finalCta.askAQuestionFirst') }} <ArrowIcon /></LocaleLink></slot></div><span class="cta-orbit" aria-hidden="true"></span></section></template>
