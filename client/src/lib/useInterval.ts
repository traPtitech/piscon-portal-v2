import { onMounted, onUnmounted } from 'vue'

export const useInterval = (callback: () => void, delay: number) => {
  let interval: number
  onMounted(() => {
    interval = setInterval(callback, delay)
  })

  onUnmounted(() => {
    clearInterval(interval)
  })
}
