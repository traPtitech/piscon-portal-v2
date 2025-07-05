<script setup lang="ts">
import 'chartjs-adapter-dayjs-4/dist/chartjs-adapter-dayjs-4.esm'
import type { Chart, ChartConfiguration } from 'chart.js'
import { onMounted, onUnmounted, useTemplateRef, watch } from 'vue'

const props = defineProps<{
  config: ChartConfiguration<'line', { x: string; y: number }[], unknown>
}>()

let chart: Chart<'line', { x: string; y: number }[], unknown> | null = null
const canvasRef = useTemplateRef('canvas')
const wrapperRef = useTemplateRef('wrapper')

const resize = () => {
  if (!canvasRef.value || !wrapperRef.value) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const parent = wrapperRef.value.parentElement!
  const rect = parent.getBoundingClientRect()
  const style = window.getComputedStyle(parent)
  const width =
    rect.width -
    parseFloat(style.paddingLeft) -
    parseFloat(style.paddingRight) -
    parseFloat(style.borderLeftWidth) -
    parseFloat(style.borderRightWidth)
  const height =
    rect.height -
    parseFloat(style.paddingTop) -
    parseFloat(style.paddingBottom) -
    parseFloat(style.borderTopWidth) -
    parseFloat(style.borderBottomWidth)
  wrapperRef.value.style.width = `${width}px`
  wrapperRef.value.style.height = `${height}px`
  canvasRef.value.width = width
  canvasRef.value.height = height
  canvasRef.value.style.width = `${width}px`
  canvasRef.value.style.height = `${height}px`

  if (chart) {
    chart.resize()
    chart.config.options!.aspectRatio = width / height
  }
}

onMounted(async () => {
  resize()
  window.addEventListener('resize', resize)

  if (!canvasRef.value) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const { Chart, LineController, LineElement, LinearScale, PointElement, TimeScale, Tooltip } =
    await import('chart.js')

  Chart.register(LineController, LineElement, LinearScale, PointElement, TimeScale, Tooltip)

  chart = new Chart(ctx, structuredClone(props.config))
  chart.update()
})

onUnmounted(() => {
  window.removeEventListener('resize', resize)
  if (chart) {
    chart.destroy()
  }
})

watch(
  () => props.config.data,
  (newData) => {
    if (chart === null) return
    if (
      chart.data.datasets.length === newData.datasets.length &&
      chart.data.datasets.every(
        (dataset, index) =>
          dataset.label === newData.datasets[index].label &&
          JSON.stringify(dataset.data) === JSON.stringify(newData.datasets[index].data),
      )
    ) {
      return
    }

    console.log('Updating chart data', newData)

    chart.data.datasets = newData.datasets
    chart.update()
  },
  { deep: true },
)
</script>

<template>
  <div ref="wrapper">
    <canvas ref="canvas"></canvas>
  </div>
</template>
