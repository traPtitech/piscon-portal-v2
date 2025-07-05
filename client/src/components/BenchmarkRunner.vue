<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { components } from '@/api/openapi'
import MainButton from '@/components/MainButton.vue'
import { useEnqueueBenchmark } from '@/lib/useServerData'
import { computed, ref, useId } from 'vue'

const id = useId()

type Bench = components['schemas']['BenchmarkListItem']
type Instance = components['schemas']['Instance']
const props = defineProps<{
  teamId: string
  benches: Bench[]
  instances: Instance[]
}>()

const { mutate: enqueueRequest, isPending } = useEnqueueBenchmark()

const instances = computed(() =>
  props.instances.map((i) => ({
    label: `サーバー${i.serverId} (${i.privateIPAddress})`,
    value: i.id,
  })),
)
const canEnqueue = computed(() => props.benches.every((b) => b.status === 'finished'))

const targetInstanceId = ref<string | null>(null)

const enqueueBenchmark = (instanceId: string | null) => {
  if (instanceId === null) return
  enqueueRequest({ teamId: props.teamId, instanceId })
}
</script>

<template>
  <div class="benchmark-runner-container">
    <div class="benchmark-runner-title">新しくベンチマークを実行する</div>
    <div class="benchmark-runner-content">
      <div>
        <Select v-model="targetInstanceId" :options="instances" option-label="label" option-value="value"
          :placeholder="'対象サーバーを選択'" :id="id" class="benchmark-runner-instance-selector" />
      </div>
      <MainButton :disabled="targetInstanceId === null || isPending || !canEnqueue" variant="primary"
        @click="enqueueBenchmark(targetInstanceId)" class="benchmark-runner-button">
        <Icon icon="mdi:thunder" width="24" height="24" />
        <span>実行</span>
      </MainButton>
    </div>
    <div class="benchmark-runner-messages">
      <div v-if="!canEnqueue" class="benchmark-runner-no-available-slot-message">
        <Icon icon="mdi:alert" width="20" height="20" />
        <span>現在実行中のベンチマークがあるため、新しいベンチマークを実行できません。</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.benchmark-runner-container {
  padding: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
}

.benchmark-runner-title {
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.benchmark-runner-content {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 0.5rem;
}

.benchmark-runner-instance-selector-label {
  display: block;
  font-size: 0.9rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.benchmark-runner-instance-selector {
  width: 100%;
}

.benchmark-runner-button {
  align-self: center;
}

.benchmark-runner-messages {
  margin-top: 0.5rem;
}

.benchmark-runner-no-available-slot-message {
  color: var(--ct-red-600);
  font-size: 0.9rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}
</style>
