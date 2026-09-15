import type { Locale } from '../i18n'
export const certificates = [
  { id: 'icta', title: 'coachTitle', subtitle: 'coachSubtitle', alt: 'coachAlt', year: '2026', images: { ru: '/images/certificates/icta-ru.png', en: '/images/certificates/icta-en.png' } },
  { id: 'numerology', title: 'numerologyTitle', subtitle: 'numerologySubtitle', alt: 'numerologyAlt', year: '2026', images: { ru: '/images/certificates/numerology-ru.png', lv: '/images/certificates/numerology-lv.png', en: '/images/certificates/numerology-en.png' } },
] satisfies Array<{ id: string; title: string; subtitle: string; alt: string; year: string; images: Partial<Record<Locale, string>> }>
export function certificateImage(certificate: typeof certificates[number], locale: string) {
  const images: Partial<Record<string, string>> = certificate.images
  return images[locale] || images.en || images.ru
}
