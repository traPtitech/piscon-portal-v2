<script setup lang="ts">
import BenchmarkList from '@/components/BenchmarkList.vue'
import { useAllBenches, useAllInstances } from '@/lib/useServerData'

const { data: benches } = useAllBenches()
const { data: instances } = useAllInstances()
</script>

<template>
  <main class="admin-benches-container">
    <PageTitle icon="mdi:thunder-circle">ベンチマーク (管理者)</PageTitle>

    <BenchmarkList :benches="benches ?? []" :instances="instances ?? []" isAdmin />
  </main>
</template>

<style scoped>
.admin-benches-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
