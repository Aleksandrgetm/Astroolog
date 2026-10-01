// This registry is the sole authority for optional integrations. No trackers are configured.
export const COOKIE_CONSENT_VERSION = 1
export const CONSENT_STORAGE_KEY = 'astroolog-cookie-consent'
export type OptionalCategory = 'functional' | 'analytics' | 'marketing'
export const optionalCategories: OptionalCategory[] = ['functional','analytics','marketing']
export interface Consent { version:number; necessary:true; functional:boolean; analytics:boolean; marketing:boolean; updatedAt:string }
export interface Integration { category:OptionalCategory; start:()=>void; stop:()=>void }
export const integrations: Integration[] = []
export const categoryAvailable=(category:OptionalCategory)=>integrations.some(i=>i.category===category)
export function makeConsent(choice:Partial<Record<OptionalCategory,boolean>>={},available=categoryAvailable):Consent {
 return {version:COOKIE_CONSENT_VERSION,necessary:true,functional:available('functional')&&choice.functional===true,analytics:available('analytics')&&choice.analytics===true,marketing:available('marketing')&&choice.marketing===true,updatedAt:new Date().toISOString()}
}
export function parseConsent(raw:string|null):Consent|null {
 try {
  const c=JSON.parse(raw||'null')
  if(!c||c.version!==COOKIE_CONSENT_VERSION||c.necessary!==true||!optionalCategories.every(k=>typeof c[k]==='boolean')||typeof c.updatedAt!=='string'||!Number.isFinite(Date.parse(c.updatedAt)))return null
  return {...makeConsent(c),updatedAt:c.updatedAt}
 }catch{return null}
}
// Every optional integration needs a stop handler for revocation. Registration alone never starts it.
export function createIntegrationController(registry:Integration[]) {
 const running=new Set<Integration>()
 return (consent:Consent|null)=>{
  for(const item of registry){
   const allowed=consent?.version===COOKIE_CONSENT_VERSION&&consent[item.category]===true
   if(allowed&&!running.has(item)){item.start();running.add(item)}
   else if(!allowed&&running.has(item)){item.stop();running.delete(item)}
  }
 }
}
export const storageInventory = [
 {name:CONSENT_STORAGE_KEY,type:'localStorage',purpose:'choice',duration:'persistent'},
 {name:'astroolog-submission-<SHA-256>',type:'sessionStorage',purpose:'submission',duration:'session'},
 {name:'__Host-astroolog_admin (production) / astroolog_admin (development)',type:'Cookie',purpose:'admin',duration:'hours'},
] as const
