import type { LocalizedFields } from '../i18n/localized'
export interface Service extends LocalizedFields<'title' | 'short_description' | 'description' | 'duration'> {
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
export interface Testimonial extends LocalizedFields<'text' | 'client_name'> {
  id: number
  client_name: string
  text: string
  rating: number
}
export interface RequestInput {
  name: string
  email: string
  phone: string
  message: string
  consent: boolean
  service_id?: number
  preferred_date?: string
}
