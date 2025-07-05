<script setup lang="ts">
import { useTeams } from '@/lib/useServerData'
import type { ChartConfiguration } from 'chart.js'
import { computed } from 'vue'

type TeamScore = {
  teamId: string
  score: number
  createdAt: string
}
const props = defineProps<{
  scores: TeamScore[]
}>()

const { data: teams } = useTeams()

const getTeamName = (teamId: string): string => {
  const team = teams.value?.find((t) => t.id === teamId)
  return team ? team.name : `Unknown Team`
}

const scoresByTeam = computed(() => {
  const scoresMap: Record<string, { x: string; y: number }[]> = {}
  props.scores.forEach((s) => {
    scoresMap[s.teamId] ??= []
    scoresMap[s.teamId].push({ x: s.createdAt, y: s.score })
  })
  return scoresMap
})

const colors = [
  '#fa525299',
  '#be4bdb99',
  '#7950f299',
  '#4c6ef599',
  '#15aabf99',
  '#12b88699',
  '#40c05799',
  '#82c91e99',
  '#fab00599',
  '#fd7e1499',
]

const config = computed<ChartConfiguration<'line', { x: string; y: number }[], unknown>>(() => ({
  type: 'line',
  data: {
    datasets: Object.entries(scoresByTeam.value).map(([teamId, scores], index) => ({
      label: getTeamName(teamId),
      data: scores,
      fill: false,
      borderColor: colors[index % colors.length],
      tension: 0.1,
    })),
  },
  options: {
    scales: {
      x: {
        type: 'time',
        time: {
          unit: 'minute',
          tooltipFormat: 'MM/DD HH:mm:ss',
        },
        display: false,
      },
      y: {
        type: 'linear',
        beginAtZero: true,
      },
    },
  },
}))
</script>

<template>
  <LineChart :config="config" />
</template>

<style scoped></style>
