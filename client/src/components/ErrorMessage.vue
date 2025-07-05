<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { effect } from 'vue'
const props = defineProps<{
  error?: Error
}>()

effect(() => {
  if (props.error) {
    console.error('ErrorMessage:', props.error)
  }
})
</script>

<template>
  <div class="error-message">
    <Icon icon="mdi:alert-octagon" width="48" height="48" />
    <span>エラーが発生しました</span>
  </div>
</template>

<style scoped>
.error-message {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.2rem;
  color: var(--ct-red-500);
  font-weight: 700;
  text-align: center;
}
</style>
