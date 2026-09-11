// Shared public contact; hide Telegram actions until the owner supplies a URL.
export const telegramUrl = (import.meta.env?.VITE_TELEGRAM_URL ?? '').trim()

export function buildTelegramUrl(message: string, contact = telegramUrl): string {
  if (!contact) return ''
  const url = new URL(contact)
  url.searchParams.delete('profile')
  url.searchParams.set('text', message)
  return url.toString()
}
