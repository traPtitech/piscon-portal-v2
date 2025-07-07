<script setup lang="ts">
import type { RouterLinkProps } from 'vue-router'
defineProps<RouterLinkProps>()
</script>

<template>
  <RouterLink v-bind="$props" :class="['navigation-link', $props.class]">
    <slot />
  </RouterLink>
</template>

<style scoped>
.navigation-link {
  text-decoration: none;
  color: var(--color-primary);
  font-weight: 700;
  transition: color 0.2s;
}

.navigation-link:hover {
  color: var(--color-primary-hover-text);
}
</style>
