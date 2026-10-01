import {chromium} from 'playwright'
import assert from 'node:assert/strict'
import {readFile,mkdir} from 'node:fs/promises'
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8086'
if(!/^http:\/\/(127\.0\.0\.1|localhost):/.test(origin))throw Error('Local isolated test server required')
const docs=JSON.parse(await readFile('src/legal/documents.json','utf8'))
const browser=await chromium.launch({channel:process.env.PLAYWRIGHT_CHANNEL||'chrome'})
try {
 const context=await browser.newContext({reducedMotion:'reduce'})
 const page=await context.newPage(),errors=[],tracking=[]
 page.on('pageerror',e=>errors.push(e.message))
 page.on('request',r=>{if(/google-analytics|googletagmanager|connect.facebook.net|facebook.com\/tr/.test(r.url()))tracking.push(r.url())})
 await page.goto(origin+'/cookies',{waitUntil:'networkidle'})
 await page.locator('[data-cookie-banner]').waitFor()
 assert.deepEqual(await context.cookies(),[],'Ordinary visitor cookies')
 const stored=()=>page.evaluate(()=>JSON.parse(localStorage.getItem('astroolog-cookie-consent')))
 assert.equal(await stored(),null)
 await page.getByRole('button',{name:'Принять',exact:true}).click()
 assert.deepEqual(Object.fromEntries(Object.entries(await stored()).filter(([k])=>['functional','analytics','marketing','necessary'].includes(k))),{necessary:true,functional:false,analytics:false,marketing:false})
 await page.reload({waitUntil:'networkidle'});assert.equal(await page.locator('[data-cookie-banner]').count(),0)
 await page.locator('.footer-legal button').click()
 const dialog=page.locator('[data-cookie-settings]');await dialog.waitFor()
 assert.equal(await dialog.locator('input:disabled').count(),4)
 assert.equal(await dialog.locator('input:checked').count(),1)
 await page.keyboard.press('Tab');assert.ok(await page.evaluate(()=>!!document.activeElement.closest('[role=dialog]')))
 await page.keyboard.press('Escape');await dialog.waitFor({state:'hidden'})
 await page.locator('.footer-legal button').click();await dialog.getByRole('button',{name:'Сохранить настройки'}).click()
 await page.evaluate(()=>{const c=JSON.parse(localStorage.getItem('astroolog-cookie-consent'));c.version=0;localStorage.setItem('astroolog-cookie-consent',JSON.stringify(c))})
 await page.reload({waitUntil:'networkidle'});await page.locator('[data-cookie-banner]').waitFor()
 await page.getByRole('button',{name:'Отклонить',exact:true}).click();assert.equal((await stored()).analytics,false)
 await page.goto(origin+'/contacts',{waitUntil:'networkidle'})
 const consent=page.locator('.v-checkbox input');assert.equal(await consent.isChecked(),false)
 assert.equal(await page.locator('.v-checkbox a').getAttribute('href'),'/privacy')
 // No backend write can occur without privacy consent; intercept to guard the isolated DB too.
 let submitted=0
 await page.route('**/api/contact',async r=>{submitted++;await r.fulfill({status:400,body:'{}'})})
 await page.locator('input[autocomplete=name]').fill('Legal test')
 await page.locator('input[type=email]').fill('legal@example.invalid')
 await page.locator('textarea').fill('Consent validation')
 await page.locator('form button[type=submit]').click()
 await page.waitForTimeout(250);assert.equal(submitted,0)
 assert.equal(await page.locator('input[autocomplete=name]').isEnabled(),true)
 await mkdir('/tmp/astroolog-legal',{recursive:true})
 for(const width of [390,430,768,1024,1440])for(const lang of ['ru','lv','en'])for(const key of ['privacy','terms','cookies']) {
  await page.setViewportSize({width,height:900})
  const path=(lang==='ru'?'':`/${lang}`)+'/'+key
  await page.goto(origin+path,{waitUntil:'networkidle'})
  assert.equal(await page.locator('h1').innerText(),docs[lang][key].title)
  assert.equal(await page.locator('html').getAttribute('lang'),lang)
  assert.equal(await page.locator('.legal-document section').count(),docs[lang][key].sections.length)
  assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false,path+' '+width)
  assert.ok(!/\[(?:адрес|номер|ссылка|дата|e-mail)|\{addonPrice\}/.test(await page.locator('.legal-document').innerText()))
  assert.match(await page.locator('meta[name=robots]').getAttribute('content'),/noindex/)
  assert.equal(await page.locator('link[hreflang]').count(),4)
  for(const target of ['privacy','terms','cookies'])assert.equal(await page.locator(`.footer-legal a[href="${lang==='ru'?'':'/'+lang}/${target}"]`).count(),1)
  if(width===390&&lang==='ru')await page.screenshot({path:`/tmp/astroolog-legal/${key}-390.png`,fullPage:true})
 }
 // SPA navigation must not leave stale legal metadata or trigger unmounted-page errors.
 await page.locator('.footer-legal a[href$="/terms"]').click();await page.waitForURL('**/terms')
 await page.locator('.footer-top a[href$="/about"]').click();await page.waitForURL('**/about')
 await page.waitForTimeout(100)
 for(const width of [390,430,768,1024,1440])for(const language of ['ru','lv','en']) {
  await page.setViewportSize({width,height:900})
  await page.evaluate(()=>localStorage.removeItem('astroolog-cookie-consent'))
  await page.goto(origin+(language==='ru'?'':`/${language}`)+'/cookies',{waitUntil:'networkidle'})
  const banner=page.locator('[data-cookie-banner]');await banner.waitFor()
  assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false)
  assert.ok((await banner.boundingBox()).height<450)
  if(width===390)await page.screenshot({path:`/tmp/astroolog-legal/banner-${language}.png`})
  await banner.locator('button').nth(2).click()
  await page.locator('[data-cookie-settings]').waitFor()
  assert.ok((await page.locator('[data-cookie-settings]').boundingBox()).height<=900)
  for(let i=0;i<6;i++){await page.keyboard.press('Tab');assert.ok(await page.evaluate(()=>!!document.activeElement.closest('[role=dialog]')))}
  if(width===390)await page.screenshot({path:`/tmp/astroolog-legal/settings-${language}.png`})
  await page.locator('[data-cookie-settings] button').first().click()
 }
 assert.deepEqual(tracking,[])
 assert.deepEqual(errors,[])
 // Built HTML has complete legal content but no initial consent banner.
 for(const lang of ['ru','lv','en'])for(const key of ['privacy','terms','cookies']){
  const html=await readFile(`dist${lang==='ru'?'':'/'+lang}/${key}/index.html`,'utf8')
  assert.ok(html.includes(docs[lang][key].sections.at(-1).title))
  assert.ok(!html.includes('data-cookie-banner'))
 }
 const admin=await browser.newPage()
 admin.on('pageerror',e=>errors.push(e.message))
 const save=async(expected=200)=>{const result=admin.waitForResponse(r=>r.request().method()==='PUT'&&r.url().includes('/api/admin/entities/sections/'));await admin.getByRole('button',{name:'Сохранить',exact:true}).click();assert.equal((await result).status(),expected)}
 await admin.goto(origin+'/admin/login')
 await admin.getByLabel('Email',{exact:true}).fill('cms-test@example.invalid')
 await admin.getByLabel('Пароль',{exact:true}).fill('cms-test-only-password-2026')
 await admin.getByRole('button',{name:'Войти',exact:true}).click();await admin.waitForURL('**/admin')
 for(const key of ['privacy','cookies','terms']) {
  await admin.goto(origin+`/admin/content/${key}/section1`,{waitUntil:'networkidle'})
  for(const language of ['ru','lv','en']) {
   await admin.getByRole('tab',{name:new RegExp('^'+language.toUpperCase()+' —')}).click()
   const body=admin.getByLabel('Текст раздела',{exact:true}),previous=await body.inputValue()
   if(key==='privacy'&&language==='ru'){await body.fill(previous+'<script>window.legalInjection=true</script>');await save(400)}
   const text=previous+'\n\nLegal CMS test\nSecond line'
   await body.fill(text)
   await save()
   await admin.getByText('Изменения сохранены',{exact:true}).waitFor()
   await page.goto(origin+(language==='ru'?'':`/${language}`)+'/'+key,{waitUntil:'networkidle'})
   assert.ok((await page.locator('.legal-document').innerText()).includes('Legal CMS test'))
   assert.equal(await page.evaluate(()=>window.legalInjection),undefined)
   await body.fill(previous);await save()
   await admin.getByText('Изменения сохранены',{exact:true}).waitFor()
  }
 }
 await admin.goto(origin+'/admin/content/cookies/section2',{waitUntil:'networkidle'})
 assert.equal(await admin.getByLabel('Текст раздела',{exact:true}).getAttribute('readonly'),'')
 await admin.close()
 assert.deepEqual(errors,[])
 console.log('Legal: 45 responsive/locale routes; consent accept/reject/save/version/reload, keyboard dialog, forms, tracking and prerender passed')
}finally{await browser.close()}
