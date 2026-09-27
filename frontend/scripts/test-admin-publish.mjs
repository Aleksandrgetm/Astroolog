import assert from 'node:assert/strict'
const origin=process.env.TEST_SITE_URL||'http://127.0.0.1:8086'
if(!/^http:\/\/(127\.0\.0\.1|localhost):/.test(origin))throw Error('Local isolated server required')
const login=await fetch(origin+'/api/admin/auth/login',{method:'POST',headers:{Origin:origin,'Content-Type':'application/json'},body:JSON.stringify({email:'cms-test@example.invalid',password:'cms-test-only-password-2026'})});assert.equal(login.status,200)
const {csrf}=await login.json(),cookie=login.headers.get('set-cookie').split(';')[0],headers={Origin:origin,Cookie:cookie,'X-CSRF-Token':csrf}
const before=await (await fetch(origin+'/')).text(),oldAsset=before.match(/src="(\/assets\/[^\"]+\.js)"/)[1]
const response=await fetch(origin+'/api/admin/publish',{method:'POST',headers});assert.equal(response.status,202);const publication=await response.json();console.log('Publication started:',publication.id)
assert.equal((await fetch(origin+'/')).status,200)
assert.equal((await fetch(origin+'/api/admin/publish',{method:'POST',headers})).status,409)
for(let i=0;i<60;i++){
 await new Promise(r=>setTimeout(r,3000));const data=await (await fetch(origin+'/api/admin/dashboard',{headers})).json()
 if(data.publication.status==='running')continue
 assert.equal(data.publication.status,'published',data.publication.error)
 assert.equal((await fetch(origin+oldAsset)).status,200,'Old content-hashed asset unavailable')
 const admin=await fetch(origin+'/admin');assert.equal(admin.status,200);assert.match(admin.headers.get('x-robots-tag'),/noindex/)
 assert.equal((await fetch(origin+'/lv/admin')).status,404)
 assert.equal((await fetch(origin+'/')).status,200)
 console.log('Atomic publication, old assets, admin noindex and localized admin rejection passed');process.exit(0)
}
throw Error('Publish timeout')
