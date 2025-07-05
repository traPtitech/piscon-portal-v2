<script setup lang="ts">
import BenchmarkList from '@/components/BenchmarkList.vue'
import { useTeamBenches, useTeamInstances } from '@/lib/useServerData'

const { teamId } = defineProps<{ teamId: string }>()

const { data: benches, error: benchesError } = useTeamBenches(teamId)
const { data: instances } = useTeamInstances(teamId)
</script>

<template>
  <ErrorMessage v-if="benchesError" :error="benchesError" />
  <BenchmarkList v-else :benches="benches ?? []" :instances="instances ?? []" />
</template>
