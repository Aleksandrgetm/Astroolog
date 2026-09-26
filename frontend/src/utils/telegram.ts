export const telegramUrl = (import.meta.env?.VITE_TELEGRAM_URL ?? '').trim()
export function buildTelegramUrl(message:string, contact=telegramUrl):string {
 try {
  const url=new URL(contact)
  if(url.protocol!=='https:' || !['t.me','telegram.me'].includes(url.hostname) || !/^\/[A-Za-z][A-Za-z0-9_]{3,31}\/?$/.test(url.pathname) || url.username || url.password || url.port)return ''
  url.searchParams.delete('profile');url.searchParams.set('text',message)
  return url.toString()
 } catch { return '' }
}
