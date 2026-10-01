<script setup lang="ts">
import { computed, watch, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { content, translation, translated } from '../stores/content'
import { useCatalog } from '../stores/catalog'
import { legalDocuments, legalPages, technicalCookieSections, legalKey, legalBlocks, type LegalPage, type LegalLocale } from '../legal/content'
import { legalUI } from '../legal/ui'
import { useCookieConsent } from '../legal/useCookieConsent'
import { setPageMeta } from '../i18n/seo'
const route=useRoute(),{locale}=useI18n(),catalog=useCatalog(),{openSettings}=useCookieConsent()
const page=computed(()=>legalPages.includes(route.meta.seo as LegalPage)?route.meta.seo as LegalPage:'privacy'),lang=computed(()=>locale.value as LegalLocale),ui=computed(()=>legalUI[lang.value])
function field(key:string,fallback:string){
 const record=content.sections.find(s=>Object.hasOwn(translation(s,lang.value).fields||{},key))
 return record?translation(record,lang.value).fields[key] as string:fallback
}
const document=computed(()=>legalDocuments[lang.value][page.value])
const title=computed(()=>field(`legal.${page.value}.title`,document.value.title))
const updated=computed(()=>field(`legal.${page.value}.updated`,document.value.updated))
const addon=computed(()=>catalog.options.find(o=>o.code==='online-meeting-40'&&o.is_addon&&o.price!==null))
const addonPrice=computed(()=>addon.value?.price!=null?ui.value.addon+new Intl.NumberFormat(lang.value,{style:'currency',currency:'EUR'}).format(addon.value.price/100):'')
const sections=computed(()=>document.value.sections.map((s,i)=>{const technical=page.value==='cookies'&&technicalCookieSections.includes(i);return {title:technical?s.title:field(legalKey(page.value,i,'title'),s.title),blocks:legalBlocks((technical?s.body:field(legalKey(page.value,i,'body'),s.body)).replaceAll('{addonPrice}',addonPrice.value))}}))
const site=import.meta.env.VITE_SITE_URL?.replace(/\/$/,'')
watch(page, value=>{if(value==='terms')catalog.load()}, {immediate:true})
watchEffect(()=>{if(!legalPages.includes(route.meta.seo as LegalPage))return;const record=content.pages.find(p=>p.key===page.value);setPageMeta(record&&translated(record,'meta_title',lang.value)||title.value,record&&translated(record,'meta_description',lang.value)||ui.value[page.value==='privacy'?'description':page.value==='terms'?'termsDescription':'cookiesDescription'],true,true)})
</script>
<template>
<article class="section container legal-document">
 <h1>{{title}}</h1><p class="legal-updated">{{ui.updated}}: <time :datetime="updated">{{updated}}</time></p>
 <p v-if="site"><a :href="site">{{site}}</a></p>
 <section v-for="(section,index) in sections" :key="`${page}-${index}`">
  <h2>{{index+1}}. {{section.title}}</h2>
  <template v-for="(block,n) in section.blocks" :key="n"><ul v-if="block.type==='ul'"><li v-for="(line,i) in block.lines" :key="i">{{line}}</li></ul><p v-else>{{block.lines.join('\n')}}</p></template>
  <address v-if="(page==='privacy'&&[1,7].includes(index))||(page==='terms'&&[0,3,11].includes(index))||(page==='cookies'&&index===9)"><a href="mailto:jelenabobrovska@gmail.com">jelenabobrovska@gmail.com</a><br>WhatsApp / Telegram: <a href="tel:+37129580232">+371 29 580 232</a> · <a href="https://wa.me/37129580232" rel="noopener noreferrer" target="_blank">WhatsApp</a></address>
  <p v-if="page==='privacy'&&index===7"><a href="https://www.dvi.gov.lv/" rel="noopener noreferrer" target="_blank">Datu valsts inspekcija</a></p>
  <p v-if="page==='privacy'&&index===10"><LocaleLink to="/cookies">{{legalDocuments[lang].cookies.title}}</LocaleLink></p>
  <p v-if="page==='terms'&&index===7"><LocaleLink to="/privacy">{{legalDocuments[lang].privacy.title}}</LocaleLink></p>
  <p v-if="page==='terms'&&index===1"><LocaleLink to="/services">{{ui.services}}</LocaleLink></p>
  <button v-if="page==='cookies'&&index===5" type="button" class="text-link" @click="openSettings">{{ui.settings}}</button>
 </section>
</article>
</template>
<style scoped>
.legal-document { max-width:960px;overflow-wrap:anywhere; }
.legal-document h1 { font-size:clamp(38px,5vw,64px);line-height:1.1; }
.legal-document section { margin-top:36px; }
.legal-document h2 { font-size:clamp(26px,3vw,36px);line-height:1.25;margin-bottom:18px; }
.legal-document p,.legal-document li,.legal-document address { line-height:1.85;white-space:pre-line; }
.legal-document p { margin:16px 0; }
.legal-document ul { padding-left:24px; }
.legal-document li { margin:8px 0; }
.legal-document a { text-decoration:underline;text-underline-offset:4px; }
.legal-document address { font-style:normal; }
.legal-updated { color:var(--muted); }
.legal-document button.text-link { background:none;border:0;border-bottom:1px solid currentColor;padding:0 0 8px;font:inherit;color:inherit;cursor:pointer; }
</style>
