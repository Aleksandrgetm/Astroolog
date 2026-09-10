import { test } from 'node:test'
import assert from 'node:assert/strict'
import { buildWhatsAppMessage, buildWhatsAppUrl } from '../src/utils/whatsapp.ts'

test('Question preserves user text and omits empty phone', () => {
 const data={name:'Anna & Co',email:'anna@example.com',phone:' ',message:'Jautājums?\nA + B & C # D'}
 const message=buildWhatsAppMessage('contact',data)
 assert.equal(message,'Новый вопрос с сайта\n\nИмя: Anna & Co\nEmail: anna@example.com\n\nВопрос:\nJautājums?\nA + B & C # D')
 const url=new URL(buildWhatsAppUrl(message,'37129580232'))
 assert.equal(url.origin,'https://wa.me');assert.equal(url.pathname,'/37129580232');assert.equal(url.searchParams.get('text'),message)
})
test('Booking includes filled optional fields and never adds empty labels',()=>{
 const data={name:'Елена',email:'test@example.com',phone:'+371 12345678',preferred_date:'2026-12-01',message:'Мой запрос'}
 const full=buildWhatsAppMessage('bookings',data,'Iepazīsti sevi')
 assert.ok(full.includes('Телефон: +371 12345678'));assert.ok(full.includes('Желаемая дата: 2026-12-01'));assert.ok(full.includes('Формат встречи: Iepazīsti sevi'));assert.ok(full.includes('О чём хочется поговорить:\nМой запрос'))
 const minimal=buildWhatsAppMessage('bookings',{name:'Anna',email:'test@example.com',phone:null,message:' ',preferred_date:undefined},'Session')
 for(const unwanted of ['undefined','null','Телефон:','Желаемая дата:','О чём хочется поговорить:'])assert.ok(!minimal.includes(unwanted))
})
test('Invalid configuration never produces a broken WhatsApp link',()=>{
 assert.equal(buildWhatsAppUrl('text',''),'');assert.equal(buildWhatsAppUrl('text','abc'),'')
 assert.ok(buildWhatsAppUrl('text','+371 29 580 232').startsWith('https://wa.me/37129580232?text='))
})
