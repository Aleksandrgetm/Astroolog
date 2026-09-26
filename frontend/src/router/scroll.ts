import { nextTick } from 'vue'
// Await async route content, not a guessed network delay. Bound only the failure case.
export async function anchorElement(hash:string):Promise<HTMLElement|null>{
 await nextTick()
 let id:string;try{id=decodeURIComponent(hash.slice(1))}catch{return null}
 const find=()=>[...document.querySelectorAll<HTMLElement>('[id]')].find(el=>el.id===id&&!el.closest('[inert]'))||null
 const target=find() || await new Promise<HTMLElement|null>(resolve=>{
  const observer=new MutationObserver(()=>{const el=find();if(el){observer.disconnect();clearTimeout(timeout);resolve(el)}})
  const timeout=setTimeout(()=>{observer.disconnect();resolve(null)},10000)
  observer.observe(document.body,{subtree:true,childList:true,attributes:true,attributeFilter:['inert']})
 })
 if(!target)return null
 await document.fonts.ready
 const scope=target.closest('.page-content')
 await Promise.all([...scope?.querySelectorAll('img')||[]].filter(img=>!!(img.compareDocumentPosition(target)&Node.DOCUMENT_POSITION_FOLLOWING)).map(img=>{img.loading='eager';return img.decode().catch(()=>{})}))
 await new Promise<void>(resolve=>requestAnimationFrame(()=>requestAnimationFrame(()=>resolve())))
 return target
}
export async function restoreReady() {
 await nextTick()
 const ready=()=>{const page=document.querySelector('.page-content:not([inert])');return !!page && !page.querySelector('[role="status"]')}
 if(!ready())await new Promise<void>(resolve=>{
  const observer=new MutationObserver(()=>{if(ready()){observer.disconnect();clearTimeout(timeout);resolve()}})
  const timeout=setTimeout(()=>{observer.disconnect();resolve()},10000)
  observer.observe(document.body,{subtree:true,childList:true,attributes:true})
 })
 await document.fonts.ready
 await new Promise<void>(resolve=>requestAnimationFrame(()=>requestAnimationFrame(()=>resolve())))
}
