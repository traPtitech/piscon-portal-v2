import { fileURLToPath, URL } from 'node:url'

import { defineConfig, UserConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import Components from 'unplugin-vue-components/vite'
import { PrimeVueResolver } from '@primevue/auto-import-resolver'

const mode = process.env.NODE_ENV || 'development'
process.env = { ...process.env, ...loadEnv(mode, process.cwd()) }

const apiProxy: NonNullable<UserConfig['server']>['proxy'] = {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    Components({
      resolvers: [PrimeVueResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: process.env.VITE_DISABLE_MSW ? apiProxy : {},
  },
})
