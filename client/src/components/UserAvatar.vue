<script setup lang="ts">
defineProps<{ name: string }>()
</script>

<template>
  <img :src="`https://q.trap.jp/api/v3/public/icon/${name}`" alt="" width="24" height="24" />
</template>
