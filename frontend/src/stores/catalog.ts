import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getServices, errorKey } from '../services/api'
import type { Service } from '../types'
export const useCatalog = defineStore('catalog', () => {
  const services = ref<Service[]>([]), loading = ref(false), error = ref('')
  async function load() { if (loading.value || services.value.length) return; loading.value = true; error.value = ''; try { services.value = await getServices() } catch (e) { error.value = errorKey(e) } finally { loading.value = false } }
  return { services, loading, error, load }
})
