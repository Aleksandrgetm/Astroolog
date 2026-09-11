<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

import { computed, onMounted, reactive, ref, watch, nextTick } from 'vue'
import { useCatalog } from '../stores/catalog'
import { errorKey, submitRequest } from '../services/api'
import { getLocalizedField } from '../i18n/localized'
import { buildWhatsAppMessage, buildWhatsAppUrl } from '../utils/whatsapp'
import { buildTelegramUrl } from '../utils/telegram'
import type { RequestInput } from '../types'
const props=defineProps<{booking?:boolean;serviceId?:number}>()
const catalog=useCatalog(),busy=ref(false),validating=ref(false),error=ref(''),success=ref(false)
const whatsappUrl=ref('')
const telegramUrl=ref('')
const successHeading=ref<HTMLElement|null>(null)
const form=ref<{isValid:boolean|null;validate:()=>Promise<{valid:boolean}>}|null>(null)
const data=reactive<RequestInput>({name:'',email:'',phone:'',message:'',consent:false,service_id:props.serviceId,preferred_date:''})
const today=new Date().toISOString().slice(0,10)
const required=(v:unknown)=>!!String(v??'').trim() || t('requestForm.pleaseFillInThisField')
const email=(v:string)=>/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)||t('requestForm.enterAValidEmailAddress')
const dateRule=(v:string)=>!v || v>=today || t('requestForm.chooseTodayOrAFutureDate')
const items=computed(()=>catalog.services.map(s=>({title:getLocalizedField(s,'title'),value:s.id})))
onMounted(()=>{if(props.booking)catalog.load()})
async function submit() {
  if (busy.value || validating.value || success.value) return
  validating.value = true
  error.value = ''
  try {
    if (!(await form.value?.validate())?.valid) return
    busy.value = true
    const submitted = { ...data }
    const service = catalog.services.find(item => item.id === submitted.service_id)
    const serviceTitle = service ? getLocalizedField(service, 'title') : ''
    await submitRequest(props.booking ? 'bookings' : 'contact', submitted)
    const message = buildWhatsAppMessage(props.booking ? 'bookings' : 'contact', submitted, serviceTitle)
    whatsappUrl.value = buildWhatsAppUrl(message)
    telegramUrl.value = buildTelegramUrl(message)
    success.value = true
    await nextTick()
    successHeading.value?.focus()
  } catch (e) {
    error.value = errorKey(e)
  } finally {
    busy.value = false
    validating.value = false
  }
}
function closeSuccess() {
  Object.assign(data, {name:'',email:'',phone:'',message:'',consent:false,service_id:props.serviceId,preferred_date:''})
  whatsappUrl.value = ''
  telegramUrl.value = ''
  success.value = false
}
watch(locale, async()=>{await nextTick();if(form.value?.isValid===false)await form.value.validate()})
</script>
<template>
<div v-if="success" class="success-panel" role="status">
  <v-icon icon="mdi-check-circle-outline" size="44" aria-hidden="true"/>
  <h3 ref="successHeading" tabindex="-1">{{ t('feedback.thankYou') }}</h3>
  <p>{{ t(booking ? 'feedback.bookingSaved' : 'feedback.questionSaved') }}</p>
  <p v-if="whatsappUrl || telegramUrl">{{ t('feedback.whatsappInvitation') }}</p>
  <div v-if="whatsappUrl || telegramUrl" class="success-messengers">
  <v-btn v-if="whatsappUrl" :href="whatsappUrl" target="_blank" rel="noopener noreferrer" color="primary" size="large" class="submit-button whatsapp-continue">
    <v-icon start icon="mdi-whatsapp"/>
    {{ t('feedback.continueWhatsApp') }}
  </v-btn>
  <v-btn v-if="telegramUrl" :href="telegramUrl" target="_blank" rel="noopener noreferrer" color="primary" variant="outlined" size="large" class="submit-button telegram-continue">
    <svg class="telegram-icon" viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="m21 3-4 18-6-5-4 3 1-7-6-2 19-7Z"/><path d="m8 12 9-6-6 10"/></svg>
    {{ t('feedback.continueTelegram') }}
  </v-btn>
  </div>
  <v-btn variant="text" @click="closeSuccess">{{ t('feedback.close') }}</v-btn>
</div>
<v-form v-else ref="form" @submit.prevent="submit" :disabled="busy" class="request-form"><div class="form-grid"><v-text-field v-model="data.name" :label="t('requestForm.yourName')" autocomplete="name" :rules="[required,v=>v.length<=100||t('requestForm.useNoMoreThan100Characters')]" maxlength="100"/><v-text-field v-model="data.email" :label="t('requestForm.emailForMyReply')" type="email" autocomplete="email" :rules="[required,email]" maxlength="254"/></div><v-text-field v-model="data.phone" :label="t('requestForm.phoneOptional')" type="tel" autocomplete="tel" maxlength="30" :rules="[v=>!v||/^\+?[0-9 ()-]{6,30}$/.test(v)||t('requestForm.pleaseCheckYourPhoneNumber')]"/><template v-if="booking"><p v-if="catalog.error" role="alert">{{t(catalog.error)}} <v-btn variant="text" @click="catalog.load">{{ t('requestForm.tryAgain') }}</v-btn></p><v-select v-model="data.service_id" :items="items" :label="t('requestForm.sessionFormat')" :loading="catalog.loading" :rules="[required]"/><v-text-field v-model="data.preferred_date" :label="t('requestForm.preferredDateOptional')" type="date" :min="today" :rules="[dateRule]"/><p class="form-hint">{{ t('requestForm.thisIsAPreferredDateOnlyWell') }}</p></template><v-textarea v-model="data.message" :label="booking?t('requestForm.whatWouldYouLikeToTalkAbout'):t('requestForm.yourQuestion')" :rules="booking?[]:[required]" maxlength="3000" rows="4" counter="3000"/><v-checkbox v-model="data.consent" :rules="[v=>v===true||t('requestForm.yourConsentIsRequiredToSendThis')]"><template #label><span>{{ t('requestForm.iAgreeToThe') }} <LocaleLink to="/privacy" @click.stop>{{ t('requestForm.dataProcessingTerms') }}</LocaleLink></span></template></v-checkbox><p v-if="error" class="form-error" role="alert">{{t(error)}}</p><v-btn type="submit" color="primary" :loading="busy || validating" size="large" class="submit-button">{{booking?t('requestForm.sendBookingRequest'):t('requestForm.sendMessage')}} <ArrowIcon class="form-submit-arrow" /></v-btn><p class="form-hint">{{ t('requestForm.yourDetailsAreUsedOnlyToRespond') }}</p></v-form></template>

<style scoped>
.success-messengers { display: grid; gap: 12px; width: 100%; }
.success-messengers .submit-button { width: 100%; min-height: 52px; height: auto; padding-block: 14px; white-space: normal; }
.success-messengers :deep(.v-btn__content) { white-space: normal; line-height: 1.5; }
.telegram-icon { flex-shrink: 0; margin-inline-end: 8px; }
</style>
