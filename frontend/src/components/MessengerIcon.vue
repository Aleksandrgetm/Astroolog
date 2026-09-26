<script setup lang="ts">
defineProps<{ messenger: 'whatsapp' | 'telegram' }>()
</script>
<template>
 <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">
  <template v-if="messenger === 'telegram'"><path d="m21 3-4 18-6-5-4 3 1-7-6-2 19-7Z"/><path d="m8 12 9-6-6 10"/></template>
  <template v-else><path d="M20 11.5a8 8 0 0 1-11.8 7L4 20l1.4-4.1A8 8 0 1 1 20 11.5Z"/><path d="M9 7.5c-.8 0-1 1.2-.6 2.4 1 2.8 3 4.6 5.7 5.1 1.2.2 2-.5 2-1.3l-2.2-1.2-.9.8a7 7 0 0 1-3-3l.7-.9L9.8 7.5Z"/></template>
 </svg>
</template>
