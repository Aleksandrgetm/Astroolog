<script setup lang="ts">
import { computed, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { VDialog, VCard } from 'vuetify/components'
import { useCookieConsent } from '../legal/useCookieConsent'
import { optionalCategories, categoryAvailable } from '../legal/consent'
import { legalUI } from '../legal/ui'
import { legalDocuments, type LegalLocale } from '../legal/content'
const {locale}=useI18n(),ui=computed(()=>legalUI[locale.value as LegalLocale])
const {choice,ready,settingsOpen,initialize,save,accept,reject}=useCookieConsent()
const draft=reactive({functional:false,analytics:false,marketing:false})
watch(settingsOpen,open=>{if(open)Object.assign(draft,{functional:false,analytics:false,marketing:false},choice.value||{})})
onMounted(initialize)
</script>
<template>
<aside v-if="ready&&!choice&&!settingsOpen" class="cookie-banner" data-cookie-banner :aria-label="legalDocuments[locale as LegalLocale].cookies.title">
 <p>{{ui.banner}} <LocaleLink to="/cookies">{{legalDocuments[locale as LegalLocale].cookies.title}}</LocaleLink></p>
 <div class="cookie-actions"><button type="button" @click="accept">{{ui.accept}}</button><button type="button" @click="reject">{{ui.reject}}</button><button type="button" @click="settingsOpen=true">{{ui.configure}}</button></div>
</aside>
<v-dialog v-model="settingsOpen" max-width="600" scrollable aria-labelledby="cookie-settings-title">
 <v-card class="cookie-settings" data-cookie-settings>
  <h2 id="cookie-settings-title">{{ui.settings}}</h2>
  <div class="cookie-options">
   <div class="cookie-category"><label><input type="checkbox" checked disabled> {{ui.necessary}} — {{ui.always}}</label><p>{{ui.necessaryText}}</p></div>
   <div v-for="category in optionalCategories" :key="category" class="cookie-category"><label><input v-model="draft[category]" type="checkbox" :disabled="!categoryAvailable(category)"> {{ui[category]}}<span v-if="!categoryAvailable(category)"> — {{ui.unused}}</span></label><p>{{ui[`${category}Text`]}}</p></div>
  </div>
  <div class="cookie-actions"><button type="button" @click="save(draft)">{{ui.save}}</button><button type="button" @click="settingsOpen=false">{{ui.close}}</button></div>
 </v-card>
</v-dialog>
</template>
<style scoped>
.cookie-banner { position:fixed;inset:auto 16px 16px;z-index:1005;max-width:920px;margin:auto;padding:20px 24px;background:var(--cream,#faf7f2);color:var(--text,#443c33);border:1px solid var(--line,#d9cabb);box-shadow:0 4px 24px #443c3320; }
.cookie-banner p { margin:0 0 12px;font-size:14px;line-height:1.6; }
.cookie-banner a { text-decoration:underline; }
.cookie-actions { display:flex;gap:12px;flex-wrap:wrap; }
.cookie-actions button { background:transparent;border:1px solid currentColor;color:inherit;padding:10px 16px;min-height:44px;font:inherit;cursor:pointer; }
.cookie-actions button:focus-visible { outline:2px solid currentColor;outline-offset:3px; }
.cookie-settings { padding:24px!important;background:#faf7f2!important;color:#443c33!important;max-height:calc(100dvh - 48px); }
.cookie-settings h2 { font-size:32px;line-height:1.2;margin-bottom:16px; }
.cookie-options { overflow-y:auto;min-height:0; }
.cookie-category { padding:16px 0;border-top:1px solid #d9cabb; }
.cookie-category label { display:block;font-weight:600; }
.cookie-category input { margin-right:8px;accent-color:#795538; }
.cookie-category p { font-size:14px;line-height:1.6;margin-top:8px; }
.cookie-settings .cookie-actions { flex-shrink:0;padding-top:16px; }
@media(max-width:600px){.cookie-banner{inset:auto 10px 10px;padding:16px;}.cookie-actions{gap:8px;}.cookie-actions button{padding:10px 12px;font-size:14px;}.cookie-settings{padding:20px!important;}}
</style>
