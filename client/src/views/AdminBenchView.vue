<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import BenchmarkDetail from '@/components/BenchmarkDetail.vue'
import { useAllInstances, useBench } from '@/lib/useServerData'

const { params } = useRoute()
const benchId = params.id as string

const { data: bench } = useBench(benchId)
const { data: instances } = useAllInstances()
</script>

<template>
  <main class="bench-container">
    <div>
      <NavigationLink to="/admin/benches" class="back-button">
        <Icon icon="mdi:chevron-left" width="24" height="24" />
        <span>ベンチマーク一覧に戻る</span>
      </NavigationLink>
    </div>
    <BenchmarkDetail v-if="bench !== undefined" :bench="bench" :instances="instances ?? []" />
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
}
.back-button svg {
  margin-top: 0.15rem;
}
</style>
