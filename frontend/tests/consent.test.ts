import {test} from 'node:test'
import assert from 'node:assert/strict'
import {COOKIE_CONSENT_VERSION,makeConsent,parseConsent,createIntegrationController,integrations} from '../src/legal/consent.ts'
test('No optional category is enabled by default or when no integration is configured',()=>{
 for(const choice of [{},{functional:true,analytics:true,marketing:true}]) {
  const consent=makeConsent(choice)
  assert.equal(consent.necessary,true)
  assert.equal(consent.functional,false);assert.equal(consent.analytics,false);assert.equal(consent.marketing,false)
 }
 assert.deepEqual(integrations,[])
})
test('Versioned choices survive storage but invalid or old records do not',()=>{
 const consent=makeConsent()
 assert.deepEqual(parseConsent(JSON.stringify(consent)),consent)
 for(const raw of [null,'broken','{}',JSON.stringify({...consent,version:COOKIE_CONSENT_VERSION-1}),JSON.stringify({...consent,analytics:'true'}),JSON.stringify({...consent,updatedAt:'bad'})])assert.equal(parseConsent(raw),null)
})
test('Configured integrations require explicit consent and stop on revocation',()=>{
 const calls:string[]=[]
 const sync=createIntegrationController([{category:'analytics',start:()=>{calls.push('start')},stop:()=>{calls.push('stop')}}])
 sync(null);sync(makeConsent());assert.deepEqual(calls,[])
 const allow=makeConsent({analytics:true},()=>true)
 sync(allow);sync(allow);assert.deepEqual(calls,['start'])
 sync(makeConsent());assert.deepEqual(calls,['start','stop'])
 sync({...allow,version:0});assert.deepEqual(calls,['start','stop'])
})

test('Legal documents have complete RU/LV/EN content and no unresolved source placeholders',async()=>{
 const {readFile}=await import('node:fs/promises')
 const docs=JSON.parse(await readFile(new URL('../src/legal/documents.json',import.meta.url),'utf8'))
 for(const language of ['ru','lv','en'])for(const [page,count] of Object.entries({privacy:13,terms:13,cookies:10})){
  const doc=docs[language][page];assert.ok(doc.title);assert.equal(doc.sections.length,count)
  for(const section of doc.sections){assert.ok(section.title);assert.ok(section.body);assert.ok(!/[\[\]]/.test(section.body))}
 }
})

test('Consent interface translations have the same complete keys',async()=>{
 const {legalUI}=await import('../src/legal/ui.ts')
 for(const language of ['lv','en'] as const){assert.deepEqual(Object.keys(legalUI[language]).sort(),Object.keys(legalUI.ru).sort());assert.ok(Object.values(legalUI[language]).every(value=>value.trim()))}
})
