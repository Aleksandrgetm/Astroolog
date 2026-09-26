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
 assert.ok(full.includes('Телефон: +371 12345678'));assert.ok(!full.includes('Желаемая дата:'));assert.ok(full.includes('Формат встречи: Iepazīsti sevi'));assert.ok(full.includes('О чём хочется поговорить:\nМой запрос'))
 const minimal=buildWhatsAppMessage('bookings',{name:'Anna',email:'test@example.com',phone:null,message:' '},'Session')
 for(const unwanted of ['undefined','null','Телефон:','Желаемая дата:','О чём хочется поговорить:'])assert.ok(!minimal.includes(unwanted))
})
test('Invalid configuration never produces a broken WhatsApp link',()=>{
 assert.equal(buildWhatsAppUrl('text',''),'');assert.equal(buildWhatsAppUrl('text','abc'),'')
 assert.ok(buildWhatsAppUrl('text','+371 29 580 232').startsWith('https://wa.me/37129580232?text='))
})

test('Booking includes base, optional bonus and server total',()=>{
 const details={name:'Anna',email:'a@example.com',service_title:'Моя матрица возможностей',option_title:'Личность',price_label:'€50',bonus_title:'Онлайн-встреча 40 минут',bonus_price_label:'€20',total_label:'€70'}
 const text=buildWhatsAppMessage('bookings',details,'Онлайн')
 for(const field of ['Разбор: Личность','Стоимость разбора: €50','Стоимость бонуса: €20','Итого: €70'])assert.ok(text.includes(field))
 const without=buildWhatsAppMessage('bookings',{...details,bonus_title:null,bonus_price_label:undefined,total_label:'€50'})
 assert.ok(without.includes('Бонус: Не выбран'));assert.ok(!without.includes('Стоимость бонуса:'))
})
