<script setup lang="ts">
import { VAlert, VCard } from 'vuetify/components'
import {ref} from 'vue'
import {useRouter} from 'vue-router'
import {adminApi,adminSession,adminError} from './api'
const email=ref(''),password=ref(''),error=ref(''),busy=ref(false),router=useRouter()
async function login(){if(busy.value)return;busy.value=true;error.value='';try{const {data}=await adminApi.post('/auth/login',{email:email.value,password:password.value});Object.assign(adminSession,data);password.value='';router.replace('/admin')}catch(e){error.value=adminError(e)}finally{busy.value=false}}
</script>
<template><v-card max-width="440" class="mx-auto mt-12 pa-6"><h1>Вход в Astroolog</h1><p class="admin-muted mb-6">Управление контентом и заявками</p><v-form @submit.prevent="login"><v-text-field v-model="email" label="Email" type="email" autocomplete="username" required/><v-text-field v-model="password" label="Пароль" type="password" autocomplete="current-password" required/><v-alert v-if="error" type="error" class="mb-4">{{error}}</v-alert><v-btn color="primary" type="submit" block :loading="busy">Войти</v-btn></v-form></v-card></template>
