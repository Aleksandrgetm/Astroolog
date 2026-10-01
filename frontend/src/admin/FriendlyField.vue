<script setup lang="ts">
import { computed } from 'vue'
import { friendly } from './contentSchema'
const props=defineProps<{modelValue?:string|null;label:string;rows?:number;help?:string;error?:string;expertName?:string;template?:boolean;counter?:boolean;required?:boolean;readonly?:boolean}>()
const emit=defineEmits<{'update:modelValue':[value:string]}>()
const value=computed(()=>friendly(props.modelValue||'',props.expertName))
const hint=computed(()=>[props.help,props.template?`Имя специалиста подставляется автоматически: ${props.expertName||'Елена Захарова'}. Его можно изменить в настройках.`:''].filter(Boolean).join(' '))
function update(value:string){emit('update:modelValue',props.template?value.replaceAll(props.expertName||'Елена Захарова','{name}'):value)}
</script>
<template><v-textarea v-if="rows" :model-value="value" @update:model-value="update" :label="label" :rows="rows" :hint="hint" :persistent-hint="!!hint" :error-messages="error" :counter="counter" :aria-required="required" :readonly="readonly" variant="outlined"/><v-text-field v-else :model-value="value" @update:model-value="update" :label="label" :hint="hint" :persistent-hint="!!hint" :error-messages="error" :counter="counter" :aria-required="required" :readonly="readonly" variant="outlined"/></template>
