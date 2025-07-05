<script setup lang="ts">
import MainLayout from '@/layouts/MainLayout.vue'
import { RouterView } from 'vue-router'
</script>

<template>
  <MainLayout>
    <RouterView :key="$route.fullPath" />
  </MainLayout>
</template>

<style scoped></style>
