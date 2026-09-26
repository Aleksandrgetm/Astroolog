import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const site = process.env.VITE_SITE_URL || loadEnv(mode, process.cwd(), '').VITE_SITE_URL
  if (mode === 'production') {
    let valid = false
    try { const url = new URL(site || ''); valid = url.protocol === 'https:' && url.origin === site?.replace(/\/$/, '') } catch {}
    if (!valid) throw new Error('Production build requires VITE_SITE_URL as an absolute HTTPS origin.')
  }
  return {
  plugins: [
    vue(),
  ],

  server: {
    host: '0.0.0.0',

    allowedHosts: [
      '.trycloudflare.com',
    ],

    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
}})