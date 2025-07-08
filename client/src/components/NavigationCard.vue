<script setup lang="ts">
import { Icon } from '@iconify/vue'
import NavigationLink from '@/components/NavigationLink.vue'

defineProps<{
  icon: string
  title: string
  description: string
  to: string
}>()
</script>

<template>
  <NavigationLink :to="to" class="nav-card">
    <div class="nav-card-content">
      <div class="nav-card-header">
        <Icon :icon="icon" width="24" height="24" />
        <span class="nav-card-title">{{ title }}</span>
      </div>
      <p class="nav-card-description">{{ description }}</p>
    </div>
    <Icon icon="mdi:chevron-right" width="20" height="20" class="nav-card-arrow" />
  </NavigationLink>
</template>

<style scoped>
.nav-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 8px;
  background: var(--ct-white);
  transition: all 0.2s;
  text-decoration: none;
  color: inherit;
}

.nav-card:hover {
  border-color: var(--color-primary);
  background: rgba(var(--color-primary-rgb), 0.05);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.nav-card-content {
  flex: 1;
}

.nav-card:hover .nav-card-header span {
  color: var(--color-primary);
}

.nav-card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.nav-card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--ct-slate-900);
}

.nav-card-description {
  margin: 0;
  font-size: 0.875rem;
  color: var(--ct-slate-600);
  line-height: 1.4;
}

.nav-card-arrow {
  color: var(--ct-slate-400);
  transition: color 0.2s;
}

.nav-card:hover .nav-card-arrow {
  color: var(--color-primary);
}
</style>
