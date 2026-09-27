<script setup lang="ts">
defineProps<{dirty:boolean;busy:boolean;error?:string}>()
defineEmits<{save:[];cancel:[]}>()
</script>
<template><div class="editor-savebar"><div><strong role="status">{{dirty?'Есть несохранённые изменения':'Все изменения сохранены'}}</strong><p v-if="error" class="editor-error" role="alert">{{error}}</p><p v-else class="admin-muted">После сохранения изменения доступны на сайте. Публикация обновляет поисковую версию страниц.</p></div><div class="editor-save-actions"><v-btn variant="text" :disabled="busy||!dirty" @click="$emit('cancel')">Отмена</v-btn><v-btn color="primary" :disabled="!dirty||busy" :loading="busy" @click="$emit('save')">Сохранить</v-btn></div></div></template>
