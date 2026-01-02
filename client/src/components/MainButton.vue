<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'destructive'
  loading?: boolean
  disabled?: boolean
}>()
</script>

<template>
  <button
    class="button"
    :class="variant ?? 'primary'"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <slot></slot>
    <Icon v-if="loading" icon="mdi:loading" width="20" height="20" class="loading-spinner" />
  </button>
</template>

<style scoped>
.button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 700;
  font-size: 0.9rem;
  transition: background-color 0.2s;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.25rem;
  position: relative;
}
.button:disabled:disabled {
  background-color: var(--ct-slate-100);
  color: var(--ct-slate-400);
  cursor: not-allowed;
}

.primary {
  background-color: rgba(var(--color-primary-rgb), 0.2);
  color: var(--color-primary);
}
.primary:hover {
  background-color: rgba(var(--color-primary-rgb), 0.3);
}

.destructive {
  background-color: var(--ct-red-100);
  color: var(--ct-red-500);
}
.destructive:hover {
  background-color: var(--ct-red-200);
}

.loading-spinner {
  display: inline-block;
  animation: spin 1s linear infinite;
  margin-left: 0.5rem;
  position: absolute;
  right: 0.5rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
