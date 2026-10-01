import { chromium } from 'playwright'
import { createServer } from 'node:http'
import { readFile, writeFile, mkdir, stat } from 'node:fs/promises'
import { resolve, dirname, extname } from 'node:path'
import { loadEnv } from 'vite'
const mode=process.env.BUILD_MODE || 'production'
const env={...loadEnv(mode,process.cwd(),''),...process.env}
const site=env.VITE_SITE_URL?.replace(/\/$/,'')
if(!site || new URL(site).protocol!=='https:' || new URL(site).origin!==site) throw Error('VITE_SITE_URL must be an absolute HTTPS origin (no path).')
const api=env.PRERENDER_API_URL || 'http://127.0.0.1:8080/api'
async function json(path){const r=await fetch(api+path,{signal:AbortSignal.timeout(15000)});if(!r.ok)throw Error(`Prerender API ${path}: ${r.status}`);return r.json()}
const services=await json('/services'),options=await json('/booking-options'),content=await json('/content')
const directions=content.directions
if(!Array.isArray(services)||!services.length||!Array.isArray(options))throw Error('Invalid live catalog')
if(services.some(s=>!s.is_active || !/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(s.slug) || s.slug==='growth-point') || new Set(services.map(s=>s.slug)).size!==services.length)throw Error('Invalid or duplicate active service slug')
const root=resolve(env.SITE_BUILD_DIR||'dist'),template=await readFile(resolve(root,'index.html'))
await writeFile(resolve(root,'admin.html'),template)
const types={'.js':'text/javascript','.css':'text/css','.webp':'image/webp','.png':'image/png','.svg':'image/svg+xml','.woff2':'font/woff2'}
const server=createServer(async(req,res)=>{try{const path=new URL(req.url,'http://localhost').pathname;const file=resolve(root,'.'+decodeURIComponent(path));if(!file.startsWith(root+'/')){res.end(template);return}if((await stat(file).catch(()=>null))?.isFile()){res.setHeader('Content-Type',types[extname(file)]||'application/octet-stream');res.end(await readFile(file))}else{res.setHeader('Content-Type','text/html');res.end(template)}}catch{res.writeHead(500).end()}})
await new Promise(r=>server.listen(0,'127.0.0.1',r))
const address=`http://127.0.0.1:${server.address().port}`
const languages=['ru','lv','en'],localized=(p,l)=>l==='ru'?p:`/${l}${p}`
const bases=['/','/about','/services','/reviews','/contacts','/privacy','/cookies','/terms',...services.map(s=>'/services/'+s.slug),...directions.map(d=>'/directions/'+d.slug)]
let browser
try{
 browser=await chromium.launch(env.PLAYWRIGHT_CHANNEL?{channel:env.PLAYWRIGHT_CHANNEL}:{})
 const page=await browser.newPage({reducedMotion:'reduce',viewport:{width:1440,height:1000}})
 await page.route('**/api/**',async route=>{const path=new URL(route.request().url()).pathname.replace(/^.*\/api/,'');const data=path==='/content'?content:path==='/services'?services:path==='/booking-options'?options:services.find(s=>path==='/services/'+s.slug);await route.fulfill({status:data?200:404,contentType:'application/json',body:JSON.stringify(data||{code:'not_found'})})})
 const paths=[]
 for(const lang of languages){for(const base of [...bases,'/__404']){
  const path=localized(base,lang)
  await page.goto(address+path,{waitUntil:'networkidle'})
  await page.locator('h1').waitFor()
  if(await page.locator('h1').count()!==1)throw Error('Expected one H1: '+path)
  await page.evaluate(()=>{document.querySelectorAll('[data-cookie-banner], [data-cookie-settings]').forEach(el=>el.remove());document.querySelectorAll('[inert]').forEach(el=>el.remove());window.scrollTo(0,0)})
  const output=base==='/__404'?localized('/404.html',lang):path.replace(/\/$/,'')+'/index.html'
  const filename=resolve(root,'.'+output);await mkdir(dirname(filename),{recursive:true})
  await writeFile(filename,'<!doctype html>\n'+await page.locator('html').evaluate(el=>el.outerHTML))
  if(base!=='/__404')paths.push(path)
 }}
 await writeFile(resolve(root,'routes.json'),JSON.stringify({site,paths,services:services.map(s=>s.slug)},null,2))
 const escape=s=>s.replace(/&/g,'&amp;').replace(/"/g,'&quot;').replace(/</g,'&lt;')
 const entries=bases.filter(p=>!['/privacy','/cookies','/terms'].includes(p)).flatMap(base=>languages.map(lang=>`<url><loc>${escape(site+localized(base,lang))}</loc>${[...languages,'x-default'].map(l=>`<xhtml:link rel="alternate" hreflang="${l}" href="${escape(site+localized(base,l==='x-default'?'ru':l))}"/>`).join('')}</url>`))
 await writeFile(resolve(root,'sitemap.xml'),'<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">'+entries.join('\n')+'</urlset>')
 await writeFile(resolve(root,'robots.txt'),`User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /api/admin\nSitemap: ${site}/sitemap.xml\n`)
 console.log(`Prerendered ${paths.length} pages and 3 localized 404 pages from live catalog.`)
}finally{await browser?.close();server.close()}
