import { chromium } from 'playwright'
import assert from 'node:assert/strict'
import { readFile, mkdir } from 'node:fs/promises'
import directions from '../src/data/directions.json' with { type:'json' }
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8087'
const browser=await chromium.launch(process.env.PLAYWRIGHT_CHANNEL?{channel:process.env.PLAYWRIGHT_CHANNEL}:{})
try {
 const page=await browser.newPage({reducedMotion:'reduce',viewport:{width:1440,height:950}})
 const errors=[];page.on('pageerror',e=>errors.push(e.message))
 const apiOptions=await (await page.request.get(origin+'/api/booking-options')).json()
 let checked=0
 for(const width of [390,430,768,1024,1440]) {
  await page.setViewportSize({width,height:950})
  for(const lang of ['ru','lv','en']) {
   const prefix=lang==='ru'?'':`/${lang}`
   await page.goto(origin+prefix+'/',{waitUntil:'networkidle'})
   assert.deepEqual(await page.locator('.direction-card').evaluateAll(cards=>cards.map(a=>a.getAttribute('href'))),directions.map(d=>`${prefix}/directions/${d.slug}`))
   for(const direction of directions) {
    const response=await page.goto(`${origin}${prefix}/directions/${direction.slug}`,{waitUntil:'networkidle'})
    assert.equal(response.status(),200)
    await page.locator('.direction-formats .service-option-row').first().waitFor()
    assert.equal(await page.locator('h1').count(),1)
    assert.equal(await page.locator('html').getAttribute('lang'),lang)
    assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false,`${width} ${lang} ${direction.slug}`)
    assert.equal((await page.locator('body').innerText()).includes('siteLayout.services'),false)
    assert.deepEqual(await page.locator('.direction-formats .service-option-row').evaluateAll(rows=>rows.map(r=>r.id)),direction.relatedBookingOptionCodes)
    for(const code of direction.relatedBookingOptionCodes){assert.ok(apiOptions.find(o=>o.code===code));assert.equal(await page.locator(`#${code} a`).getAttribute('href'),`${prefix}/contacts?booking=1&option=${code}`)}
    assert.equal(await page.locator('.final-cta .button').getAttribute('href'),`${prefix}/contacts?booking=1&option=${direction.relatedBookingOptionCodes[0]}`)
    assert.equal(await page.locator('.breadcrumbs a').nth(1).getAttribute('href'),`${prefix}/#directions`)
    assert.ok((await page.locator('meta[property="og:image"]').getAttribute('content')).endsWith(direction.image))
    await page.locator('.direction-hero img').evaluate(img=>img.decode())
    if(lang==='ru'&&[390,1440].includes(width)) {await mkdir('../artifacts/directions',{recursive:true});await page.screenshot({path:`../artifacts/directions/${direction.slug}-${width}.png`,fullPage:true})}
    checked++
   }
  }
  console.log(`Directions responsive passed at ${width}`)
 }
 // All cards navigate; each booking action selects the right API variant in every locale.
 for(const lang of ['ru','lv','en']) {
  const prefix=lang==='ru'?'':`/${lang}`
  for(const direction of directions) {
   await page.goto(origin+prefix+'/',{waitUntil:'networkidle'})
   await page.locator(`.direction-card[href="${prefix}/directions/${direction.slug}"]`).click()
   await page.waitForURL(`**${prefix}/directions/${direction.slug}`)
   await page.locator('.direction-formats .service-option-row').first().waitFor()
   for(const code of direction.relatedBookingOptionCodes) {
    await page.locator(`#${code} a`).click();await page.waitForURL(`**${prefix}/contacts?booking=1&option=${code}`)
    await page.locator(`input[value="${code}"]`).waitFor();assert.equal(await page.locator(`input[value="${code}"]`).isChecked(),true)
    await page.goBack();await page.waitForURL(`**${prefix}/directions/${direction.slug}`);await page.locator(`#${code}`).waitFor()
   }
   await page.locator('.breadcrumbs a').nth(1).click();await page.waitForURL(`**${prefix}/#directions`)
   await page.waitForFunction(()=>{const el=document.getElementById('directions');return el && Math.abs(el.getBoundingClientRect().top-document.querySelector('header').getBoundingClientRect().height-20)<5})
  }
  await page.goto(origin+prefix+'/services/personal-matrix',{waitUntil:'networkidle'})
  await page.locator('.breadcrumbs').waitFor()
  assert.equal((await page.locator('.breadcrumbs').innerText()).includes('siteLayout.services'),false)
  assert.equal(await page.locator('.breadcrumbs a').nth(1).innerText(),{ru:'Услуги',lv:'Pakalpojumi',en:'Services'}[lang])
 }
 // A new API price must be rendered, with no editorial price overriding it.
 await page.route('**/api/booking-options',async route=>{await route.fulfill({json:apiOptions.map(o=>o.code==='personality'?{...o,price:6100}:o)})})
 await page.goto(origin+'/directions/potential',{waitUntil:'networkidle'})
 assert.match(await page.locator('#personality .option-price').innerText(),/61/)
 await page.unroute('**/api/booking-options')
 // Missing/inactive variants and network errors leave useful editorial content and safe CTAs.
 await page.route('**/api/booking-options',route=>route.fulfill({json:[]}))
 await page.reload({waitUntil:'networkidle'})
 assert.equal(await page.locator('.direction-formats .service-option-row').count(),0)
 assert.equal(await page.locator('.final-cta .button').getAttribute('href'),'/contacts?booking=1')
 await page.unroute('**/api/booking-options')
 await page.route('**/api/services',route=>route.fulfill({status:503,json:{code:'server_error'}}))
 await page.reload({waitUntil:'networkidle'})
 assert.equal(await page.locator('h1').count(),1);await page.locator('.direction-formats [role="alert"]').waitFor()
 await page.unroute('**/api/services')
 const manifest=JSON.parse(await readFile('dist/routes.json','utf8'))
 assert.equal(manifest.paths.filter(p=>p.includes('/directions/')).length,12)
 assert.deepEqual(errors,[])
 console.log(`${checked} direction layouts, 12 card flows, 24 booking links, API prices/fallbacks, breadcrumbs and Back passed.`)
} finally {await browser.close()}
