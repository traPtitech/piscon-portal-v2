<script setup lang="ts">
import BenchmarkStatusChip from '@/components/BenchmarkStatusChip.vue'
import { formatDate, formatRelativeDate } from '@/lib/formatDate'
import { useTeam } from '@/lib/useServerData'
import { computed } from 'vue'
import { formatScore } from '@/lib/formatScore'
import { useUsers } from '@/lib/useUsers'
import type { components } from '@/api/openapi'

type Bench =
  | (components['schemas']['Benchmark'] & { adminLog?: string })
  | components['schemas']['BenchmarkAdminResult']
type Instance = components['schemas']['Instance']
const { bench, instances } = defineProps<{ bench: Bench; instances: Instance[] }>()

const { data: team } = useTeam(bench.teamId)
const { getUserById } = useUsers()

const instance = computed(() => instances?.find((i) => i.id === bench?.instanceId))
const user = computed(() => getUserById(bench?.userId))
</script>

<template>
  <div v-if="bench" class="bench-detail-container">
    <div class="bench-score-container">
      <div class="bench-score-element">
        <div class="bench-score-label">スコア</div>
        <div
          v-if="bench.status === 'running' || bench.status === 'finished'"
          class="bench-score-content"
        >
          {{ formatScore(bench.score) }}
        </div>
        <div v-else class="bench-score-content-dimmed">未計測</div>
      </div>
      <div class="bench-status">
        <BenchmarkStatusChip :status="bench.status" v-if="bench.status !== 'finished'" />
        <BenchmarkResultChip :result="bench.result" v-else />
      </div>
    </div>
    <div class="bench-detail-element-container">
      <div class="bench-detail-element">
        <div class="bench-detail-label">チーム</div>
        <div class="bench-detail-content">{{ team?.name }}</div>
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
          <time
            :datetime="bench.createdAt"
            :title="formatDate(bench.createdAt, 'YYYY/MM/DD hh:mm:ss.SSS')"
          >
            {{ formatRelativeDate(bench.createdAt) }}
          </time>
        </div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">開始時刻</div>
        <div
          class="bench-detail-content"
          v-if="bench.status === 'running' || bench.status === 'finished'"
        >
          <time
            :datetime="bench.startedAt"
            :title="formatDate(bench.startedAt, 'YYYY/MM/DD hh:mm:ss.SSS')"
          >
            {{ formatRelativeDate(bench.startedAt) }}
          </time>
        </div>
        <div class="bench-detail-content-dimmed" v-else>まだ開始していません</div>
      </div>
      <div class="bench-detail-element">
        <div class="bench-detail-label">終了時刻</div>
        <div class="bench-detail-content" v-if="bench.status === 'finished'">
          <time
            :datetime="bench.finishedAt"
            :title="formatDate(bench.finishedAt, 'YYYY/MM/DD hh:mm:ss.SSS')"
          >
            {{ formatRelativeDate(bench.finishedAt) }}
          </time>
        </div>
        <div class="bench-detail-content-dimmed" v-else>まだ終了していません</div>
      </div>
    </div>
    <div class="bench-log-container">
      <div class="bench-log-label">ベンチマーカーログ</div>
      <pre><code>{{ bench.log }}</code></pre>
    </div>
    <div class="bench-log-container error" v-if="bench.adminLog !== undefined">
      <div class="bench-log-label">ベンチマーカーエラーログ</div>
      <pre><code>{{ bench.adminLog }}</code></pre>
    </div>
  </div>
</template>

<style scoped>
.bench-detail-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  container-type: inline-size;
}

.bench-score-container {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 0.5rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 2px;
  padding: 0.5rem;
}
.bench-score-element {
  display: flex;
  flex-direction: column;
}
.bench-status {
  display: flex;
  align-items: center;
  justify-content: flex-end;
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

@container (max-width: 600px) {
  .bench-detail-element-container {
    grid-template-columns: repeat(2, 1fr);
  }
}

@container (max-width: 400px) {
  .bench-detail-element-container {
    grid-template-columns: 1fr;
  }
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
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

@container (max-width: 600px) {
  .bench-log-container {
    font-size: 0.7rem;
  }
}

.bench-log-label {
  font-weight: 600;
  font-size: 0.9rem;
}

.bench-log-container.error {
  color: var(--ct-red-300);
}
</style>
