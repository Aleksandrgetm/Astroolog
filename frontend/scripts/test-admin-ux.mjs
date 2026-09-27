import {chromium} from 'playwright'
import assert from 'node:assert/strict'
import {mkdir} from 'node:fs/promises'
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8086'
if(!/^http:\/\/(127\.0\.0\.1|localhost):/.test(origin))throw Error('Local isolated server required')
const browser=await chromium.launch({channel:process.env.PLAYWRIGHT_CHANNEL||'chrome'})
const page=await browser.newPage({viewport:{width:1440,height:1000}}),errors=[],writes=[],mediaRequests=[]
page.on('pageerror',e=>errors.push(e.message))
page.on('request',r=>{if(r.url().includes('/api/admin/')&&!['GET','HEAD'].includes(r.method()))writes.push(r.url());if(r.url().includes('/api/admin/media'))mediaRequests.push(r.url())})
const go=async path=>{await page.goto(origin+path,{waitUntil:'networkidle'});await page.locator('h1').first().waitFor()}
const check=async label=>{assert.equal(await page.evaluate(()=>[...document.querySelectorAll('*')].filter(el=>el.tagName.startsWith('V-')).length),0,'Unregistered component '+label);assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+1),false,'Overflow '+label)}
try{
 await go('/admin');await page.waitForURL('**/admin/login');await page.getByLabel('Email',{exact:true}).fill('cms-test@example.invalid');await page.getByLabel('Пароль',{exact:true}).fill('cms-test-only-password-2026');await page.getByRole('button',{name:'Войти',exact:true}).click();await page.waitForURL('**/admin');await page.getByText('Последние заявки',{exact:true}).waitFor()
 const list=async kind=>(await (await page.request.get(origin+'/api/admin/entities/'+kind+'?limit=100')).json()).items
 const ids={};for(const kind of ['services','directions','reviews','certificates','pages'])ids[kind]=(await list(kind))[0].id
 const option=(await list('options'))[0];
 const baseline=await (await page.request.get(origin+'/api/content')).json()
 const routes=['/admin','/admin/content/home','/admin/content/home/hero','/admin/content/home/numerology','/admin/content/home/how-it-works','/admin/services','/admin/services/'+ids.services,'/admin/directions/'+ids.directions,'/admin/reviews/'+ids.reviews,'/admin/certificates/'+ids.certificates,'/admin/seo/'+ids.pages,'/admin/settings','/admin/media','/admin/requests','/admin/audit']
 await mkdir('/tmp/astroolog-admin-ux',{recursive:true})
 for(const width of [390,430,768,1440]){
  await page.setViewportSize({width,height:1000})
  for(const path of routes){await go(path);await check(path+' '+width);assert.equal(await page.locator('.site-header').count(),0);if(path.endsWith('/hero'))await page.screenshot({path:`/tmp/astroolog-admin-ux/hero-${width}.png`,fullPage:true})}
 }
 await page.setViewportSize({width:1440,height:1000});await go('/admin/content/home/hero')
 const bar=await page.locator('.editor-savebar').boundingBox();assert.ok(bar.y+bar.height<=1002,'Save bar must be visible before scrolling');
 const initialWrites=writes.length,initialMedia=mediaRequests.length
 const heading=page.getByLabel('Заголовок — первая строка',{exact:true}),old=await heading.inputValue()
 await heading.fill('Проверка локального предпросмотра')
 assert.ok((await page.locator('.section-draft').innerText()).includes('Проверка локального предпросмотра'))
 assert.equal(writes.length,initialWrites);assert.equal(mediaRequests.length,initialMedia)
 assert.deepEqual(await (await page.request.get(origin+'/api/content')).json(),baseline,'Draft leaked to public API')
 const dismissed=page.waitForEvent('dialog').then(d=>d.dismiss());await page.getByRole('link',{name:'Главная',exact:true}).click();await dismissed;assert.ok(page.url().endsWith('/home/hero'))
 await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor();assert.ok(page.url().endsWith('/home/hero'))
 await go('/admin/content/home/hero');assert.equal(await heading.inputValue(),'Проверка локального предпросмотра');await heading.fill(old);await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor()
 for(const lang of ['LV','EN','RU']){await page.getByRole('tab',{name:new RegExp('^'+lang)}).click();assert.ok((await heading.inputValue()).length>0);assert.equal(await page.getByRole('button',{name:'Сохранить',exact:true}).isDisabled(),true)}
 await page.getByRole('button',{name:'Выбрать из медиатеки',exact:true}).click();await page.getByRole('dialog').waitFor();await page.locator('.media-tile').first().waitFor();assert.ok(mediaRequests.length>initialMedia);await page.locator('.media-tile').first().click();await page.getByRole('button',{name:'Выбрать изображение',exact:true}).waitFor();await page.getByRole('button',{name:'Закрыть медиатеку'}).click();assert.equal(await page.getByRole('button',{name:'Сохранить',exact:true}).isDisabled(),true)
 await go('/admin/certificates/'+ids.certificates);for(const lang of ['RU','LV','EN']){await page.getByRole('tab',{name:new RegExp('^'+lang)}).click();await check('certificate '+lang)}
 await go('/admin/content/about/intro');assert.equal((await page.locator('body').innerText()).includes('{name}'),false);assert.equal(await page.getByRole('button',{name:'Сохранить',exact:true}).isDisabled(),true)
 for(const path of ['/admin/content/home','/admin/services','/admin/media','/admin/seo/'+ids.pages]){await go(path);await page.screenshot({path:'/tmp/astroolog-admin-ux/'+path.split('/').slice(2).join('-')+'.png',fullPage:true})}
 // Exercise the dedicated editors against the isolated database, restoring each value.
 for(const [path,label] of [['/admin/services/'+ids.services,'Название'],['/admin/directions/'+ids.directions,'Название'],['/admin/reviews/'+ids.reviews,'Полная версия для страницы отзывов'],['/admin/certificates/'+ids.certificates,'Название сертификата'],['/admin/certificates/'+ids.certificates,'Год обучения'],['/admin/seo/'+ids.pages,'SEO-заголовок'],['/admin/settings','Профессиональная подпись'],['/admin/options/'+option.code,'Название варианта']]){
  console.log('Save and restore:',path);await go(path);const field=page.getByLabel(label,{exact:true}),before=await field.inputValue();await field.fill(before+' — UX test');await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor();await go(path);assert.equal(await field.inputValue(),before+' — UX test');await field.fill(before);await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor()
 }
 await go('/admin/options/'+option.code);const price=page.getByLabel('Цена, €',{exact:true});const beforePrice=await price.inputValue();await price.fill('25,50');await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor();assert.equal((await (await page.request.get(origin+'/api/admin/entities/options/'+option.code)).json()).price,2550);await price.fill(beforePrice);await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.getByText('Изменения сохранены',{exact:true}).waitFor()
 await go('/admin/reviews/new');await page.getByLabel('Полная версия для страницы отзывов',{exact:true}).fill('Изолированная проверка нового редактора. Этот отзыв остаётся в архиве.');await page.getByRole('button',{name:'Сохранить',exact:true}).click();await page.waitForURL(/\/admin\/reviews\/\d+$/);await page.getByText('Изменения сохранены',{exact:true}).waitFor()
 await go('/admin/media');await page.getByLabel('Найти изображение',{exact:true}).fill('icta');const searchResponse=page.waitForResponse(r=>r.url().includes('/api/admin/media?')&&r.url().includes('q=icta'));await page.getByRole('button',{name:'Найти',exact:true}).click();await searchResponse;assert.ok(await page.locator('.media-tile').count()>0);await page.locator('.media-tile').first().click();assert.equal(await page.getByRole('button',{name:'Удалить',exact:true}).isDisabled(),true)
 await page.setViewportSize({width:390,height:1000});await page.getByRole('button',{name:'Открыть меню',exact:true}).click();await page.getByRole('link',{name:'Контент сайта',exact:true}).click();await page.waitForURL('**/admin/content');await check('mobile drawer navigation')
 await page.setViewportSize({width:1440,height:1000});await page.getByRole('button',{name:'Выйти',exact:true}).click();await page.waitForURL('**/admin/login');assert.equal((await page.request.get(origin+'/api/admin/entities/services/'+ids.services)).status(),401)
 assert.deepEqual(errors,[]);console.log('Admin UX: 60 responsive layouts, local preview isolation, RU/LV/EN, save/reload/restore, dirty leave guard, lazy media picker passed')
}finally{await browser.close()}
