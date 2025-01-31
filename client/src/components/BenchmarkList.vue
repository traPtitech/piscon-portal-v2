<script setup lang="ts">
import { useTeamBenches, useTeamInstances, useUsers } from '@/lib/useServerData'
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { formatDate } from '@/lib/formatDate'
import { Icon } from '@iconify/vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import BenchmarkStatusChip from '@/components/BenchmarkStatusChip.vue'
import { formatScore } from '@/lib/formatScore'
import UserAvatar from '@/components/UserAvatar.vue'

const { teamId } = defineProps<{ teamId: string }>()

const { data: benches, error: benchesError } = useTeamBenches(teamId)
const sortedBenches = computed(() =>
  [...(benches.value ?? [])].sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
  ),
)
const { data: instances } = useTeamInstances(teamId)
const { data: users } = useUsers()
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
        <span>日時</span>
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
        <div v-else class="bench-score-loading">計測中</div>
        <div class="bench-date list-datetime">
          {{ formatDate(bench.createdAt, 'YYYY/MM/DD hh:mm:ss.SSS') }}
        </div>
        <div class="bench-server list-server">
          サーバー{{ instances?.find((i) => i.id === bench.instanceId)?.serverId ?? '?' }}
        </div>
        <div class="bench-user list-user">
          <UserAvatar :name="users?.find((u) => u.id === bench.userId)?.name ?? ''" />
          <span>@{{ users?.find((u) => u.id === bench.userId)?.name }}</span>
        </div>
        <div>
          <BenchmarkStatusChip :status="bench.status" />
        </div>
        <div>
          <RouterLink :to="`/benches/${bench.id}`" class="bench-link">
            <span>詳細を見る</span>
            <Icon icon="mdi:chevron-right" width="24" height="24" />
          </RouterLink>
        </div>
      </template>
    </div>
    <ErrorMessage v-if="benchesError" />
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
  grid-template-columns: repeat(6, auto);
}

.bench-skeleton {
  height: 2rem;
  border-bottom: 1px solid var(--ct-slate-300);
  padding: 0.5rem 1rem;
  display: flex;
  align-items: center;
}

.bench-list > div {
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
  text-decoration: none;
  color: var(--color-primary);
  font-weight: 700;
}

.bench-link svg {
  margin-top: 0.15rem;
}

@container (max-width: 900px) {
  .bench-list {
    grid-template-columns: repeat(5, auto);
  }
  .list-datetime.list-datetime {
    display: none;
  }
}

@container (max-width: 780px) {
  .bench-list {
    grid-template-columns: repeat(4, auto);
  }
  .list-server.list-server {
    display: none;
  }
}

@container (max-width: 560px) {
  .bench-list {
    grid-template-columns: repeat(3, auto);
  }
  .list-user.list-user {
    display: none;
  }
}
</style>
