import { i18n, type Locale } from './index'
import { pickLocalizedField, type LocalizedFields } from './fields'
export type { LocalizedFields } from './fields'

/** Reads the reactive locale so templates and computed values update immediately. */
export function getLocalizedField<Field extends string>(
  item: Record<Field, string> & LocalizedFields<Field>,
  field: Field,
  locale: Locale = i18n.global.locale.value,
): string {
  return pickLocalizedField(item, field, locale)
}
