import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
const icons = ['favicon.svg', 'favicon.ico', 'favicon-32x32.png', 'favicon-192x192.png', 'apple-touch-icon.png']
for (const name of icons) {
 const source = await readFile('public/' + name)
 assert.deepEqual(await readFile('dist/' + name), source, `${name}: build copy`)
 if (name.endsWith('.png')) {
  const size = name.includes('32x32') ? 32 : name.includes('192x192') ? 192 : 180
  assert.equal(source.readUInt32BE(16), size)
  assert.equal(source.readUInt32BE(20), size)
 }
}
assert.ok(!(await readFile('public/favicon.svg', 'utf8')).includes('<text'), 'SVG must not depend on installed fonts')
const { paths } = JSON.parse(await readFile('dist/routes.json', 'utf8'))
const pages = [...paths.map(p => p.replace(/\/$/, '') + '/index.html'), '/404.html', '/lv/404.html', '/en/404.html', '/admin.html']
for (const path of pages) {
 const html = await readFile('dist' + path, 'utf8')
 for (const name of icons) assert.ok(html.includes(`href="/${name}?v=ez2"`), `${path}: ${name}`)
 assert.ok(!html.includes('/vite.svg'), `${path}: default icon`)
}
console.log(`Brand icons verified in ${pages.length} generated pages across RU/LV/EN.`)
