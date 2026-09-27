<script setup lang="ts">
import { content } from '../stores/content'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RequestForm from '../components/RequestForm.vue'
const route=useRoute(),router=useRouter(),booking=computed(()=>route.query.booking==='1')
function changeMode(isBooking:boolean){
 const query={...route.query}
 if(isBooking)query.booking='1';else delete query.booking
 router.replace({path:route.path,query,hash:route.hash})
}
</script>
<template><section class="section container contacts-layout"><div><span class="eyebrow">{{ t('contacts.letsKeepInTouch') }}</span><h1>{{ t('contacts.everyJourney') }}<br>{{ t('contacts.startsWith') }}<br><em>{{ t('contacts.aConversation') }}</em></h1><p class="lead">{{ t('contacts.tellMeWhatMattersToYouRight') }}<br>{{ t('contacts.illHelpYouFindYourNextStep') }}</p><div class="contact-detail"><v-icon icon="mdi-earth"/><div><h2>{{ t('contacts.weMeetOnline') }}</h2><p>{{ t('contacts.fromAnywhereInTheWorldInA') }}</p></div></div><div class="owner-contacts">
  <div class="owner-contact">
    <v-icon icon="mdi-phone-outline" aria-hidden="true" />
    <div class="owner-contact-copy">
      <p class="owner-contact-label">{{ t('contacts.phoneLabel') }}</p>
      <a class="owner-contact-value" :href="'tel:'+(content.settings.phone || '+371 29 580 232').replace(/[^+0-9]/g,'')">{{ content.settings.phone || '+371 29 580 232' }}</a>
      <a class="owner-whatsapp" :href="'https://wa.me/'+(content.settings.whatsapp_phone || '37129580232')" target="_blank" rel="noopener noreferrer">{{ t('contacts.messageWhatsApp') }}</a>
    </div>
  </div>
  <div class="owner-contact">
    <v-icon icon="mdi-email-outline" aria-hidden="true" />
    <div class="owner-contact-copy">
      <p class="owner-contact-label">{{ t('contacts.emailLabel') }}</p>
      <a class="owner-contact-value" :href="'mailto:'+(content.settings.email || 'jelenabobrovska@gmail.com')">{{ content.settings.email || 'jelenabobrovska@gmail.com' }}</a>
    </div>
  </div>
</div>
</div><div class="contact-form-panel"><div class="form-tabs" role="group" :aria-label="t('contacts.enquiryType')"><button :class="{active:!booking}" @click="changeMode(false)" :aria-pressed="!booking">{{ t('contacts.askAQuestion') }}</button><button :class="{active:booking}" @click="changeMode(true)" :aria-pressed="booking">{{ t('siteLayout.bookASession') }}</button></div><RequestForm :key="String(booking)" :booking="booking"/></div></section></template>

<style scoped>
.owner-contacts { display: grid; gap: 28px; margin-top: 38px; }
.owner-contact { display: flex; align-items: flex-start; gap: 18px; }
.owner-contact > .v-icon { color: #a98c6b; font-size: 24px; margin-top: 2px; flex-shrink: 0; }
.owner-contact-copy { min-width: 0; }
.owner-contact-label { margin: 0 0 8px; font-size: 10px; line-height: 1.5; letter-spacing: 1.7px; text-transform: uppercase; color: #826b55; }
.owner-contact-value { display: inline-block; font-size: 16px; line-height: 1.65; font-weight: 400; color: var(--ink); overflow-wrap: anywhere; }
.owner-whatsapp { display: block; width: fit-content; margin-top: 7px; padding-block: 5px; font-size: 11px; line-height: 1.6; color: var(--brown); text-decoration: underline; text-underline-offset: 4px; text-decoration-color: #baa38b; }
@media (hover: hover) and (pointer: fine) {
  .owner-contact a { transition: color 250ms ease; }
  .owner-contact a:hover { color: #a16c3d; }
}
@media (max-width: 600px) {
  .owner-contacts { margin-top: 30px; gap: 24px; }
  .owner-contact-value { font-size: 15px; }
}
</style>
