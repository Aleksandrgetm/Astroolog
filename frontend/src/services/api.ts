import axios from 'axios'
import type { Service, Testimonial, RequestInput } from '../types'
export const api = axios.create({ baseURL: import.meta.env.VITE_API_URL || '/api', timeout: 15000, headers: { 'Content-Type': 'application/json' } })
export const getServices = () => api.get<Service[]>('/services').then(r => r.data)
export const getService = (slug: string) => api.get<Service>(`/services/${encodeURIComponent(slug)}`).then(r => r.data)
export const getTestimonials = () => api.get<Testimonial[]>('/testimonials').then(r => r.data)
export const submitRequest = (kind: 'bookings' | 'contact', data: RequestInput) => api.post<{ message: string; code: string }>(`/${kind}`, data).then(r => r.data)
const knownCodes = ['server_error', 'validation_error', 'invalid_request', 'not_found'] as const
export type ErrorKey = `feedback.${typeof knownCodes[number] | 'network'}`
// Persist message keys, not translated strings, so visible errors follow locale changes.
export function errorKey(error: unknown): ErrorKey {
  if (!axios.isAxiosError<{ code?: string }>(error)) return 'feedback.server_error'
  if (!error.response) return 'feedback.network'
  const code = knownCodes.find(code => code === error.response?.data?.code)
  if (code) return `feedback.${code}`
  if (error.response.status === 404) return 'feedback.not_found'
  if (error.response.status === 400) return 'feedback.validation_error'
  return 'feedback.server_error'
}
