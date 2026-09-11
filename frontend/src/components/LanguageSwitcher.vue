<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { supportedLocales, type Locale } from '../i18n'
import { useRoute, useRouter } from 'vue-router'
const { t, locale } = useI18n()
const route = useRoute(), router = useRouter()
function changeLocale(language: Locale) {
  return router.push({ path: route.path, query: { ...route.query, lang: language }, hash: route.hash })
}
</script>

<template>
  <div class="language-switcher" role="group" :aria-label="t('language.label')">
    <button
      v-for="language in supportedLocales"
      :key="language"
      type="button"
      :lang="language"
      :aria-label="t(`language.${language}`)"
      :aria-pressed="locale === language"
      :class="{ active: locale === language }"
      @click="changeLocale(language)"
    >{{ language.toUpperCase() }}</button>
  </div>
</template>
