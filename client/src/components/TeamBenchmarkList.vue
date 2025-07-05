<script setup lang="ts">
import BenchmarkList from '@/components/BenchmarkList.vue'
import { useTeamBenches, useTeamInstances } from '@/lib/useServerData'
import { onMounted, onUnmounted } from 'vue'

const { teamId } = defineProps<{ teamId: string }>()

const { data: benches, error: benchesError, refetch } = useTeamBenches(teamId)
const { data: instances } = useTeamInstances(teamId)

let interval: number
onMounted(() => {
  interval = setInterval(() => {
    const loadingStatuses = ['waiting', 'running']
    if (benches.value?.some((b) => loadingStatuses.includes(b.status))) {
      void refetch()
    }
  }, 500)
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>

<template>
  <ErrorMessage v-if="benchesError" :error="benchesError" />
  <BenchmarkList v-else :benches="benches ?? []" :instances="instances ?? []" />
</template>
