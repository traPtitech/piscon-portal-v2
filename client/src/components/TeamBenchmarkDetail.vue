<script setup lang="ts">
import { useTeamBench, useTeamInstances } from '@/lib/useServerData'
import { watch } from 'vue'

const { teamId, benchId } = defineProps<{
  teamId: string
  benchId: string
}>()

const { data: bench, error, refetch } = useTeamBench(teamId, benchId)
const { data: instances } = useTeamInstances(teamId)

watch(
  bench,
  () => {
    if (bench.value?.status === 'running' || bench.value?.status === 'waiting') {
      const interval = setInterval(() => {
        refetch()
      }, 1000)
      return () => clearInterval(interval)
    }
  },
  { immediate: true },
)
</script>

<template>
  <ErrorMessage v-if="error" />
  <BenchmarkDetail v-if="bench !== undefined" :bench="bench" :instances="instances ?? []" />
</template>
