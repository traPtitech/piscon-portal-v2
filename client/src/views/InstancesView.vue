<script setup lang="ts">
import PageTitle from '@/components/PageTitle.vue'
import TeamInstanceManagement from '@/components/TeamInstanceManagement.vue'
import { useMe } from '@/lib/useServerData'

const { data: me } = useMe()
</script>

<template>
  <main class="instances-container">
    <PageTitle icon="mdi:server-network">インスタンス</PageTitle>
    <div v-if="me?.teamId !== undefined">
      <TeamInstanceManagement :teamId="me.teamId" />
    </div>
    <div v-if="me !== undefined && me?.teamId === undefined">
      <p>チームに所属していません</p>
    </div>
  </main>
</template>

<style scoped>
.instances-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
