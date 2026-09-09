import type { Locale } from './index'
export type LocalizedFields<Field extends string> = Partial<Record<`${Field}_${Locale}`, string | null>>

/** Empty or missing translations fall back to Russian, then the original API field. */
export function pickLocalizedField<Field extends string>(
  item: Record<Field, string> & LocalizedFields<Field>,
  field: Field,
  locale: Locale,
): string {
  const localized = item[`${field}_${locale}` as `${Field}_${Locale}`]
  const russian = item[`${field}_ru` as `${Field}_${Locale}`]
  if (typeof localized === 'string' && localized.trim()) return localized
  if (typeof russian === 'string' && russian.trim()) return russian
  return item[field]
}
