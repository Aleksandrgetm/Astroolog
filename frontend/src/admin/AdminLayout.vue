<script setup lang="ts">
import { VLayout, VNavigationDrawer, VMain, VList, VListItem } from 'vuetify/components'
import './admin.css'
import {ref,computed} from 'vue'
import {useRoute,useRouter} from 'vue-router'
import {useDisplay} from 'vuetify'
import {adminApi,adminSession} from './api'
const route=useRoute(),router=useRouter(),{lgAndUp}=useDisplay(),drawer=ref(false)
const login=computed(()=>route.path==='/admin/login')
const links=[['','Обзор'],['content','Контент сайта'],['services','Услуги'],['directions','Направления'],['reviews','Отзывы'],['certificates','Сертификаты'],['requests','Заявки'],['media','Медиа'],['seo','SEO'],['settings','Настройки'],['audit','Журнал действий']]
async function logout(){await adminApi.post('/auth/logout');adminSession.user=null;adminSession.csrf='';router.replace('/admin/login')}
</script>
<template><div class="admin-surface"><v-layout>
 <v-navigation-drawer v-if="!login" :model-value="lgAndUp || drawer" @update:model-value="drawer=$event" :permanent="lgAndUp" :temporary="!lgAndUp" width="252">
  <div class="admin-brand">Astroolog <small>Управление сайтом</small></div>
  <v-list nav><v-list-item v-for="[path,title] in links" :key="path" :to="'/admin'+(path?'/'+path:'')" :title="title" :exact="!path" :active="path ? route.path.startsWith('/admin/'+path)||(path==='services'&&route.path.startsWith('/admin/options/')) : route.path==='/admin'" @click="drawer=false" /></v-list>
  <template #append><div class="admin-account"><small>{{adminSession.user?.email}}</small><v-btn variant="text" block @click="logout">Выйти</v-btn></div></template>
 </v-navigation-drawer>
 <v-main><header v-if="!login" class="admin-toolbar"><v-btn v-if="!lgAndUp" aria-label="Открыть меню" icon="mdi-menu" variant="text" @click="drawer=true"/><span>Панель администратора</span><a href="/" target="_blank" rel="noopener">Открыть сайт</a></header><main class="admin-main"><slot/></main></v-main>
 </v-layout></div></template>
<style>
.admin-surface { background:#f5f4f1;min-height:100vh;font-family:Arial,sans-serif;color:#302c28; }
.admin-surface h1 { font:500 30px/1.25 Arial,sans-serif;margin:0 0 24px; }
.admin-surface h2 { font:500 21px/1.4 Arial,sans-serif; }
.admin-surface h3 { font:500 18px/1.4 Arial,sans-serif; }
.admin-brand { padding:28px 24px;font-size:24px; }.admin-brand small {display:block;font-size:12px;color:#777;margin-top:8px;}
.admin-toolbar { display:flex;align-items:center;gap:14px;background:white;border-bottom:1px solid #ddd;padding:16px 28px; }.admin-toolbar a {margin-left:auto;font-size:13px;}
.admin-surface .v-main {min-width:0;}
.admin-main {min-width:0;padding:32px;max-width:1500px;margin:auto;}.admin-account {padding:20px;overflow-wrap:anywhere;}
.admin-actions > .v-input {min-width:0;flex:1 1 240px;}
.admin-actions { display:flex;gap:12px;flex-wrap:wrap;margin:18px 0; }.admin-table {background:white;border:1px solid #e0ddd7;overflow:auto;}.admin-table table {width:100%;border-collapse:collapse;}.admin-table th,.admin-table td {text-align:left;padding:14px;border-bottom:1px solid #eee;vertical-align:top;font-size:14px;}.admin-table button {cursor:pointer;color:#795438;}.admin-title-cell {min-width:180px;max-width:500px;overflow-wrap:anywhere;}
.admin-dialog > .v-card {overflow-y:auto!important;}
.admin-editor {padding:24px;background:white;}.admin-editor h2 {margin-bottom:20px;}.admin-editor .v-tab {letter-spacing:0;}.admin-editor-grid {display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:20px;}.admin-editor-grid .wide {grid-column:1/-1;}.admin-preview {display:block;max-width:100%;height:180px;object-fit:contain;margin-bottom:16px;}.admin-muted {color:#777;font-size:13px;line-height:1.6;}.admin-cards {display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:16px;}.admin-card {min-width:0;overflow-wrap:anywhere;background:#fff;padding:24px;border:1px solid #ddd;}.admin-card strong {display:block;font-size:30px;margin-top:14px;}
@media(max-width:600px){.admin-surface .v-main {min-width:0;}
.admin-main {min-width:0;padding:20px 12px;}.admin-toolbar {padding:12px;}.admin-editor {padding:16px;}.admin-editor-grid {grid-template-columns:1fr;}.admin-surface h1 {font-size:25px;}}
</style>
