<script setup lang="ts">
import { VTabs,VTab } from 'vuetify/components'
import { languages,type Language } from './contentSchema'
defineProps<{modelValue:Language;complete:Record<string,boolean>}>()
defineEmits<{ 'update:modelValue':[value:Language] }>()
</script>
<template><div class="editor-languages"><v-tabs :model-value="modelValue" @update:model-value="value=>$emit('update:modelValue',value as Language)" aria-label="Язык редактирования"><v-tab v-for="lang in languages" :key="lang" :value="lang" :aria-label="`${lang.toUpperCase()} — ${complete[lang]?'перевод заполнен':'есть незаполненные поля'}`">{{lang.toUpperCase()}}<v-icon :icon="complete[lang]?'mdi-check-circle-outline':'mdi-alert-circle-outline'" :color="complete[lang]?'success':'warning'" size="17" class="ml-2"/></v-tab></v-tabs><p class="admin-muted">{{complete[modelValue]?'Обязательные поля этого языка заполнены.':'Есть незаполненные поля. Добавьте перевод перед публикацией.'}}</p></div></template>
