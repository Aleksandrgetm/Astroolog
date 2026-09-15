export interface WhatsAppRequest {
  service_title?: string
  price_label?: string
  name: string
  email: string
  phone?: string | null
  message?: string | null
  preferred_date?: string | null
}

/** Owner-facing labels are Russian; submitted values are preserved in their original language. */
export function buildWhatsAppMessage(kind: 'contact' | 'bookings', data: WhatsAppRequest, serviceTitle = ''): string {
  const lines = [kind === 'bookings' ? 'Новая запись с сайта' : 'Новый вопрос с сайта', '', `Имя: ${data.name}`, `Email: ${data.email}`]
  if (data.phone?.trim()) lines.push(`Телефон: ${data.phone}`)
  if (kind === 'bookings') {
    if (data.service_title?.trim()) lines.push(`Услуга: ${data.service_title}`)
    if (data.price_label?.trim()) lines.push(`Стоимость: ${data.price_label}`)
    if (serviceTitle.trim()) lines.push(`Формат встречи: ${serviceTitle}`)
    if (data.preferred_date?.trim()) lines.push(`Желаемая дата: ${data.preferred_date}`)
  }
  if (data.message?.trim()) lines.push('', kind === 'bookings' ? 'О чём хочется поговорить:' : 'Вопрос:', data.message)
  return lines.join('\n')
}

export function buildWhatsAppUrl(message: string, phone = import.meta.env?.VITE_WHATSAPP_PHONE ?? ''): string {
  const digits = phone.replace(/[^0-9]/g, '')
  if (!/^[1-9][0-9]{6,14}$/.test(digits)) return ''
  return `https://wa.me/${digits}?text=${encodeURIComponent(message)}`
}
