// Run after changing bundled legal content. Existing CMS edits are never overwritten.
import {readFile,writeFile} from 'node:fs/promises'
const docs=JSON.parse(await readFile(new URL('../src/legal/documents.json',import.meta.url),'utf8'))
const sections=[]
for(const page of ['privacy','terms','cookies']) {
 const definitions=[{id:'document',fields:['title','updated'],label:docs.ru[page].title},...docs.ru[page].sections.map((s,i)=>({id:'section'+(i+1),fields:['title','body'],label:s.title,index:i}))]
 for(const def of definitions){
  const prefix='legal.'+page+(def.id==='document'?'':'.'+def.id)
  sections.push({section_key:'legal.'+page+'.'+def.id,page,section_type:'text',label:def.label,critical:true,read_only:page==='cookies'&&[1,2,3,4,7].includes(def.index),settings:{},fields:def.fields.map(f=>({key:prefix+'.'+f,label:f==='title'?'Заголовок':f==='updated'?'Дата обновления':'Текст раздела'})),translations:['ru','lv','en'].map(language=>({language,fields:Object.fromEntries(def.fields.map(f=>[prefix+'.'+f,(def.id==='document'?docs[language][page]:docs[language][page].sections[def.index])[f]]))}))})
 }
}
const descriptions={ru:{privacy:'Информация об обработке персональных данных, их защите и ваших правах.',terms:'Условия заказа, оплаты, проведения встреч и возврата средств.',cookies:'Фактическое браузерное хранение и управление выбором cookies.'},lv:{privacy:'Informācija par personas datu apstrādi, aizsardzību un jūsu tiesībām.',terms:'Pasūtīšanas, samaksas, tikšanās un atmaksas noteikumi.',cookies:'Faktiskā datu glabāšana pārlūkā un sīkdatņu izvēles pārvaldība.'},en:{privacy:'Information about personal data processing, protection and your rights.',terms:'Terms for ordering, payment, meetings and refunds.',cookies:'Actual browser storage and management of cookie choices.'}}
const pages=['privacy','cookies','terms'].map(key=>({key,slug:'/'+key,is_active:true,translations:['ru','lv','en'].map(language=>({language,title:docs[language][key].title,meta_title:docs[language][key].title,meta_description:descriptions[language][key]}))}))
const output=JSON.stringify({pages,sections},null,2)+'\n'
const dest=new URL('../../backend/internal/database/legal_seed.json',import.meta.url)
if(process.argv.includes('--check')){if(await readFile(dest,'utf8')!==output)throw Error('Legal seed is out of sync')}else await writeFile(dest,output)
