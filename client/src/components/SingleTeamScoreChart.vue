<script setup lang="ts">
import type { components } from '@/api/openapi'
import { computed } from 'vue'

type Bench = components['schemas']['BenchmarkListItem']
const props = defineProps<{
  benches: Bench[]
}>()

const teamScores = computed(() => {
  if (!props.benches) return []
  return props.benches
    .filter((b) => b.status === 'finished')
    .flatMap((b) => ({
      teamId: b.teamId,
      score: b.score,
      createdAt: b.finishedAt,
    }))
})
</script>

<template>
  <div class="single-team-score-chart-container">
    <div class="single-team-score-chart-title">スコアの推移</div>
    <div class="score-chart-container">
      <ScoreChart :scores="teamScores" class="team-scores-chart" />
    </div>
  </div>
</template>

<style scoped>
.single-team-score-chart-container {
  height: 200px;
  min-width: 0;
  padding: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  display: flex;
  flex-direction: column;
}

.single-team-score-chart-title {
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.score-chart-container {
  flex: 1;
  position: relative;
  min-height: 0;
}
</style>
