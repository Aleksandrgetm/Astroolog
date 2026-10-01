import { ref } from 'vue'
import { CONSENT_STORAGE_KEY, integrations, makeConsent, parseConsent, createIntegrationController, type Consent, type OptionalCategory } from './consent'
const choice=ref<Consent|null>(null),ready=ref(false),settingsOpen=ref(false)
const applyIntegrations=createIntegrationController(integrations)
let listening=false
function read(){try{return parseConsent(localStorage.getItem(CONSENT_STORAGE_KEY))}catch{return null}}
function refresh(){choice.value=read();applyIntegrations(choice.value)}
export function useCookieConsent(){
 function initialize(){
  refresh();ready.value=true
  if(!listening){window.addEventListener('storage',e=>{if(e.key===CONSENT_STORAGE_KEY||e.key===null)refresh()});listening=true}
 }
 function save(values:Partial<Record<OptionalCategory,boolean>>={}){
  choice.value=makeConsent(values)
  try{localStorage.setItem(CONSENT_STORAGE_KEY,JSON.stringify(choice.value))}catch{ /* Keep this visit usable when persistence is blocked. */ }
  applyIntegrations(choice.value);settingsOpen.value=false
 }
 return {choice,ready,settingsOpen,initialize,save,accept:()=>save({functional:true,analytics:true,marketing:true}),reject:()=>save(),openSettings:()=>{settingsOpen.value=true}}
}
