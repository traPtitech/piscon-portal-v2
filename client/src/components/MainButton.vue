<script setup lang="ts">
import { defineProps } from 'vue'

defineProps<{
  variant?: 'primary' | 'destructive'
}>()
</script>

<template>
  <button class="button" :class="variant ?? 'primary'" v-bind="$attrs">
    <slot></slot>
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
</style>
