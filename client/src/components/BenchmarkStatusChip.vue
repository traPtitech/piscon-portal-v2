<script setup lang="ts">
import type { components } from '@/api/openapi'
type Status = components['schemas']['BenchmarkStatus']
const { status } = defineProps<{ status: Status }>()

const labels: Record<Status, string> = {
  waiting: 'WAITING',
  readying: 'READY',
  running: 'RUNNING',
  finished: 'FINISHED',
}
</script>

<template>
  <div class="status-chip" :class="status">
    {{ labels[status] }}
  </div>
</template>

<style scoped>
.status-chip {
  line-height: 2;
  border-radius: 2px;
  width: 72px;
  font-weight: 600;
  font-size: 0.7rem;
  text-align: center;
}
.status-chip.waiting {
  background-color: #f0ad4e33;
  color: #f0ad4e;
}
.status-chip.readying {
  background-color: #f0ad4e33;
  color: #f0ad4e;
}
.status-chip.running {
  background-color: #5bc0de33;
  color: #5bc0de;
}
.status-chip.finished {
  background-color: #5cb85c33;
  color: #5cb85c;
}
</style>
