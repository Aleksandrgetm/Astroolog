<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useCatalog } from '../stores/catalog'
import { getLocalizedField } from '../i18n/localized'
import { formatPrice } from '../utils/prices'
import type { BookingOption } from '../types'
const selection = defineModel<BookingOption | null>({ required: true })
const props = defineProps<{ invalid: boolean; disabled: boolean; serviceId?: number }>()
const { t, locale } = useI18n()
const route = useRoute()
const catalog = useCatalog()
const groups = computed(() => catalog.services.filter(s => !props.serviceId || s.id === props.serviceId))
onMounted(() => catalog.load())
watch(() => [catalog.options, route.query.option, props.serviceId], () => {
 const option = catalog.options.find(o => o.code === route.query.option && (!props.serviceId || o.service_id === props.serviceId))
 if (option) selection.value = option
}, { immediate: true })

function price(value: number | null) { return formatPrice(value, locale.value, t('catalog.free'), t('bookingChoice.pricePending')) }
</script>
<template>
  <fieldset class="booking-choice" :disabled="disabled" :aria-invalid="invalid" :aria-describedby="invalid ? 'booking-choice-error' : undefined">
    <legend>{{ t('bookingChoice.service') }}</legend>
    <p v-if="catalog.loading" role="status">{{ t('serviceCards.loadingSessionOptions') }}</p>
    <div v-else-if="catalog.error" role="alert"><p>{{ t('feedback.network') }}</p><button type="button" @click="catalog.load">{{ t('requestForm.tryAgain') }}</button></div>
    <div v-for="group in groups" :key="group.id" class="booking-group"><p class="booking-group-title">{{ getLocalizedField(group, 'title') }}</p>
    <label v-for="option in catalog.options.filter(o => o.service_id === group.id)" :key="option.code" class="booking-option" :class="{ selected: selection?.code === option.code }">
      <input type="radio" name="booking-service-type" :value="option.code" :checked="selection?.code === option.code" @change="selection = option" />
      <span><span class="booking-option-title">{{ getLocalizedField(option, 'title') }}</span><span class="booking-option-price">{{ t('bookingChoice.cost') }}: {{ price(option.price) }}</span></span>
    </label></div>
    <p v-if="invalid" id="booking-choice-error" class="form-error" role="alert">{{ t('bookingChoice.choose') }}</p>
  </fieldset>
</template>
<style scoped>
.booking-group-title { font-family: var(--serif); font-size: 21px; margin: 20px 0 10px; }
.booking-choice { padding: 0; border: 0; margin: 6px 0 24px; min-width: 0; }
.booking-choice legend { font-size: 11px; letter-spacing: 1.4px; text-transform: uppercase; margin-bottom: 12px; }
.booking-option { display: flex; gap: 12px; align-items: flex-start; padding: 16px; border: 1px solid var(--line); margin-top: 10px; cursor: pointer; }
.booking-option.selected { border-color: var(--brown); background: #f5eee5; }
.booking-option input { appearance: none; width: 17px; height: 17px; flex: 0 0 17px; border: 1px solid #a98c6b; border-radius: 50%; margin-top: 3px; }
.booking-option input:checked { background: var(--brown); box-shadow: inset 0 0 0 4px #f5eee5; }
.booking-option-title, .booking-option-price { display: block; line-height: 1.6; }
.booking-option-title { font-size: 13px; }
.booking-option-price { margin-top: 6px; font-size: 12px; color: var(--muted); }
</style>
