<script setup lang="ts">
import { content } from '../stores/content'
import { useI18n } from 'vue-i18n'
import { computed, onMounted, reactive, ref, watch, nextTick } from 'vue'
import { useCatalog } from '../stores/catalog'
import { errorKey, submitRequest } from '../services/api'
import { getLocalizedField } from '../i18n/localized'
import { validPhone } from '../utils/validation'
import { formatPrice } from '../utils/prices'
import { buildWhatsAppMessage, buildWhatsAppUrl } from '../utils/whatsapp'
import { buildTelegramUrl } from '../utils/telegram'
import BookingServiceChoice from './BookingServiceChoice.vue'
import MessengerIcon from './MessengerIcon.vue'
import type { RequestInput, BookingOption } from '../types'
const { t, locale } = useI18n()
type Messenger = 'whatsapp' | 'telegram'
const props=defineProps<{booking?:boolean;serviceId?:number}>()
const catalog=useCatalog(),busy=ref(false),validating=ref(false),error=ref(''),success=ref(false)
const selectedOption = ref<BookingOption | null>(null)
const bonusSelected = ref(false)
const bonus = computed(() => selectedOption.value && catalog.options.find(o => o.is_addon && o.service_id === selectedOption.value!.service_id))
const total = computed(() => selectedOption.value?.price == null ? null : selectedOption.value.price + (bonusSelected.value ? bonus.value?.price ?? 0 : 0))
const choiceAttempted = ref(false)
const whatsappUrl=ref(''), telegramUrl=ref('')
const chosenMessenger=ref<Messenger>('whatsapp'), messengerFailed=ref(false)
const successHeading=ref<HTMLElement|null>(null)
const form=ref<{isValid:boolean|null;validate:()=>Promise<{valid:boolean}>}|null>(null)
const data=reactive<RequestInput>({name:'',email:'',phone:'',message:'',consent:false,service_id:props.serviceId})
let attemptKey='', attemptPayload=''
const required=(v:unknown)=>!!String(v??'').trim() || t('requestForm.pleaseFillInThisField')
const email=(v:string)=>/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)||t('requestForm.enterAValidEmailAddress')
const format = computed(() => selectedOption.value ? [getLocalizedField(selectedOption.value, 'format'), ...(bonusSelected.value && bonus.value ? [getLocalizedField(bonus.value, 'format')] : [])].join(' / ') : '')
const price = (value: number | null) => formatPrice(value, locale.value, t('catalog.free'), t('bookingChoice.pricePending'))
onMounted(()=>{if(props.booking)catalog.load()})
watch(() => selectedOption.value?.service_id, () => { bonusSelected.value=false })
watch(bonus, value => { if(!value) bonusSelected.value=false })
function openMessenger(popup: Window | null = null) {
 const url=chosenMessenger.value==='whatsapp'?whatsappUrl.value:telegramUrl.value
 messengerFailed.value=false
 try {
  if (!url) throw new Error('Missing messenger URL')
  if (popup && !popup.closed) { popup.location.replace(url) }
  else { const opened=window.open(url,'_blank'); if (!opened) throw new Error('Popup blocked'); opened.opener=null }
 } catch { popup?.close(); messengerFailed.value=true }
}
async function submit(messenger: Messenger = 'whatsapp') {
 if (busy.value || validating.value || success.value) return
 chosenMessenger.value=messenger
 validating.value=true; error.value=''
 let popup: Window | null=null
 try {
  choiceAttempted.value=true
  const valid=(await form.value?.validate())?.valid
  if(!valid || (props.booking && !selectedOption.value)) return
  busy.value=true
  // Reserve a blank window during the click; no messenger navigation before a successful save.
  if(props.booking) { try { popup=window.open('about:blank','_blank'); if(popup)popup.opener=null } catch { /* retry remains available after save */ } }
  const language=locale.value as 'ru' | 'lv' | 'en'
  const submitted:RequestInput={...data,language,...(props.booking && selectedOption.value ? {
   service_id:selectedOption.value.service_id,service_type:selectedOption.value.code,price:selectedOption.value.price,
   ...(bonusSelected.value && bonus.value ? {bonus_code:bonus.value.code,bonus_price:bonus.value.price}:{}),
  }: {})}
  const fingerprint=JSON.stringify(submitted)
  if(!attemptKey || fingerprint!==attemptPayload){
   attemptPayload=fingerprint; attemptKey=crypto.randomUUID()
   // Persist only a digest and random key, never personal form fields. Survives returning from a messenger.
   try {
    const digest=Array.from(new Uint8Array(await crypto.subtle.digest('SHA-256',new TextEncoder().encode((props.booking?'booking:':'contact:')+fingerprint)))).map(b=>b.toString(16).padStart(2,'0')).join('')
    const storageKey='astroolog-submission-'+digest
    attemptKey=sessionStorage.getItem(storageKey)||attemptKey
    sessionStorage.setItem(storageKey,attemptKey)
   } catch { /* in-memory idempotency also works with storage disabled */ }
  }
  const response=await submitRequest(props.booking?'bookings':'contact',submitted,attemptKey)
  success.value=true
  // Link failures after this point never become submission errors.
  try {
   const saved=response.booking
   const details=saved ? {...submitted,service_title:saved.service_title,option_title:saved.option_title,
    price_label:price(saved.price),bonus_title:saved.bonus_title,bonus_price_label:saved.bonus_price==null?undefined:price(saved.bonus_price),total_label:price(saved.total_price)}:submitted
   const message=buildWhatsAppMessage(props.booking?'bookings':'contact',details,saved?.option_format||'')
   if(props.booking && !saved)throw new Error('Missing saved snapshot')
   whatsappUrl.value=buildWhatsAppUrl(message,content.settings.whatsapp_phone || undefined); telegramUrl.value=buildTelegramUrl(message,content.settings.telegram_url || undefined)
  } catch { whatsappUrl.value='';telegramUrl.value='' }
  if(props.booking)openMessenger(popup)
  await nextTick();successHeading.value?.focus()
 } catch(e) {
  popup?.close()
  if(success.value){messengerFailed.value=true;return}
  error.value=errorKey(e)
  if(error.value==='feedback.stale_price'){await catalog.load({force:true});attemptKey='';attemptPayload=''}
 } finally {busy.value=false;validating.value=false}
}
function closeSuccess() {
 // Preserve booking payload and its key so another messenger does not create another request.
 if(!props.booking){selectedOption.value=null;attemptKey='';attemptPayload='';Object.assign(data,{name:'',email:'',phone:'',message:'',consent:false,service_id:props.serviceId})}
 choiceAttempted.value=false;success.value=false
}
watch(locale,async()=>{await nextTick();if(form.value?.isValid===false)await form.value.validate()})
</script>
<template>
<div v-if="success" class="success-panel" role="status">
  <v-icon icon="mdi-check-circle-outline" size="44" aria-hidden="true"/>
  <h3 ref="successHeading" tabindex="-1">{{ t(booking ? 'bookingFlow.thanks' : 'feedback.thankYou') }}</h3>
  <p>{{ t(booking ? 'bookingFlow.saved' : 'feedback.questionSaved') }}</p>
  <p v-if="!booking && (whatsappUrl || telegramUrl)">{{ t('feedback.whatsappInvitation') }}</p>
  <div v-if="!booking && (whatsappUrl || telegramUrl)" class="success-messengers">
  <v-btn v-if="whatsappUrl" :href="whatsappUrl" target="_blank" rel="noopener noreferrer" color="primary" size="large" class="submit-button whatsapp-continue">
    <v-icon start icon="mdi-whatsapp"/>
    {{ t('feedback.continueWhatsApp') }}
  </v-btn>
  <v-btn v-if="telegramUrl" :href="telegramUrl" target="_blank" rel="noopener noreferrer" color="primary" variant="outlined" size="large" class="submit-button telegram-continue">
    <svg class="telegram-icon" viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false"><path d="m21 3-4 18-6-5-4 3 1-7-6-2 19-7Z"/><path d="m8 12 9-6-6 10"/></svg>
    {{ t('feedback.continueTelegram') }}
  </v-btn>
  </div>
  <template v-if="booking">
    <p>{{ t(messengerFailed ? 'bookingFlow.unavailable' : 'bookingFlow.retryHint') }}</p>
    <v-btn type="button" variant="outlined" color="primary" class="messenger-retry" @click="openMessenger()"><MessengerIcon :messenger="chosenMessenger" />{{ t(chosenMessenger === 'whatsapp' ? 'bookingFlow.retryWhatsApp' : 'bookingFlow.retryTelegram') }}</v-btn>
  </template>
  <v-btn variant="text" class="success-close" @click="closeSuccess">{{ t('feedback.close') }}</v-btn>
