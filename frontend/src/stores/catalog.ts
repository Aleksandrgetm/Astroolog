import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getServices, getBookingOptions, errorKey } from '../services/api'
import type { Service, BookingOption } from '../types'
export const useCatalog = defineStore('catalog', () => {
 const options=ref<BookingOption[]>([]), services=ref<Service[]>([]), loading=ref(false), error=ref('')
 let loadedAt=0, pending:Promise<void>|null=null
 function load(settings: {force?:boolean}|unknown={}) {
  const force=!!(settings && typeof settings==='object' && 'force' in settings && settings.force)
  if(pending)return pending
  if(!force && loadedAt && Date.now()-loadedAt<60000)return Promise.resolve()
  loading.value=true;error.value=''
  pending=(async()=>{try{[services.value,options.value]=await Promise.all([getServices(),getBookingOptions()]);loadedAt=Date.now()}catch(e){error.value=errorKey(e)}finally{loading.value=false;pending=null}})()
  return pending
 }
 function invalidate(){loadedAt=0}
 return {options,services,loading,error,load,invalidate}
})
