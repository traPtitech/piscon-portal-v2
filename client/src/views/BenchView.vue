<script setup lang="ts">
import { useMe } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import TeamBenchmarkDetail from '@/components/TeamBenchmarkDetail.vue'

const { data: me } = useMe()
const { params } = useRoute()
const benchId = params.id as string
</script>

<template>
  <main class="bench-container">
    <div>
      <NavigationLink to="/benches" class="back-button">
        <Icon icon="mdi:chevron-left" width="24" height="24" />
        <span>ベンチマーク一覧に戻る</span>
      </NavigationLink>
    </div>
    <div v-if="me?.teamId !== undefined">
      <TeamBenchmarkDetail :teamId="me.teamId" :benchId="benchId" />
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
}
.back-button svg {
  margin-top: 0.15rem;
}
</style>
