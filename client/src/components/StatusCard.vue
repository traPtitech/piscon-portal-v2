<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  icon: string
  title: string
  value: string | number
  description?: string
  variant?: 'default' | 'primary' | 'highlight'
}>()
</script>

<template>
  <div :class="['status-card', variant ?? 'default']">
    <div class="status-card-header">
      <Icon :icon="icon" width="24" height="24" />
      <span class="status-card-title">{{ title }}</span>
    </div>
    <div class="status-card-value">{{ value }}</div>
  </div>
</template>

<style scoped>
.status-card {
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--ct-slate-300);
  background: var(--ct-white);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.status-card.primary {
  border-color: var(--color-primary);
  background: rgba(var(--color-primary-rgb), 0.05);
}

.status-card.highlight {
  border-color: var(--ct-orange-500);
  background: rgba(255, 165, 0, 0.05);
}

.status-card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--ct-slate-600);
}

.status-card-title {
  font-size: 0.875rem;
  font-weight: 500;
}

.status-card-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--ct-slate-900);
}

.status-card.primary .status-card-value {
  color: var(--color-primary);
}

.status-card.highlight .status-card-value {
  color: var(--ct-orange-600);
}
</style>
