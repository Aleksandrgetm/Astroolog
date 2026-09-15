<script setup lang="ts">
import { expertImage } from '../utils/expertImage'
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()

import { ref, watch, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { getService, errorKey } from '../services/api'
import type { Service } from '../types'
import ServiceOptions from '../components/ServiceOptions.vue'
import RequestForm from '../components/RequestForm.vue'
import axios from 'axios'
import { getLocalizedField } from '../i18n/localized'
import { setPageMeta } from '../i18n/seo'
const route=useRoute(),service=ref<Service|null>(null),error=ref(''),loading=ref(true),notFound=ref(false)
let request=0
async function load(){const id=++request;loading.value=true;error.value='';service.value=null;notFound.value=false;try{const result=await getService(String(route.params.slug));if(id!==request)return;service.value=result}catch(e){if(id!==request)return;notFound.value=axios.isAxiosError(e)&&e.response?.status===404;error.value=errorKey(e)}finally{if(id===request)loading.value=false}}
watch(()=>route.params.slug,load,{immediate:true})
watchEffect(()=>{ void locale.value; setPageMeta(service.value?getLocalizedField(service.value,'title'):t('seo.serviceTitle'),service.value?getLocalizedField(service.value,'short_description'):t('seo.serviceDescription')) })
</script>
<template><section class="section container"><LocaleLink to="/services" class="back-link"><ArrowIcon direction="left" /> {{ t('service.allSessions') }}</LocaleLink><p v-if="loading" class="state" role="status">{{ t('service.loadingSessionDetails') }}</p><div v-else-if="error" class="state" role="alert"><h1>{{notFound?t('service.sessionNotFound'):t('service.unableToLoadThisSession')}}</h1><p>{{t(error)}}</p><v-btn v-if="!notFound" @click="load">{{ t('requestForm.tryAgain') }}</v-btn><LocaleLink v-else to="/services">{{ t('service.chooseAnotherSession') }}</LocaleLink></div><template v-else-if="service"><div class="service-detail"><div><span class="eyebrow">{{ t('service.yourOwnSpaceForChange') }}</span><h1>{{getLocalizedField(service, 'title')}}</h1><p class="lead">{{getLocalizedField(service, 'short_description')}}</p><span class="duration-chip"><v-icon icon="mdi-clock-outline" size="18"/> {{getLocalizedField(service, 'duration')}}</span><p class="detail-description">{{getLocalizedField(service, 'description')}}</p></div><img :src="expertImage(service.image)" :class="{ 'expert-photo': expertImage(service.image) === '/images/Hero-new.png' }" :alt="getLocalizedField(service, 'title')" width="1536" height="1024"/></div><ServiceOptions :service-id="service.id"/><div id="booking" class="booking-layout"><div><span class="eyebrow">{{ t('service.trustYourFirstStep') }}</span><h2>{{ t('service.letsBegin') }}<br><em>{{ t('service.withYourQuestion') }}</em></h2><p>{{ t('service.sendARequestWellDiscussTheDetails') }}</p></div><RequestForm :key="service.id" booking :service-id="service.id"/></div></template></section></template>

<style scoped>
.detail-description { white-space: pre-line; }
</style>
