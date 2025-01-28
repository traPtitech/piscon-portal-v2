<script setup lang="ts">
import MainButton from '@/components/MainButton.vue'
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const { action } = defineProps<{
  icon: string
  title: string
  inputPlaceholder: string
  action: (value: string) => void
  actionIcon: string
  actionLabel: string
}>()

const value = ref('')

const actionHandler = () => {
  if (value.value === '') return
  action(value.value)
  value.value = ''
}
</script>

<template>
  <div class="action-form-card">
    <div class="action-form-card-header">
      <Icon :icon="icon" width="20" height="20" />
      <span>{{ title }}</span>
    </div>
    <form class="action-form-card-body" @submit.prevent="actionHandler">
      <InputText v-model="value" :placeholder="inputPlaceholder" class="action-form-input" />
      <MainButton type="submit">
        <Icon :icon="actionIcon" width="20" height="20" />
        <span>{{ actionLabel }}</span>
      </MainButton>
    </form>
  </div>
</template>

<style scoped>
.action-form-card {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
  container-type: inline-size;
}

.action-form-card-header {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
}

.action-form-card-body {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.action-form-input {
  flex: 1;
  font-size: 0.8rem;
}

@container (max-width: 320px) {
  .action-form-card-body {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
