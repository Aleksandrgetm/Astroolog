import { reactive } from 'vue'
import { i18n } from '../i18n'
export interface Translation { language: string; [key: string]: any }
export interface ContentRecord { id:number; translations:Translation[]; [key:string]:any }
export const content=reactive<{ready:boolean;pages:ContentRecord[];sections:ContentRecord[];steps:ContentRecord[];directions:ContentRecord[];reviews:ContentRecord[];certificates:ContentRecord[];settings:Record<string,string>;media_variants:Record<string,Record<string,string>>}>({ready:false,pages:[],sections:[],steps:[],directions:[],reviews:[],certificates:[],settings:{},media_variants:{}})
export function translation(record:ContentRecord,locale:string):Translation{return record.translations.find(t=>t.language===locale)||record.translations.find(t=>t.language==='ru')||{language:locale}}
export function translated(record:ContentRecord,key:string,locale:string):string{return translation(record,locale)[key]||translation(record,'ru')[key]||''}
function assign(target:Record<string,any>,path:string,value:any){const parts=path.split('.');let node=target;for(const part of parts.slice(0,-1))node=node[part]||=( {});node[parts.at(-1)!]=value}
export async function loadContent(){
 try{
  const response=await fetch((import.meta.env.VITE_API_URL||'/api')+'/content',{signal:AbortSignal.timeout(10000)});if(!response.ok)throw Error('Content unavailable')
  const data=await response.json();Object.assign(content,data,{ready:true})
  for(const lang of ['ru','lv','en'] as const){
   const messages:Record<string,any>={}
   for(const section of content.sections){const ru=translation(section,'ru').fields||{},local=translation(section,lang).fields||{};for(const key of new Set([...Object.keys(ru),...Object.keys(local)]))assign(messages,key,local[key]||ru[key]||'')}
   for(const page of content.pages){assign(messages,`seo.${page.key}Title`,translated(page,'meta_title',lang));assign(messages,`seo.${page.key}Description`,translated(page,'meta_description',lang))}
   for(const d of content.directions){const tr=translation(d,lang);const key=d.content_key||d.slug;for(const [field,mapped]of Object.entries({title:'title',eyebrow:'eyebrow',lead:'lead',context:'context',format_note:'formatNote',seo_title:'seoTitle',seo_description:'seoDescription',alt:'imageAlt'}))assign(messages,`directions.${key}.${mapped}`,tr[field]||translation(d,'ru')[field]||'')}
   if(content.settings.expert_name){assign(messages,'siteLayout.expertName',content.settings.expert_name);assign(messages,'siteLayout.expertNameUppercase',content.settings.expert_name.toUpperCase())}
   if(content.settings['professional_title_'+lang])assign(messages,'siteLayout.numerologistCoach',content.settings['professional_title_'+lang])
   i18n.global.mergeLocaleMessage(lang,messages)
  }
 }catch{ /* Existing bundled content remains available when the API is temporarily unavailable. */ }
}
export function sectionImage(key:string,fallback:string){return content.sections.find(s=>s.section_key===key)?.settings?.image||fallback}
const destinations:Record<string,string>={booking:'/contacts?booking=1',question:'/contacts',services:'/services',reviews:'/reviews',about:'/about'}
export function ctaDestination(key:string,which='primary_destination',fallback='booking'){return destinations[content.sections.find(s=>s.section_key===key)?.settings?.[which]||fallback]||destinations[fallback]!}

export function imageSrcset(path:string){const variants=content.media_variants[path];return variants?Object.entries(variants).map(([width,url])=>`${url} ${width}w`).join(', '):undefined}
