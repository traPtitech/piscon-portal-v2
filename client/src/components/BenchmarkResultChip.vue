<script setup lang="ts">
import type { components } from '@/api/openapi'
type Result = components['schemas']['FinishedBenchmark']['result']
const { result } = defineProps<{ result: Result }>()

const labels: Record<Result, string> = {
  passed: 'PASSED',
  failed: 'FAILED',
  error: 'ERROR',
}
</script>

<template>
  <div class="result-chip" :class="result">
    {{ labels[result] }}
  </div>
</template>

<style scoped>
.result-chip {
  line-height: 2;
  border-radius: 2px;
  width: 72px;
  font-weight: 600;
  font-size: 0.7rem;
  text-align: center;
}
.result-chip.passed {
  background-color: #5cb85c33;
  color: #5cb85c;
}
.result-chip.failed {
  background-color: #e52b2b33;
  color: #e52b2b;
}
.result-chip.error {
  background-color: #ddbf0033;
  color: #bca408;
}
</style>
