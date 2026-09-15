import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getServices, getBookingOptions, errorKey } from '../services/api'
import type { Service, BookingOption } from '../types'
export const useCatalog = defineStore('catalog', () => {
  const options = ref<BookingOption[]>([])
  const services = ref<Service[]>([]), loading = ref(false), error = ref('')
  async function load() { if (loading.value || services.value.length) return; loading.value = true; error.value = ''; try { [services.value, options.value] = await Promise.all([getServices(), getBookingOptions()]) } catch (e) { error.value = errorKey(e) } finally { loading.value = false } }
  return { options, services, loading, error, load }
})
