import { test } from 'node:test'
import assert from 'node:assert/strict'
import ru from '../src/i18n/locales/ru.ts'
import lv from '../src/i18n/locales/lv.ts'
import en from '../src/i18n/locales/en.ts'
import { pickLocalizedField } from '../src/i18n/fields.ts'

function flatten(value: Record<string, unknown>, prefix = ''): Record<string, string> {
  return Object.fromEntries(Object.entries(value).flatMap(([key, child]) => {
    const name = prefix ? `${prefix}.${key}` : key
    return typeof child === 'string' ? [[name, child]] : Object.entries(flatten(child as Record<string, unknown>, name))
  }))
}
test('Every interface message is present and nonempty in all three languages', () => {
  const source = flatten(ru)
  for (const messages of [lv, en]) {
    const flat = flatten(messages)
    assert.deepEqual(Object.keys(flat).sort(), Object.keys(source).sort())
    for (const [key, value] of Object.entries(flat)) assert.ok(value.trim(), key)
  }
})
test('Expert name and monogram are never translated', () => {
  for (const messages of [lv, en]) {
    for (const key of ['expertName', 'expertNameUppercase', 'monogramFirst', 'monogramSecond'] as const) {
      assert.equal(messages.siteLayout[key], ru.siteLayout[key])
    }
  }
})
test('API content uses the chosen language and supports old API responses', () => {
  const item = { title: 'Original', title_ru: 'Русский', title_lv: 'Latviski', title_en: 'English' }
  assert.equal(pickLocalizedField(item, 'title', 'ru'), 'Русский')
  assert.equal(pickLocalizedField(item, 'title', 'lv'), 'Latviski')
  assert.equal(pickLocalizedField(item, 'title', 'en'), 'English')
  for (const missing of [undefined, null, '', '  ']) {
    assert.equal(pickLocalizedField({ ...item, title_lv: missing }, 'title', 'lv'), 'Русский')
  }
  assert.equal(pickLocalizedField({ title: 'Legacy' }, 'title', 'en'), 'Legacy')
  assert.equal(pickLocalizedField({ title: 'Legacy', title_ru: '', title_en: '' }, 'title', 'en'), 'Legacy')
})