</div>
<v-form v-else ref="form" @submit.prevent="submit()" :disabled="busy" class="request-form"><div class="form-grid"><v-text-field v-model="data.name" :label="t('requestForm.yourName')" autocomplete="name" :rules="[required,v=>v.length<=100||t('requestForm.useNoMoreThan100Characters')]" maxlength="100"/><v-text-field v-model="data.email" :label="t('requestForm.emailForMyReply')" type="email" autocomplete="email" :rules="[required,email]" maxlength="254"/></div><v-text-field v-model="data.phone" :label="t('requestForm.phoneOptional')" type="tel" autocomplete="tel" maxlength="30" :rules="[v=>validPhone(v)||t('requestForm.pleaseCheckYourPhoneNumber')]"/><template v-if="booking"><BookingServiceChoice v-model="selectedOption" :service-id="serviceId" :invalid="choiceAttempted && !selectedOption" :disabled="busy || validating"><template #additional="{ serviceId: groupServiceId }"><div v-if="bonus && bonus.service_id === groupServiceId" class="booking-bonus">
 <p class="bonus-label">{{ t('bookingFlow.additional') }}</p>
 <label class="bonus-card" :class="{ selected: bonusSelected }">
  <input v-model="bonusSelected" type="checkbox" :disabled="busy || validating" />
  <span><strong>{{ getLocalizedField(bonus, 'title') }}</strong><span class="bonus-description">{{ t('bookingFlow.bonusDescription') }}</span><span>+{{ price(bonus.price) }}</span></span>
 </label>
