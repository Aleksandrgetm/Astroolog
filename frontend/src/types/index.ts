import type { LocalizedFields } from '../i18n/localized'
export interface Service extends LocalizedFields<'title' | 'short_description' | 'description' | 'duration' | 'alt'> {
  id: number
  title: string
  slug: string
  short_description: string
  description: string
  image: string
  price: number | null
  duration: string
  is_active: boolean
  sort_order: number
}
export interface RequestInput {
  language?: string
  service_type?: string
  price?: number | null
  name: string
  email: string
  phone: string
  message: string
  consent: boolean
  service_id?: number
  bonus_code?: string
  bonus_price?: number | null
}

export interface BookingSnapshot {
 service_title: string; option_title: string; option_format: string; price: number | null; bonus_code: string | null; bonus_title: string | null; bonus_price: number | null; total_price: number | null
}

export interface BookingOption extends LocalizedFields<'title' | 'format'> {
 image?: string
 is_addon: boolean
 code: string
 service_id: number
 price: number | null
 title: string
 format: string
 sort_order: number
}
