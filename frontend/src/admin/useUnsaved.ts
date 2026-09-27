import { onBeforeUnmount, type Ref } from 'vue'
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router'
export function useUnsaved(dirty:Ref<boolean>) {
  const leave=()=>!dirty.value||window.confirm('У вас есть несохранённые изменения. Покинуть страницу?')
  onBeforeRouteLeave(leave)
  onBeforeRouteUpdate((to,from)=>to.path===from.path||leave())
  const unload=(event:BeforeUnloadEvent)=>{if(dirty.value){event.preventDefault();event.returnValue=''}}
  window.addEventListener('beforeunload',unload)
  onBeforeUnmount(()=>window.removeEventListener('beforeunload',unload))
}
