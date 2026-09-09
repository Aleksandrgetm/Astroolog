<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

import { computed, onMounted, reactive, ref, watch, nextTick } from 'vue'
import { useCatalog } from '../stores/catalog'
import { errorKey, submitRequest } from '../services/api'
import { getLocalizedField } from '../i18n/localized'
import type { RequestInput } from '../types'
const props=defineProps<{booking?:boolean;serviceId?:number}>()
const catalog=useCatalog(),busy=ref(false),error=ref(''),success=ref(false)
const form=ref<{isValid:boolean|null;validate:()=>Promise<{valid:boolean}>}|null>(null)
const data=reactive<RequestInput>({name:'',email:'',phone:'',message:'',consent:false,service_id:props.serviceId,preferred_date:''})
const today=new Date().toISOString().slice(0,10)
const required=(v:unknown)=>!!String(v??'').trim() || t('requestForm.pleaseFillInThisField')
const email=(v:string)=>/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)||t('requestForm.enterAValidEmailAddress')
const dateRule=(v:string)=>!v || v>=today || t('requestForm.chooseTodayOrAFutureDate')
const items=computed(()=>catalog.services.map(s=>({title:getLocalizedField(s,'title'),value:s.id})))
onMounted(()=>{if(props.booking)catalog.load()})
async function submit(){if(busy.value)return;error.value='';if(!(await form.value?.validate())?.valid)return;busy.value=true;try{await submitRequest(props.booking?'bookings':'contact',data);success.value=true}catch(e){error.value=errorKey(e)}finally{busy.value=false}}
watch(locale, async()=>{await nextTick();if(form.value?.isValid===false)await form.value.validate()})
</script>
<template><div v-if="success" class="success-panel" role="status"><v-icon icon="mdi-check-circle-outline" size="44"/><h3>{{ t('requestForm.youveTakenTheFirstStep') }}</h3><p>{{t('feedback.success')}}</p><v-btn variant="text" @click="success=false;data.message=''">{{ t('requestForm.sendAnotherMessage') }}</v-btn></div><v-form v-else ref="form" @submit.prevent="submit" :disabled="busy" class="request-form"><div class="form-grid"><v-text-field v-model="data.name" :label="t('requestForm.yourName')" autocomplete="name" :rules="[required,v=>v.length<=100||t('requestForm.useNoMoreThan100Characters')]" maxlength="100"/><v-text-field v-model="data.email" :label="t('requestForm.emailForMyReply')" type="email" autocomplete="email" :rules="[required,email]" maxlength="254"/></div><v-text-field v-model="data.phone" :label="t('requestForm.phoneOptional')" type="tel" autocomplete="tel" maxlength="30" :rules="[v=>!v||/^\+?[0-9 ()-]{6,30}$/.test(v)||t('requestForm.pleaseCheckYourPhoneNumber')]"/><template v-if="booking"><p v-if="catalog.error" role="alert">{{t(catalog.error)}} <v-btn variant="text" @click="catalog.load">{{ t('requestForm.tryAgain') }}</v-btn></p><v-select v-model="data.service_id" :items="items" :label="t('requestForm.sessionFormat')" :loading="catalog.loading" :rules="[required]"/><v-text-field v-model="data.preferred_date" :label="t('requestForm.preferredDateOptional')" type="date" :min="today" :rules="[dateRule]"/><p class="form-hint">{{ t('requestForm.thisIsAPreferredDateOnlyWell') }}</p></template><v-textarea v-model="data.message" :label="booking?t('requestForm.whatWouldYouLikeToTalkAbout'):t('requestForm.yourQuestion')" :rules="booking?[]:[required]" maxlength="3000" rows="4" counter="3000"/><v-checkbox v-model="data.consent" :rules="[v=>v===true||t('requestForm.yourConsentIsRequiredToSendThis')]"><template #label><span>{{ t('requestForm.iAgreeToThe') }} <router-link to="/privacy" @click.stop>{{ t('requestForm.dataProcessingTerms') }}</router-link></span></template></v-checkbox><p v-if="error" class="form-error" role="alert">{{t(error)}}</p><v-btn type="submit" color="primary" :loading="busy" size="large" class="submit-button">{{booking?t('requestForm.sendBookingRequest'):t('requestForm.sendMessage')}} <v-icon end icon="mdi-arrow-top-right"/></v-btn><p class="form-hint">{{ t('requestForm.yourDetailsAreUsedOnlyToRespond') }}</p></v-form></template>
