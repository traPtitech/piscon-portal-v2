<script setup lang="ts">
import {
  useTeamBench,
  useTeamBenches,
  useEnqueueBenchmark,
  useTeamInstances,
} from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import { computed, watch } from 'vue'

const { teamId, benchId } = defineProps<{
  teamId: string
  benchId: string
}>()

const { data: bench, error: benchError, refetch } = useTeamBench(teamId, benchId)
const { data: benches } = useTeamBenches(teamId)
const { mutate: enqueueBenchmark } = useEnqueueBenchmark({ redirect: true })
const { data: instances } = useTeamInstances(teamId)

const canReEnqueue = computed(() => benches.value?.every((b) => b.status === 'finished'))

const reEnqueueBenchmark = () => {
  if (bench.value === undefined) return

  enqueueBenchmark({ teamId, instanceId: bench.value?.instanceId })
}

watch(
  bench,
  () => {
    if (bench.value?.status === 'running' || bench.value?.status === 'waiting') {
      const interval = setInterval(() => {
        void refetch()
      }, 1000)
      return () => clearInterval(interval)
    }
  },
  { immediate: true },
)
</script>

<template>
  <div>
    <div class="bench-detail-actions">
      <MainButton
        @click="reEnqueueBenchmark"
        :disabled="!canReEnqueue"
        v-tooltip.bottom="
          !canReEnqueue
            ? {
                value: '現在実行中のベンチマークがあるため、新しいベンチマークを実行できません',
                pt: {
                  arrow: {
                    style: {
                      borderBottomColor: 'rgba(var(--ct-slate-800-rgb), 0.9)',
                    },
                  },
                  text: {
                    style: {
                      color: 'var(--ct-slate-100)',
                      backgroundColor: 'rgba(var(--ct-slate-800-rgb), 0.9)',
                      fontSize: '0.875rem',
                    },
                  },
                },
              }
            : ''
        "
      >
        <Icon icon="mdi:thunder" width="20" height="20" />
        <span>ベンチマーク再実行</span>
      </MainButton>
    </div>
    <ErrorMessage v-if="benchError" :error="benchError" />
    <BenchmarkDetail v-else-if="bench !== undefined" :bench="bench" :instances="instances ?? []" />
  </div>
</template>

<style scoped>
.bench-detail-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}
</style>
