import type { Locale } from '../i18n'
export const certificates = [
  { id: 'icta', title: 'coachTitle', subtitle: 'coachSubtitle', alt: 'coachAlt', year: '2026', images: { ru: '/images/optimized/icta-ru-full.webp', en: '/images/optimized/icta-en-full.webp' } },
  { id: 'numerology', title: 'numerologyTitle', subtitle: 'numerologySubtitle', alt: 'numerologyAlt', year: '2026', images: { ru: '/images/optimized/numerology-ru-full.webp', lv: '/images/optimized/numerology-lv-full.webp', en: '/images/optimized/numerology-en-full.webp' } },
] satisfies Array<{ id: string; title: string; subtitle: string; alt: string; year: string; images: Partial<Record<Locale, string>> }>
export function certificateImage(certificate: typeof certificates[number], locale: string) {
  const images: Partial<Record<string, string>> = certificate.images
  return images[locale] || images.en || images.ru
}
