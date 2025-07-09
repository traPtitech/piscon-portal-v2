<script lang="ts" setup>
import { BENCHMARK_DURATION_SECONDS } from '@/lib/benchmark'
import { computed, onMounted, useTemplateRef } from 'vue'

const duration = `${BENCHMARK_DURATION_SECONDS}s`

const props = defineProps<{
  size: number
  startedAt: string
}>()

const ref = useTemplateRef('progress-circle')
const sizePx = computed(() => `${props.size}px`)
const delay = computed(() => {
  const startedAt = new Date(props.startedAt)
  const now = new Date()
  const elapsedSeconds = Math.max((now.getTime() - startedAt.getTime()) / 1000, 0)
  const delay = Math.min(elapsedSeconds, BENCHMARK_DURATION_SECONDS)
  return `-${delay}s`
})
</script>

<template>
  <svg
    :width="size"
    :height="size"
    :viewBox="`0 0 ${size} ${size}`"
    class="progress-circle"
    ref="progress-circle"
  >
    <circle class="bg"></circle>
    <circle class="fg"></circle>
  </svg>
</template>

<style scoped>
.progress-circle {
  --size: v-bind(sizePx);
  --half-size: calc(var(--size) / 2);
  --stroke-width: calc(var(--size) / 4);
  --radius: calc((var(--size) - var(--stroke-width)) / 2);
  --circumference: calc(var(--radius) * pi * 2);
  --dash: calc((var(--progress) * var(--circumference)) / 100);
  animation: progress-animation v-bind(duration) linear v-bind(delay) 1 forwards;
}
.progress-circle circle {
  cx: var(--half-size);
  cy: var(--half-size);
  r: var(--radius);
  stroke-width: var(--stroke-width);
  fill: none;
  stroke-linecap: round;
}
.progress-circle circle.bg {
  stroke: var(--ct-slate-200);
}
.progress-circle circle.fg {
  transform: rotate(-90deg);
  transform-origin: var(--half-size) var(--half-size);
  stroke-dasharray: var(--dash) calc(var(--circumference) - var(--dash));
  transition: stroke-dasharray 0.3s linear 0s;
  stroke: var(--color-primary);
}
@property --progress {
  syntax: '<number>';
  inherits: false;
  initial-value: 0;
}
@keyframes progress-animation {
  from {
    --progress: 0;
  }
  to {
    --progress: 100;
  }
}
</style>
