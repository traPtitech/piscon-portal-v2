<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useBenchmarkQueue } from '@/lib/useServerData'
import { useInterval } from '@/lib/useInterval'
import type { components } from '@/api/openapi'
import { ref, watch } from 'vue'
import { useTeams } from '@/lib/useTeams'

const { data: queue, refetch } = useBenchmarkQueue()
const { getTeamName } = useTeams()

type QueueItem =
  | components['schemas']['BenchmarkListItem']
  | (Omit<components['schemas']['RunningBenchmark'], 'status'> & { status: 'removed' })
const queueItems = ref<QueueItem[] | undefined>(undefined)

watch(
  queue,
  (newQueue, prevQueue) => {
    if (newQueue === undefined) return
    if (prevQueue === undefined) {
      queueItems.value = newQueue
      return
    }

    const addedItems = newQueue.filter((item) => {
      return !prevQueue.some((prevItem) => prevItem.id === item.id)
    })
    const removedItems = prevQueue.filter((item) => {
      return !newQueue.some((newItem) => newItem.id === item.id)
    })

    queueItems.value = [
      ...(queueItems.value?.map(
        (item) =>
          newQueue.find((newItem) => newItem.id === item.id) ??
          ({ ...(item as components['schemas']['RunningBenchmark']), status: 'removed' } as const),
      ) ?? []),
      ...addedItems,
    ]

    setTimeout(() => {
      queueItems.value = queueItems.value?.filter((item) =>
        removedItems.every((i) => i.id !== item.id),
      )
    }, 600)
  },
  { immediate: true },
)

useInterval(() => void refetch(), 500)
</script>

<template>
  <div class="benchmark-queue-container">
    <div v-if="queueItems !== undefined" class="benchmark-queue-items">
      <div v-if="queueItems.length === 0" class="">ベンチマークキューは空です</div>
      <div
        v-for="item in queueItems"
        :key="item.id"
        class="benchmark-queue-item"
        :class="item.status"
      >
        <div class="benchmark-queue-item-content">
          <div class="benchmark-queue-item-icon">
            <Icon
              :icon="item.status === 'waiting' ? 'mdi:box-variant-closed' : 'mdi:box-variant'"
              width="48"
              height="48"
              class="benchmark-queue-icon"
            />
            <div class="benchmark-queue-side-icon">
              <BenchmarkProgressCircle
                v-if="item.status === 'running'"
                :size="24"
                :startedAt="item.startedAt"
              />
              <BenchmarkProgressCircle
                v-else-if="item.status === 'removed'"
                :size="24"
                :startedAt="new Date(Date.now() - 1000 * 60).toISOString()"
              />
            </div>
          </div>
          <div class="benchmark-queue-item-details">
            <div class="benchmark-queue-item-name">{{ getTeamName(item.teamId) }}</div>
            <BenchmarkStatusChip :status="item.status === 'removed' ? 'finished' : item.status" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.benchmark-queue-container {
  padding: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  overflow-x: hidden;
}

.benchmark-queue-item {
  display: grid;
  grid-template-rows: 1fr;
  animation: wipe-in 0.3s ease-in-out forwards;
}

.benchmark-queue-item-icon {
  position: relative;
  width: fit-content;
  height: fit-content;
}

.benchmark-queue-item.removed {
  animation:
    wipe-out 0.3s ease-in-out forwards,
    shrink 0.2s ease-in-out 0.4s forwards;
}

@keyframes wipe-in {
  from {
    opacity: 0;
    transform: translateX(64px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes wipe-out {
  from {
    opacity: 1;
    transform: translateX(0);
  }
  to {
    opacity: 0;
    transform: translateX(64px);
  }
}

@keyframes shrink {
  from {
    grid-template-rows: 1fr;
  }
  to {
    grid-template-rows: 0fr;
  }
}

.benchmark-queue-item-content {
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.benchmark-queue-icon {
  color: var(--ct-slate-500);
}

.benchmark-queue-side-icon {
  position: absolute;
  right: 0;
  bottom: 0;
}

.benchmark-queue-item-details {
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
}
</style>
