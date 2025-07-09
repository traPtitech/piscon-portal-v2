<script setup lang="ts">
import { apiBaseUrl } from '@/api'
import MainLayout from '@/layouts/MainLayout.vue'
import { useMe } from '@/lib/useServerData'
import { computed, effect } from 'vue'
import { RouterView, useRouter } from 'vue-router'

const { data: me, isLoading } = useMe()
const router = useRouter()

const status = computed(() => {
  if (isLoading.value) return 'loading'
  return me.value ? 'authenticated' : 'unauthenticated'
})

const loginHandler = () => {
  sessionStorage.setItem('redirect', window.location.pathname)
  window.location.href = apiBaseUrl + '/oauth2/code'
}

effect(() => {
  const redirect = sessionStorage.getItem('redirect')
  if (redirect && status.value === 'authenticated') {
    sessionStorage.removeItem('redirect')
    if (!redirect.startsWith('/')) return
    void router.push(redirect)
  }
})
</script>

<template>
  <MainLayout>
    <div v-if="status === 'loading'" class="loading-message">
      <LoaderAnimation />
      <div>読み込み中...</div>
    </div>
    <div v-if="status === 'unauthenticated'" class="unauthenticated-message">
      <div class="portal-title">PISCON Portal v2</div>
      <MainButton @click="loginHandler">traQでログイン</MainButton>
      <div class="caption-message">利用するにはログインが必要です</div>
    </div>
    <RouterView v-if="status === 'authenticated'" :key="$route.fullPath" />
  </MainLayout>
</template>

<style scoped>
.loading-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 11rem);
  gap: 1rem;
}
.unauthenticated-message {
  display: grid;
  place-content: center;
  height: calc(100vh - 11rem);
  gap: 1rem;
}

.portal-title {
  font-size: 2rem;
  font-weight: bold;
}

.caption-message {
  color: var(--ct-slate-500);
  font-size: 0.875rem;
  text-align: center;
  margin-top: -0.5rem;
}
</style>
