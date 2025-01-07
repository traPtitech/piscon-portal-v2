<script setup lang="ts">
import BenchmarkList from '@/components/BenchmarkList.vue'
import { useMe } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'

const { data: me } = useMe()
</script>

<template>
  <main class="benches-container">
    <h1 class="title">
      <Icon icon="mdi:thunder" width="32" height="32" />
      <span>ベンチマーク</span>
    </h1>
    <div v-if="me?.teamId !== undefined">
      <BenchmarkList :teamId="me.teamId" />
    </div>
    <div v-if="me !== undefined && me?.teamId === undefined">
      <p>チームに所属していません</p>
    </div>
  </main>
</template>

<style scoped>
.title {
  display: flex;
  align-items: center;
  font-size: 1.5rem;
  margin-bottom: 2rem;
}

.benches-container {
  padding: 1rem;
}
</style>
