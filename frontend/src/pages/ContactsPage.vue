<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RequestForm from '../components/RequestForm.vue'
const route=useRoute(),router=useRouter(),booking=computed(()=>route.query.booking==='1')
const telegram=import.meta.env.VITE_TELEGRAM_URL,whatsapp=import.meta.env.VITE_WHATSAPP_URL,contactEmail=import.meta.env.VITE_CONTACT_EMAIL
</script>
<template><section class="section container contacts-layout"><div><span class="eyebrow">{{ t('contacts.letsKeepInTouch') }}</span><h1>{{ t('contacts.everyJourney') }}<br>{{ t('contacts.startsWith') }}<br><em>{{ t('contacts.aConversation') }}</em></h1><p class="lead">{{ t('contacts.tellMeWhatMattersToYouRight') }}<br>{{ t('contacts.illHelpYouFindYourNextStep') }}</p><div class="contact-detail"><v-icon icon="mdi-earth"/><div><h3>{{ t('contacts.weMeetOnline') }}</h3><p>{{ t('contacts.fromAnywhereInTheWorldInA') }}</p></div></div><div v-if="telegram||whatsapp||contactEmail" class="contact-links"><a v-if="telegram" :href="telegram" target="_blank" rel="noopener noreferrer">Telegram ↗</a><a v-if="whatsapp" :href="whatsapp" target="_blank" rel="noopener noreferrer">WhatsApp ↗</a><a v-if="contactEmail" :href="`mailto:${contactEmail}`">{{contactEmail}}</a></div><p v-else class="form-hint">{{ t('contacts.directContactDetailsWillBeAddedOnce') }}</p></div><div class="contact-form-panel"><div class="form-tabs" role="group" :aria-label="t('contacts.enquiryType')"><button :class="{active:!booking}" @click="router.replace('/contacts')" :aria-pressed="!booking">{{ t('contacts.askAQuestion') }}</button><button :class="{active:booking}" @click="router.replace('/contacts?booking=1')" :aria-pressed="booking">{{ t('siteLayout.bookASession') }}</button></div><RequestForm :key="String(booking)" :booking="booking"/></div></section></template>
