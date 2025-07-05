<script setup lang="ts">
import { formatDate, formatRelativeDate } from '@/lib/formatDate'
import { Icon } from '@iconify/vue'
import BenchmarkStatusChip from '@/components/BenchmarkStatusChip.vue'
import { formatScore } from '@/lib/formatScore'
import { useUsers } from '@/lib/useUsers'
import type { components } from '@/api/openapi'
import { computed } from 'vue'
import { useTeams } from '@/lib/useServerData'

type Bench = components['schemas']['BenchmarkListItem']
type Instance = components['schemas']['Instance']
const { benches, instances, isAdmin } = defineProps<{
  benches: Bench[]
  instances: Instance[]
  isAdmin?: boolean
}>()

const sortedBenches = computed(() =>
  [...(benches ?? [])].sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
  ),
)

const { getUserById } = useUsers()
const { data: teams } = useTeams()

const columns = computed(() => (isAdmin ? 7 : 6))

const getTeamName = (teamId: string) => teams.value?.find((t) => t.id === teamId)?.name ?? ''
const getInstanceServerId = (instanceId: string) =>
  instances.find((i) => i.id === instanceId)?.serverId ?? '?'
const getUserName = (userId: string) => getUserById(userId)?.name ?? ''
</script>

<template>
  <div class="bench-list-container">
    <div class="bench-list">
      <div class="list-label">
        <Icon icon="mdi:score" width="24" height="24" />
        <span>スコア</span>
      </div>
      <div class="list-label list-datetime">
        <Icon icon="mdi:calendar-clock" width="24" height="24" />
        <span>リクエスト時刻</span>
      </div>
      <div class="list-label" v-if="isAdmin">
        <Icon icon="mdi:account-group" width="24" height="24" />
        <span>チーム</span>
      </div>
      <div class="list-label list-server">
        <Icon icon="mdi:server-network" width="24" height="24" />
        <span>対象サーバー</span>
      </div>
      <div class="list-label list-user">
        <Icon icon="mdi:account" width="24" height="24" />
        <span>実行ユーザー</span>
      </div>
      <div class="list-label">
        <Icon icon="mdi:progress-clock" width="24" height="24" />
        <span>ステータス</span>
      </div>
      <div class="list-label"></div>
      <template v-for="bench in sortedBenches" :key="bench.id">
        <div v-if="bench.status === 'running' || bench.status === 'finished'" class="bench-score">
          {{ formatScore(bench.score) }}
        </div>
        <div v-else class="bench-score-loading">待機中</div>
        <div class="bench-date list-datetime">
          <time :datetime="bench.createdAt" :title="formatDate(bench.createdAt, 'YYYY/MM/DD hh:mm:ss.SSS')">
            {{ formatRelativeDate(bench.createdAt) }}
          </time>
        </div>
        <div v-if="isAdmin" class="bench-team">
          {{ getTeamName(bench.teamId) }}
        </div>
        <div class="bench-server list-server">
          サーバー{{ getInstanceServerId(bench.instanceId) }}
        </div>
        <div class="bench-user list-user">
          <UserChip :name="getUserName(bench.userId)" />
        </div>
        <div>
          <BenchmarkStatusChip :status="bench.status" />
        </div>
        <div>
          <NavigationLink :to="isAdmin ? `/admin/benches/${bench.id}` : `/benches/${bench.id}`" class="bench-link">
            <span>詳細を見る</span>
            <Icon icon="mdi:chevron-right" width="24" height="24" />
          </NavigationLink>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.bench-list-container {
  width: 100%;
  container-type: inline-size;
}

.bench-list {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(v-bind(columns), auto);
}

.bench-skeleton {
  height: 2rem;
  border-bottom: 1px solid var(--ct-slate-300);
  padding: 0.5rem 1rem;
  display: flex;
  align-items: center;
}

.bench-list>div {
  border-bottom: 1px solid var(--ct-slate-300);
  padding: 0.5rem 1rem;
  display: flex;
  align-items: center;
}

.list-label {
  font-weight: 700;
  gap: 0.25rem;
}

.bench-score {
  font-weight: 700;
}

.bench-score-loading {
  color: var(--ct-slate-500);
  font-size: 0.8rem;
}

.bench-date {
  font-size: 0.9rem;
}

.bench-server {
  font-size: 0.9rem;
}

.bench-user {
  font-weight: 600;
  font-size: 0.9rem;
  gap: 0.25rem;
}

.bench-link {
  display: flex;
  width: fit-content;
  align-items: center;
  gap: 0rem;
}

.bench-link svg {
  margin-top: 0.15rem;
}

@container (max-width: 900px) {
  .bench-list {
    grid-template-columns: repeat(v-bind(columns-1), auto);
  }

  .list-datetime.list-datetime {
    display: none;
  }
}

@container (max-width: 780px) {
  .bench-list {
    grid-template-columns: repeat(v-bind(columns-2), auto);
  }

  .list-server.list-server {
    display: none;
  }
}

@container (max-width: 560px) {
  .bench-list {
    grid-template-columns: repeat(v-bind(columns-3), auto);
  }

  .list-user.list-user {
    display: none;
  }
}
</style>
