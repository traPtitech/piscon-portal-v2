<script setup lang="ts">
import PageTitle from '@/components/PageTitle.vue'
import TeamBenchmarkList from '@/components/TeamBenchmarkList.vue'
import { useMe } from '@/lib/useServerData'

const { data: me } = useMe()
</script>

<template>
  <main class="benches-container">
    <PageTitle icon="mdi:thunder">ベンチマーク</PageTitle>
    <div v-if="me?.teamId !== undefined">
      <TeamBenchmarkList :teamId="me.teamId" />
    </div>
    <div v-if="me !== undefined && me?.teamId === undefined">
      <p>チームに所属していません</p>
    </div>
  </main>
</template>

<style scoped>
.benches-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