</div>
</template></BookingServiceChoice><p v-if="catalog.error" role="alert">{{t(catalog.error)}} <v-btn variant="text" @click="catalog.load">{{ t('requestForm.tryAgain') }}</v-btn></p><div v-if="selectedOption" class="booking-total" aria-live="polite">
 <p>{{ t('bookingFlow.base') }}: {{ getLocalizedField(selectedOption, 'title') }} — {{ price(selectedOption.price) }}</p>
 <p v-if="bonusSelected && bonus">{{ t('bookingFlow.bonus') }}: {{ getLocalizedField(bonus, 'title') }} — +{{ price(bonus.price) }}</p>
 <strong>{{ t('bookingFlow.total') }}: {{ price(total) }}</strong>
</div>
<v-text-field :model-value="format" :label="t('requestForm.sessionFormat')" readonly/></template><v-textarea v-model="data.message" :label="booking?t('requestForm.whatWouldYouLikeToTalkAbout'):t('requestForm.yourQuestion')" :rules="booking?[]:[required]" maxlength="3000" rows="4" counter="3000"/><v-checkbox v-model="data.consent" :rules="[v=>v===true||t('requestForm.yourConsentIsRequiredToSendThis')]"><template #label><span>{{ t('requestForm.iAgreeToThe') }} <LocaleLink to="/privacy" @click.stop>{{ t('requestForm.dataProcessingTerms') }}</LocaleLink></span></template></v-checkbox><p v-if="error" class="form-error" role="alert">{{t(error)}}</p><div v-if="booking" class="form-messengers">
 <p class="form-hint">{{ t('bookingFlow.beforeSend') }}</p>
 <v-btn v-for="messenger in (['whatsapp', 'telegram'] as const)" :key="messenger" type="button" :data-messenger="messenger" color="primary" :variant="messenger === 'whatsapp' ? 'flat' : 'outlined'" :disabled="busy || validating" :loading="(busy || validating) && chosenMessenger === messenger" size="large" class="submit-button" @click="submit(messenger)">
  <MessengerIcon :messenger="messenger" />{{ t((busy || validating) && chosenMessenger === messenger ? 'bookingFlow.saving' : messenger === 'whatsapp' ? 'bookingFlow.sendWhatsApp' : 'bookingFlow.sendTelegram') }}
 </v-btn>
</div><v-btn v-else type="submit" color="primary" :loading="busy || validating" size="large" class="submit-button">{{booking?t('requestForm.sendBookingRequest'):t('requestForm.sendMessage')}} <ArrowIcon class="form-submit-arrow" /></v-btn><p class="form-hint">{{ t('requestForm.yourDetailsAreUsedOnlyToRespond') }}</p></v-form></template>

<style scoped>
.success-messengers { display: grid; gap: 12px; width: 100%; }
.success-messengers .submit-button { width: 100%; min-height: 52px; height: auto; padding-block: 14px; white-space: normal; }
.success-messengers :deep(.v-btn__content) { white-space: normal; line-height: 1.5; }
.telegram-icon { flex-shrink: 0; margin-inline-end: 8px; }

.form-messengers { display:grid;gap:12px; }
.form-messengers .submit-button { width:100%;height:auto;min-height:52px;padding:14px 18px; }
.form-messengers :deep(.v-btn__content) { white-space:normal;line-height:1.5;gap:10px; }
.form-messengers svg { flex-shrink:0; }
.booking-bonus { margin:24px 0; }
.bonus-label { font-size:11px;letter-spacing:1.4px;text-transform:uppercase;margin-bottom:12px; }
.bonus-card { display:flex;align-items:flex-start;gap:12px;border:1px solid var(--line);padding:16px;cursor:pointer; }
.bonus-card.selected { border-color:var(--brown);background:#f5eee5; }
.bonus-card input { accent-color:var(--brown);width:18px;height:18px;flex-shrink:0;margin-top:3px; }
.bonus-card span { min-width:0;line-height:1.6;font-size:13px; }
.bonus-card strong { font-weight:500; }
.bonus-description { display:block;color:var(--muted);margin:8px 0; }
.booking-total { padding:18px 0;margin-bottom:20px;border-block:1px solid var(--line);overflow-wrap:anywhere; }
.booking-total p { margin-bottom:8px;font-size:13px; }
.booking-total strong { font-size:18px;color:var(--brown); }
.messenger-retry { display:flex;width:fit-content;max-width:100%;height:auto;min-height:44px;margin:0 auto 12px;padding:12px 18px;font-size:13px; }
.messenger-retry :deep(.v-btn__content) { gap:10px;white-space:normal;line-height:1.5; }
.messenger-retry svg { flex-shrink:0; }
.success-close { display:flex;width:fit-content;margin-inline:auto; }
@media(max-width:600px) { .messenger-retry { width:100%; } }
</style>
