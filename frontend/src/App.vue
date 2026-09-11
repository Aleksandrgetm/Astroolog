<script setup lang="ts">
import SiteLayout from './layouts/SiteLayout.vue'

function disableLeavingPage(element: Element) {
  // The outgoing page remains visible briefly, but must no longer receive input.
  element.setAttribute('inert', '')
  element.setAttribute('aria-hidden', 'true')
}
</script>

<template>
  <v-app>
    <SiteLayout>
      <div class="page-stage">
        <RouterView v-slot="{ Component, route }">
          <Transition name="page" appear @before-leave="disableLeavingPage">
            <!-- A single element also supports pages with multiple root nodes. -->
            <div v-if="Component" :key="route.path" class="page-content">
              <component :is="Component" />
            </div>
          </Transition>
        </RouterView>
      </div>
    </SiteLayout>
  </v-app>
</template>

<style scoped>
.page-stage { position: relative; }
.page-content { display: flow-root; }
.page-enter-active {
  transition: opacity 550ms cubic-bezier(0.22, 1, 0.36, 1),
    transform 550ms cubic-bezier(0.22, 1, 0.36, 1);
}
.page-enter-from { opacity: 0; transform: translateY(10px); }
.page-enter-to { opacity: 1; transform: translateY(0); }
.page-leave-active {
  position: absolute;
  inset: 0 0 auto;
  width: 100%;
  pointer-events: none;
  transition: opacity 140ms ease-out;
}
.page-leave-to { opacity: 0; }
@media (prefers-reduced-motion: reduce) {
  .page-enter-active, .page-leave-active { transition: none; }
  .page-enter-from, .page-enter-to, .page-leave-to {
    opacity: 1;
    transform: none;
  }
}
</style>
