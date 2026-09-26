import { chromium } from 'playwright'
import assert from 'node:assert/strict'
import { readFile, mkdir, writeFile } from 'node:fs/promises'
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8087'
const manifest=JSON.parse(await readFile('dist/routes.json','utf8'))
const browser=await chromium.launch(process.env.PLAYWRIGHT_CHANNEL?{channel:process.env.PLAYWRIGHT_CHANNEL}:{})
try{
 const page=await browser.newPage({reducedMotion:'reduce'})
 const errors=[];page.on('pageerror',error=>errors.push(error.message))
 let checks=0
 for(const width of [390,430,768,1024,1440]){
  await page.setViewportSize({width,height:950})
  for(const path of manifest.paths){
   const response=await page.goto(origin+path,{waitUntil:'networkidle'})
   assert.equal(response.status(),200,path)
   await page.locator('h1').waitFor()
   const overflow=await page.evaluate(()=>document.documentElement.scrollWidth>window.innerWidth+1)
   assert.equal(overflow,false,`${width} ${path}`)
   assert.equal(await page.locator('.page-content:not([inert]) h1').count(),1,path)
   checks++
  }
  console.log('Responsive passed at',width)
 }
 for(const [path,status,location] of [['/services?lang=lv&option=personality',301,'/lv/services?option=personality'],['/en/services/growth-point',301,'/en/services'],['/lv/no-such-route',404],['/services/not-real',404]]){
  const r=await page.request.get(origin+path,{maxRedirects:0});assert.equal(r.status(),status,path);if(location)assert.equal(r.headers().location,location)
 }
 assert.deepEqual(errors,[])
 await mkdir('../artifacts/pre-admin',{recursive:true})
 await writeFile('../artifacts/pre-admin/browser-results.json',JSON.stringify({checks,errors,http:true},null,2))
 console.log(`${checks} responsive route checks passed; HTTP 301/404 passed.`)
}finally{await browser.close()}
