import { defineComponent, Fragment, h } from 'vue'
import { content } from '../stores/content'
export default defineComponent({setup(_props,{slots}){return ()=>{
 const children=slots.default?.()||[]
 const settings=(v:any)=>content.sections.find(s=>s.section_key===v.props?.['data-cms-section'])
 return h(Fragment,children.filter(v=>settings(v)?.is_active!==false).sort((a,b)=>(settings(a)?.sort_order??0)-(settings(b)?.sort_order??0)))
}}})
