import axios from 'axios'
import { reactive } from 'vue'
export const adminSession=reactive<{user:null|{id:number;email:string};csrf:string}>({user:null,csrf:''})
export const adminApi=axios.create({baseURL:'/api/admin',withCredentials:true,timeout:30000})
adminApi.interceptors.request.use(config=>{if(config.method!=='get'&&adminSession.csrf)config.headers['X-CSRF-Token']=adminSession.csrf;return config})
export async function loadSession(){try{const {data}=await adminApi.get('/auth/me');Object.assign(adminSession,data);return true}catch{adminSession.user=null;adminSession.csrf='';return false}}
export function adminError(e:unknown){return axios.isAxiosError(e)?e.response?.data?.error||'Не удалось выполнить запрос. Попробуйте ещё раз.':'Не удалось выполнить запрос.'}
