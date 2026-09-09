<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
<template><section class="section container not-found"><span class="eyebrow">{{ t('notFound.404ALittleOffThePath') }}</span><h1>{{ t('notFound.thisPageIsntHere') }}<br>{{ t('notFound.butTheresAlways') }} <em>{{ t('notFound.aNewBeginning') }}</em></h1><p>{{ t('notFound.theLinkMayHaveChangedLetsGo') }}</p><router-link class="button" to="/">{{ t('notFound.backToHome') }}</router-link></section></template>
