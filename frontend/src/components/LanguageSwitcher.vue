<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { setLocale, supportedLocales } from '../i18n'
const { t, locale } = useI18n()
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
      @click="setLocale(language)"
    >{{ language.toUpperCase() }}</button>
  </div>
</template>
