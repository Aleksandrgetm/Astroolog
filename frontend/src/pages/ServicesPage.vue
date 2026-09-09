<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
import ServiceCards from '../components/ServiceCards.vue';import FinalCta from '../components/FinalCta.vue'</script>
<template><section class="section container"><div class="page-heading"><span class="eyebrow">{{ t('services.individualSessions') }}</span><h1>{{ t('services.yourQuestion') }}<br><em>{{ t('services.yourFirstStep') }}</em></h1><p>{{ t('services.fromASingleConversationToOngoingSupport') }}<br>{{ t('services.chooseASpaceWhereYoudLikeTo') }}</p></div><ServiceCards/><div class="service-guidance"><h2>{{ t('services.notSure') }} <em>{{ t('services.whatToChoose') }}</em></h2><p>{{ t('services.tellMeALittleAboutYourSituation') }}</p><router-link class="text-link" to="/contacts">{{ t('services.askAQuestion') }}</router-link></div></section><FinalCta/></template>
