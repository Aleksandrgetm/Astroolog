import { chromium } from 'playwright'
import { readFile } from 'node:fs/promises'
import assert from 'node:assert/strict'
const manifest=JSON.parse(await readFile('dist/routes.json','utf8'))
const browser=await chromium.launch(process.env.PLAYWRIGHT_CHANNEL?{channel:process.env.PLAYWRIGHT_CHANNEL}:{})
try{
 const titles=new Set(), descriptions=new Set()
 const page=await browser.newPage({javaScriptEnabled:false})
 for(const path of manifest.paths){
  const html=await readFile('dist'+path.replace(/\/$/,'')+'/index.html','utf8')
  await page.setContent(html)
  assert.equal(await page.locator('h1').count(),1,path)
  const title=(await page.title()).trim();assert.ok(title,path);assert.ok(!titles.has(title),'Duplicate title: '+path);titles.add(title)
  const description=await page.locator('meta[name="description"]').getAttribute('content');assert.ok(description,path);assert.ok(!descriptions.has(description),'Duplicate description: '+path);descriptions.add(description)
  assert.equal(await page.locator('link[rel="canonical"]').count(),1,path)
  assert.equal(await page.locator('link[rel="canonical"]').getAttribute('href'),manifest.site+path,path)
  const lang=path.match(/^\/(lv|en)(\/|$)/)?.[1]||'ru'
  assert.equal(await page.locator('html').getAttribute('lang'),lang,path)
  for(const language of ['ru','lv','en','x-default'])assert.equal(await page.locator(`link[hreflang="${language}"]`).count(),1,path)
  if(!path.endsWith('/privacy'))assert.ok(!(await page.locator('meta[name="robots"]').getAttribute('content')).includes('noindex'),path)
  const errors=await page.evaluate(()=>({alt:[...document.images].filter(i=>!i.hasAttribute('alt')).length,links:[...document.querySelectorAll('a[href]')].map(a=>a.getAttribute('href')),schemas:[...document.querySelectorAll('script[type="application/ld+json"]')].map(s=>JSON.parse(s.textContent))}))
  if(path.includes('/directions/')) {
   const breadcrumbs=errors.schemas.find(s=>s['@type']==='BreadcrumbList')
   assert.equal(breadcrumbs?.itemListElement.length,3,path)
   assert.equal(breadcrumbs.itemListElement[2].item,manifest.site+path,path)
   assert.ok(breadcrumbs.itemListElement[1].item.endsWith('/#directions'),path)
   assert.equal(breadcrumbs.itemListElement[2].name,await page.locator('h1').innerText(),path)
  }
  assert.equal(errors.alt,0,path)
  assert.ok(errors.links.some(h=>h.includes('/services')),path)
  assert.ok(!errors.schemas.some(s=>s.aggregateRating),path)
 }
 const sitemap=await readFile('dist/sitemap.xml','utf8')
 assert.ok(!/\?lang|growth-point|\/404|\/privacy|\/api\//.test(sitemap))
 for(const path of manifest.paths.filter(p=>!p.endsWith('/privacy')))assert.ok(sitemap.includes('<loc>'+manifest.site+path+'</loc>'),path)
 console.log(`SEO HTML checks passed: ${manifest.paths.length} routes (JavaScript disabled), sitemap.`)
}finally{await browser.close()}
