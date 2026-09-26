import { test } from 'node:test'
import assert from 'node:assert/strict'
import { validPhone, rigaToday } from '../src/utils/validation.ts'
import { languagePath, basePath, localeFromPath } from '../src/i18n/urls.ts'
test('Phone requires 7–15 digits and accepts international formatting',()=>{
 for(const p of ['', '+371 29 580 232','(020) 1234-5678','123.4567'])assert.equal(validPhone(p),true,p)
 for(const p of ['------','12345','1234567890123456','call me','++12345678'])assert.equal(validPhone(p),false,p)
})
test('Riga midnight differs from UTC date',()=>{assert.equal(rigaToday(new Date('2026-09-15T21:01:00Z')),'2026-09-16');assert.equal(rigaToday(new Date('2026-01-15T22:01:00Z')),'2026-01-16')})
test('Locale URLs preserve stable service slugs and canonical home slash',()=>{assert.equal(languagePath('/services/personal-matrix','lv'),'/lv/services/personal-matrix');assert.equal(languagePath('/en/about','ru'),'/about');assert.equal(languagePath('/','en'),'/en/');assert.equal(basePath('/lv/about/'),'/about');assert.equal(localeFromPath('/en/reviews'),'en')})
