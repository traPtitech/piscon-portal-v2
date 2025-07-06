<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { components } from '@/api/openapi'
import { useTeams } from '@/lib/useServerData'

type RankingItem = components['schemas']['RankingItem']

defineProps<{
  ranking: RankingItem[]
  highlightTeamId?: string
}>()

const { data: teams } = useTeams()

const getTeamName = (teamId: string): string => {
  const team = teams.value?.find((t) => t.id === teamId)
  return team ? team.name : `Unknown Team`
}
</script>

<template>
  <div class="ranking-list">
    <div class="ranking-list-header">
      <div>順位</div>
      <div>チーム</div>
      <div>スコア</div>
    </div>
    <div class="ranking-list-body">
      <div v-for="item in ranking" :key="item.teamId"
        :class="['ranking-list-row', { highlight: item.teamId === highlightTeamId }]">
        <div class="ranking-list-rank">
          <span v-if="item.rank <= 3" class="crown">
            <Icon icon="mdi:crown" width="24" height="24" :title="`${item.rank}`" :style="{
              color: ['orange', 'silver', 'indianred'][item.rank - 1] || 'var(--ct-slate-900)'
            }" />
          </span>
          <span v-else>
            {{ item.rank }}
          </span>
        </div>
        <div class="ranking-list-team">{{ getTeamName(item.teamId) }}</div>
        <div class="ranking-list-score">{{ item.score }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ranking-list {
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  overflow: auto;
  background: var(--ct-white);
}

.ranking-list-header {
  display: grid;
  grid-template-columns: 40px 1fr 120px;
  gap: 1rem;
  padding: 0.75rem 1rem;
  background: var(--ct-slate-100);
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--ct-slate-700);
}

.ranking-list-body {
  max-height: 400px;
  overflow-y: auto;
}

.ranking-list-row {
  display: grid;
  grid-template-columns: 40px 1fr 120px;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--ct-slate-200);
  transition: background-color 0.1s;
}

.ranking-list-row:last-child {
  border-bottom: none;
}

.ranking-list-row:hover {
  background: var(--ct-slate-50);
}

.ranking-list-row.highlight {
  background: rgba(var(--color-primary-rgb), 0.1);
  border-left: 4px solid var(--color-primary);
  padding-left: calc(1rem - 4px);
}

.ranking-list-row.highlight:hover {
  background: rgba(var(--color-primary-rgb), 0.15);
}

.ranking-list-rank {
  font-weight: 600;
  color: var(--ct-slate-900);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.25em;
}

.crown {
  margin-right: 0.2em;
}

.ranking-list-team {
  color: var(--ct-slate-800);
}

.ranking-list-row.highlight .ranking-list-team {
  color: var(--color-primary);
  font-weight: 600;
}

.ranking-list-score {
  font-weight: 600;
  text-align: right;
  color: var(--ct-slate-900);
}
</style>
