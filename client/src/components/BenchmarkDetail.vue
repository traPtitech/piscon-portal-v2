<script setup lang="ts">
import BenchmarkStatusChip from '@/components/BenchmarkStatusChip.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import { formatDate } from '@/lib/formatDate'
import { useTeamBench, useTeamInstances, useUsers } from '@/lib/useServerData'
import { computed, watch } from 'vue'
import { formatScore } from '@/lib/formatScore'

const { teamId, benchId } = defineProps<{
  teamId: string
  benchId: string
}>()

const { data: bench, error: benchError, refetch } = useTeamBench(teamId, benchId)
const { data: instances } = useTeamInstances(teamId)
const { data: users } = useUsers()

const instance = computed(() => instances.value?.find((i) => i.id === bench.value?.instanceId))
const user = computed(() => users.value?.find((u) => u.id === bench.value?.userId))

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
  <ErrorMessage v-if="benchError"></ErrorMessage>
  <div v-if="bench" class="bench-detail-container">
    <div class="bench-score-container">
      <div class="bench-score-label">スコア</div>
      <div
        v-if="bench.status === 'running' || bench.status === 'finished'"
        class="bench-score-content"
      >
        {{ formatScore(bench.score) }}
      </div>
      <div v-else class="bench-score-content-dimmed">未計測</div>
    </div>
    <div class="bench-detail-element-container">
      <div class="bench-detail-element">
        <div class="bench-detail-label">ステータス</div>
        <BenchmarkStatusChip :status="bench.status" />
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">実行ユーザー</div>
        <div class="bench-detail-content">@{{ user?.name }}</div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">対象サーバー</div>
        <div class="bench-detail-content">サーバー{{ instance?.serverId }}</div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">リクエスト時刻</div>
        <div class="bench-detail-content">
          {{ formatDate(bench.createdAt, 'YYYY/MM/DD hh:mm:ss') }}
        </div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">開始時刻</div>
        <div
          class="bench-detail-content"
          v-if="bench.status === 'running' || bench.status === 'finished'"
        >
          {{ formatDate(bench.startedAt, 'YYYY/MM/DD hh:mm:ss') }}
        </div>
        <div class="bench-detail-content-dimmed" v-else>まだ開始していません</div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">終了時刻</div>
        <div class="bench-detail-content" v-if="bench.status === 'finished'">
          {{ formatDate(bench.finishedAt, 'YYYY/MM/DD hh:mm:ss') }}
        </div>
        <div class="bench-detail-content-dimmed" v-else>まだ終了していません</div>
      </div>
    </div>
    <div class="bench-log-container">
      <pre><code>{{ bench.log }}</code></pre>
    </div>
  </div>
</template>

<style scoped>
.bench-detail-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.bench-score-container {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--ct-slate-300);
  border-radius: 2px;
  padding: 0.5rem;
}
.bench-score-label {
  font-weight: 700;
  font-size: 0.9rem;
}
.bench-score-content {
  font-size: 1.5rem;
  font-weight: 700;
}
.bench-score-content-dimmed {
  font-size: 1.5rem;
  color: var(--ct-slate-500);
}

.bench-detail-element-container {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.5rem;
}
.bench-detail-element {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--ct-slate-300);
  border-radius: 2px;
  padding: 0.5rem;
}
.bench-detail-label {
  font-weight: 700;
  font-size: 0.8rem;
}
.bench-detail-content {
  font-size: 0.9rem;
}
.bench-detail-content-dimmed {
  font-size: 0.9rem;
  color: var(--ct-slate-500);
}

.bench-log-container {
  border-radius: 2px;
  background-color: var(--ct-slate-800);
  color: var(--ct-slate-200);
  font-size: 0.8rem;
  padding: 1rem;
  overflow-x: auto;
}
</style>
