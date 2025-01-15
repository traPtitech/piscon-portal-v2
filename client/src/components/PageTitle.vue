<script setup lang="ts">
import { Icon } from '@iconify/vue'

const { icon } = defineProps<{
  icon: string
}>()
</script>

<template>
  <h1 class="title">
    <Icon :icon="icon" width="32" height="32" />
    <span>
      <slot />
    </span>
  </h1>
</template>

<style scoped>
.title {
  display: flex;
  align-items: center;
  font-size: 1.5rem;
}
</style>
