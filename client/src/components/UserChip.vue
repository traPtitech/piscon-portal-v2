<script setup lang="ts">
defineProps<{ name: string }>()
</script>

<template>
  <div class="user-chip">
    <UserAvatar :name="name" class="user-avatar" />
    <span>{{ name }}</span>
  </div>
</template>

<style scoped>
.user-chip {
  display: flex;
  align-items: center;
  gap: 0.25rem;

  border-radius: 4px;
  padding: 0.15rem 0.25rem;
  background-color: var(--ct-slate-100);
}

.user-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

.user-chip span {
  font-size: 0.8rem;
  font-weight: 600;
}
</style>
