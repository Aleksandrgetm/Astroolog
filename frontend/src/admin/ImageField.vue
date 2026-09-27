<script setup lang="ts">
import {ref} from 'vue'
import MediaLibrary from './MediaLibrary.vue'
withDefaults(defineProps<{modelValue?:string;label?:string;missing?:string}>(),{label:'Фотография'})
const emit=defineEmits<{'update:modelValue':[value:string];selected:[asset:any]}>()
const open=ref(false)
function choose(asset:any){emit('update:modelValue',asset.storage_path);emit('selected',asset);open.value=false}
</script>
<template><div class="image-field"><h3>{{label}}</h3><img v-if="modelValue" :src="modelValue" :alt="'Предпросмотр: '+label" class="image-field-preview"/><div v-else class="image-field-empty"><v-icon icon="mdi-image-outline" size="32"/><p>{{missing||'Изображение пока не выбрано'}}</p></div><div class="admin-actions"><v-btn variant="outlined" @click="open=true">{{modelValue?'Заменить изображение':'Загрузить изображение'}}</v-btn><v-btn variant="text" @click="open=true">Выбрать из медиатеки</v-btn></div><v-dialog v-model="open" max-width="1120" content-class="admin-media-dialog"><MediaLibrary v-if="open" picker @select="choose" @close="open=false"/></v-dialog></div></template>
