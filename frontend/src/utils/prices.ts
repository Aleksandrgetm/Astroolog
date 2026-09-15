export function formatPrice(value: number | null, locale: string, free: string, pending: string): string {
  if (value === null) return pending
  if (value === 0) return free
  return new Intl.NumberFormat(locale, { style: 'currency', currency: 'EUR', minimumFractionDigits: 0, maximumFractionDigits: 2 }).format(value / 100)
}
