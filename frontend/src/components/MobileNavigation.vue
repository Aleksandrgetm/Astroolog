<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LanguageSwitcher from './LanguageSwitcher.vue'
const props = defineProps<{ links: string[][] }>()
const { t } = useI18n()
const route = useRoute()
const open = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const closeButton = ref<HTMLButtonElement | null>(null)
const mobileLinks = computed(() => [[t('siteLayout.home'), '/'], ...props.links])
function active(target: string) {
  const [path, hash] = target.split('#')
  if (path === '/') return route.path === '/' && route.hash === (hash ? `#${hash}` : '')
  return route.path === path || route.path.startsWith(`${path}/`)
}
watch(() => `${route.path}${route.hash}`, () => { open.value = false })
let breakpoint: MediaQueryList | undefined
function closeOnDesktop() { if (breakpoint?.matches) open.value = false }
onMounted(() => {
  breakpoint = matchMedia('(min-width: 1051px)')
  breakpoint.addEventListener('change', closeOnDesktop)
})
onUnmounted(() => breakpoint?.removeEventListener('change', closeOnDesktop))
function restoreFocus() { if (trigger.value?.getClientRects().length) trigger.value.focus({ preventScroll: true }) }
</script>

<template>
  <button ref="trigger" class="mobile-toggle menu-toggle" type="button" :aria-label="t('siteLayout.openMenu')" :aria-expanded="open" aria-controls="mobile-navigation" @click="open = true">
    <span class="menu-lines" :class="{ 'is-open': open }" aria-hidden="true"><i></i><i></i><i></i></span>
  </button>
  <v-dialog v-model="open" class="mobile-navigation-dialog" transition="mobile-navigation" scroll-strategy="block" :aria-label="t('siteLayout.mobileNavigation')" @after-enter="closeButton?.focus()" @after-leave="restoreFocus">
    <div id="mobile-navigation" class="mobile-menu">
      <div class="menu-top">
        <span class="eyebrow">{{ t('siteLayout.expertNameUppercase') }}</span>
        <button ref="closeButton" class="menu-toggle" type="button" :aria-label="t('siteLayout.closeMenu')" @click="open = false">
          <span class="menu-lines is-open" aria-hidden="true"><i></i><i></i><i></i></span>
        </button>
      </div>
      <nav :aria-label="t('siteLayout.mobileNavigation')">
        <LocaleLink v-for="item in mobileLinks" :key="item[1]" :to="item[1]!" active-class="" exact-active-class="" :class="{ 'is-nav-active': active(item[1]!) }" :aria-current="active(item[1]!) ? 'location' : undefined" @click="open = false">
          {{ item[0] }} <ArrowIcon />
        </LocaleLink>
      </nav>
      <div class="mobile-menu-bottom">
        <LanguageSwitcher class="mobile-language" />
        <LocaleLink class="button mobile-menu-cta" to="/contacts?booking=1" @click="open = false">{{ t('siteLayout.bookASession') }} <ArrowIcon /></LocaleLink>
      </div>
    </div>
  </v-dialog>
</template>

<style>
.menu-toggle {
  display: inline-flex; align-items: center; justify-content: center;
  width: 46px; height: 46px; flex: 0 0 46px;
  border: 0; border-radius: 2px; background: #f1ebe2; color: #765239; cursor: pointer;
}
.menu-lines { position: relative; width: 23px; height: 18px; }
.menu-lines i { position: absolute; left: 0; width: 23px; height: 1px; background: currentColor; transition: transform 350ms cubic-bezier(0.22,1,0.36,1), opacity 250ms; }
.menu-lines i:nth-child(1) { top: 1px; }
.menu-lines i:nth-child(2) { top: 8px; }
.menu-lines i:nth-child(3) { top: 15px; }
.menu-lines.is-open i:nth-child(1) { transform: translateY(7px) rotate(45deg); }
.menu-lines.is-open i:nth-child(2) { opacity: 0; }
.menu-lines.is-open i:nth-child(3) { transform: translateY(-7px) rotate(-45deg); }
.mobile-navigation-dialog .v-overlay__scrim { background: #40372f; opacity: .22; }
.mobile-navigation-dialog .v-overlay__content {
  position: fixed; top: 0; right: 0; margin: 0;
  width: min(100%, 560px); max-width: 100%; height: 100dvh; max-height: 100dvh;
  border-radius: 0; box-shadow: none;
}
.mobile-navigation-dialog .mobile-menu {
  display: flex; flex-direction: column; height: 100%; min-height: 0;
  padding: max(22px, env(safe-area-inset-top)) max(28px, env(safe-area-inset-right)) max(26px, env(safe-area-inset-bottom)) max(28px, env(safe-area-inset-left));
  background: #faf7f2; overflow-y: auto; overscroll-behavior: contain;
}
.mobile-navigation-dialog .menu-top { gap: 20px; }
.mobile-navigation-dialog .menu-top .eyebrow { margin: 0; line-height: 1.6; }
.mobile-navigation-dialog nav { margin: 28px 0; }
.mobile-navigation-dialog nav a {
  align-items: center; min-height: 57px; padding: 11px 0;
  font-size: clamp(29px, 4.2vw, 36px); line-height: 1.2; color: var(--ink);
}
.mobile-navigation-dialog nav a.is-nav-active { color: var(--brown); }
.mobile-navigation-dialog nav a .arrow-icon { opacity: .55; }
.mobile-navigation-dialog nav a.is-nav-active .arrow-icon { opacity: 1; }
.mobile-menu-bottom { margin-top: auto; padding-top: 12px; }
.mobile-navigation-dialog .mobile-language { margin-bottom: 24px; gap: 14px; }
.mobile-navigation-dialog .mobile-language button { min-width: 44px; min-height: 44px; }
.mobile-navigation-dialog .mobile-menu-cta { width: 100%; justify-content: space-between; }
.mobile-navigation-enter-active { transition: opacity 420ms cubic-bezier(0.22,1,0.36,1), transform 420ms cubic-bezier(0.22,1,0.36,1); }
.mobile-navigation-leave-active { transition: opacity 250ms ease, transform 250ms ease; }
.mobile-navigation-enter-from, .mobile-navigation-leave-to { opacity: 0; transform: translateX(20px); }
@media (max-width: 600px) {
  .mobile-navigation-dialog .v-overlay__content { width: 100%; }
  .mobile-navigation-enter-from, .mobile-navigation-leave-to { transform: translateY(-8px); }
}
@media (hover: none) {
  .mobile-navigation-dialog .button:hover { transform: none; box-shadow: none; background: var(--brown); }
}
@media (prefers-reduced-motion: reduce) {
  .menu-lines i, .mobile-navigation-enter-active, .mobile-navigation-leave-active { transition: none; }
  .mobile-navigation-enter-from, .mobile-navigation-leave-to { transform: none; }
}
</style>
