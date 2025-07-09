import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'

const app = createApp(App)

if (import.meta.env.DEV && !import.meta.env.VITE_DISABLE_MSW) {
  const { worker } = await import('@/mock/browser')
  await worker.start({ onUnhandledRequest: 'bypass' })
}

const preset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#e6eef6',
      100: '#c2d8ec',
      200: '#9ac0e1',
      300: '#6fa7d6',
      400: '#4b93ce',
      500: '#005bac',
      600: '#004a9a',
      700: '#003c7c',
      800: '#002d5c',
      900: '#001e3c',
      950: '#00122a',
    },
  },
})

app.use(router)
app.use(VueQueryPlugin)
app.use(PrimeVue, {
  theme: {
    preset: preset,
    options: { darkModeSelector: '.dark' },
  },
})
app.use(ConfirmationService)

app.mount('#app')
