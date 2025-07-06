import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice';
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'

const app = createApp(App)

if (import.meta.env.DEV) {
  const { worker } = await import('@/mock/browser')
  await worker.start({ onUnhandledRequest: 'bypass' })
}

const preset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{cyan.50}',
      100: '{cyan.100}',
      200: '{cyan.200}',
      300: '{cyan.300}',
      400: '{cyan.400}',
      500: '{cyan.500}',
      600: '{cyan.600}',
      700: '{cyan.700}',
      800: '{cyan.800}',
      900: '{cyan.900}',
      950: '{cyan.950}',
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
app.use(ConfirmationService);

app.mount('#app')
