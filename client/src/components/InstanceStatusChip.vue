<script setup lang="ts">
import type { components } from '@/api/openapi'
type Status = components['schemas']['InstanceStatus']
const { status } = defineProps<{ status: Status }>()

const labels: Record<Status, string> = {
  building: 'BUILDING',
  starting: 'STARTING',
  running: 'RUNNING',
  stopping: 'STOPPING',
  stopped: 'STOPPED',
  deleting: 'DELETING',
  deleted: 'DELETED',
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
.status-chip.building {
  background-color: #f0ad4e33;
  color: #f0ad4e;
}
.status-chip.starting {
  background-color: #5bc0de33;
  color: #5bc0de;
}
.status-chip.running {
  background-color: #5cb85c33;
  color: #5cb85c;
}
.status-chip.stopping {
  background-color: #f0ad4e33;
  color: #f0ad4e;
}
.status-chip.stopped {
  background-color: #de5b5b33;
  color: #de5b5b;
}
.status-chip.deleting {
  background-color: #f0ad4e33;
  color: #f0ad4e;
}
.status-chip.deleted {
  background-color: #77777733;
  color: #777777;
}
</style>
