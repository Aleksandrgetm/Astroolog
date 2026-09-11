<script setup lang="ts">
withDefaults(defineProps<{ direction?: 'up-right' | 'right' | 'left' }>(), { direction: 'up-right' })
const paths = {
  'up-right': 'M2 22 22 2M15 2h7v7',
  right: 'M2 12h20m-5-5 5 5-5 5',
  left: 'M22 12H2m5-5-5 5 5 5',
}
</script>

<template>
  <svg class="arrow-icon" :class="`arrow-icon--${direction}`" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="0.9" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">
    <path :d="paths[direction]" />
  </svg>
</template>

<style scoped>
.arrow-icon { display: inline-block; width: 22px; height: 22px; flex: 0 0 auto; vertical-align: middle; }
@media (max-width: 600px) {
  .arrow-icon { width: 20px; height: 20px; }
}
@media (hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference) {
  .arrow-icon { transition: transform 275ms cubic-bezier(0.22, 1, 0.36, 1); }
  a:hover .arrow-icon--up-right, button:hover .arrow-icon--up-right { transform: translate(2px, -2px); }
}
</style>
