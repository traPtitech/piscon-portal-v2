<script setup lang="ts">
import { useMe } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import { RouterLink, useRoute } from 'vue-router'
import BenchmarkDetail from '@/components/BenchmarkDetail.vue'

const { data: me } = useMe()
const { params } = useRoute()
const benchId = params.id as string
</script>

<template>
  <main class="bench-container">
    <div>
      <RouterLink to="/benches" class="back-button">
        <Icon icon="mdi:chevron-left" width="24" height="24" />
        <span>ベンチマーク一覧に戻る</span>
      </RouterLink>
    </div>
    <div v-if="me?.teamId !== undefined">
      <BenchmarkDetail :teamId="me.teamId" :benchId="benchId" />
    </div>
    <div v-if="me !== undefined && me?.teamId === undefined">
      <p>チームに所属していません</p>
    </div>
  </main>
</template>

<style scoped>
.bench-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.back-button {
  width: fit-content;
  padding: 0.25rem 0.5rem;
  display: flex;
  align-items: center;
  color: var(--color-primary);
  text-decoration: none;
  font-size: 1rem;
  font-weight: 600;
}
.back-button svg {
  margin-top: 0.15rem;
}
</style>
