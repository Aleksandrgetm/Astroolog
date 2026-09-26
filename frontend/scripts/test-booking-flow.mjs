import {chromium} from 'playwright'
import assert from 'node:assert/strict'
import {mkdir} from 'node:fs/promises'
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8089'
const browser=await chromium.launch({channel:process.env.PLAYWRIGHT_CHANNEL||'chrome'})
try {
 const page=await browser.newPage({reducedMotion:'reduce'})
 const errors=[];page.on('pageerror',e=>errors.push(e.message))
 // Do not open a real messenger or send anything externally.
 await page.addInitScript(()=>{
  window.messengerAttempts=[]
  window.open=()=>({closed:false,opener:null,close(){this.closed=true},location:{replace(url){window.messengerAttempts.push(url)}}})
 })
 const fill=async()=>{
  await page.locator('input[autocomplete="name"]').fill('Booking flow test')
  await page.locator('input[type="email"]').fill('booking-flow@example.invalid')
  await page.locator('textarea').fill('Test — isolated database only')
  await page.locator('.v-checkbox input').check()
 }
 const options=await (await page.request.get(origin+'/api/booking-options')).json()
 assert.ok(!options.some(o=>o.code==='reading-call-40'))
 assert.equal(options.find(o=>o.code==='online-meeting-40').is_addon,true)
 const responses=[]
 page.on('response',async r=>{if(r.url().endsWith('/api/bookings')&&r.status()===201)responses.push(await r.json())})
 for(const width of [390,430,768,1024,1440])for(const lang of ['ru','lv','en']) {
  await page.setViewportSize({width,height:950})
  await page.goto(origin+(lang==='ru'?'':`/${lang}`)+'/contacts?booking=1&option=personality',{waitUntil:'networkidle'})
  await page.locator('.booking-bonus').waitFor()
  assert.equal(await page.locator('input[type="date"]').count(),0)
  assert.equal(await page.locator('input[type="radio"][value="online-meeting-40"]').count(),0)
  assert.match(await page.locator('.booking-total strong').innerText(),/50/)
  await page.locator('.booking-bonus input').check()
  assert.match(await page.locator('.booking-total strong').innerText(),/70/)
  assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false)
  await fill()
  const messenger=lang==='lv'?'telegram':'whatsapp'
  const req=page.waitForRequest(r=>r.url().endsWith('/api/bookings'))
  await page.locator(`[data-messenger="${messenger}"]`).click()
  const request=await req;assert.ok(!('preferred_date' in request.postDataJSON()))
  await page.locator('.success-panel').waitFor()
  assert.equal(await page.locator('.success-messengers').count(),0)
  assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false)
  const attempts=await page.evaluate(()=>window.messengerAttempts)
  assert.equal(attempts.length,1)
  const url=new URL(attempts[0]);assert.equal(url.hostname,messenger==='whatsapp'?'wa.me':'t.me')
  assert.match(url.searchParams.get('text'),/Итого:.*70/)
  assert.ok(!/undefined|null|Желаемая дата/.test(url.searchParams.get('text')))
  if(lang==='ru'&&[390,1440].includes(width)){
   await mkdir('../artifacts/booking-bonus',{recursive:true})
   await page.screenshot({path:`../artifacts/booking-bonus/success-${width}.png`,fullPage:true})
  }
  // Closing success returns to the same payload; choosing another messenger reuses its key and row.
  const firstId=responses.at(-1).booking.id
  await page.locator('.success-panel .v-btn').last().click()
  if(lang==='ru'&&[390,1440].includes(width))await page.screenshot({path:`../artifacts/booking-bonus/form-${width}.png`,fullPage:true})
  const secondResponse=page.waitForResponse(r=>r.url().endsWith('/api/bookings')&&r.status()===201)
  await page.locator('[data-messenger="telegram"]').click()
  assert.equal((await (await secondResponse).json()).booking.id,firstId)
  await page.locator('.success-panel').waitFor()
 }
 console.log('Booking form, success and both messengers passed: 15 locale/viewport combinations')
 // Manual selection without a deep-link must also survive closing success.
 await page.goto(origin+'/contacts?booking=1',{waitUntil:'networkidle'});await fill()
 await page.locator('input[value="personality"]').check()
 await page.locator('.booking-bonus input').check()
 await page.locator('[data-messenger="whatsapp"]').click();await page.locator('.success-panel').waitFor()
 await page.locator('.success-panel .v-btn').last().click()
 assert.equal(await page.locator('input[value="personality"]').isChecked(),true)
 assert.equal(await page.locator('.booking-bonus input').isChecked(),true)
 // Failure before save must not navigate to a messenger and must preserve input.
 await page.goto(origin+'/contacts?booking=1&option=personality',{waitUntil:'networkidle'});await fill()
 await page.route('**/api/bookings',route=>route.fulfill({status:500,json:{code:'server_error'}}))
 await page.locator('[data-messenger="whatsapp"]').click();await page.locator('.form-error').waitFor()
 assert.deepEqual(await page.evaluate(()=>window.messengerAttempts),[])
 assert.equal(await page.locator('input[autocomplete="name"]').inputValue(),'Booking flow test')
 await page.unroute('**/api/bookings')
 // Stale quote refreshes the catalog and requires another explicit click.
 let postCount=0
 await page.route('**/api/bookings',route=>{postCount++;return route.fulfill({status:409,json:{code:'stale_price'}})})
 await page.route('**/api/booking-options',route=>route.fulfill({json:options.map(o=>o.code==='personality'?{...o,price:6100}:o)}))
 await page.locator('[data-messenger="whatsapp"]').click()
 await page.waitForFunction(()=>document.querySelector('.booking-total strong')?.textContent.includes('61'))
 assert.equal(postCount,1);assert.deepEqual(await page.evaluate(()=>window.messengerAttempts),[])
 await page.unroute('**/api/bookings');await page.unroute('**/api/booking-options')
 // Popup failure after a real successful save stays success; retry has no POST.
 await page.goto(origin+'/contacts?booking=1&option=finances',{waitUntil:'networkidle'});await fill()
 await page.evaluate(()=>{window.open=()=>null})
 let saves=0;page.on('request',r=>{if(r.url().endsWith('/api/bookings'))saves++})
 await page.locator('[data-messenger="telegram"]').click();await page.locator('.success-panel').waitFor()
 await page.locator('.messenger-retry').click();assert.equal(saves,1)
 assert.equal(await page.locator('.form-error').count(),0)
 // Explicitly malformed URL navigation also stays saved.
 await page.locator('.success-panel .v-btn').last().click()
 await page.locator('textarea').fill('Distinct payload for navigation failure')
 await page.evaluate(()=>{window.open=()=>({closed:false,close(){},opener:null,location:{replace(){throw new Error('URL error')}}})})
 await page.locator('[data-messenger="whatsapp"]').click();await page.locator('.success-panel').waitFor()
 assert.equal(await page.locator('.form-error').count(),0)
 // Malicious total is ignored; backend returns its own calculation.
 const base=options.find(o=>o.code==='personality')
 const response=await page.request.post(origin+'/api/bookings',{headers:{'Idempotency-Key':crypto.randomUUID()},data:{name:'Total test',email:'total@example.invalid',consent:true,language:'ru',service_id:base.service_id,service_type:base.code,price:base.price,bonus_code:'online-meeting-40',bonus_price:2000,total:1,total_price:1}})
 assert.equal(response.status(),201);assert.equal((await response.json()).booking.total_price,7000)
 assert.deepEqual(errors,[])
 console.log('Errors, stale price, idempotency, no date, popup/URL failure and backend total passed')
}finally{await browser.close()}
