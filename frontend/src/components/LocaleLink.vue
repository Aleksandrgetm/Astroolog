<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter, type RouteLocationRaw } from 'vue-router'
import { localizedRoute } from '../i18n/routing'
const props = defineProps<{ to: RouteLocationRaw }>()
const router = useRouter()
const route = useRoute()
const destination = computed(() => {
  // Track query-only language changes as well as path changes.
  void route.query.lang
  return localizedRoute(router, props.to)
})
</script>

<template>
  <RouterLink :to="destination"><slot /></RouterLink>
</template>
