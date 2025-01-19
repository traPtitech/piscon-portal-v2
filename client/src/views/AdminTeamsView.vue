<script setup lang="ts">
import PageTitle from '@/components/PageTitle.vue'
import AdminTeamManagement from '@/components/AdminTeamManagement.vue'
import { useTeams } from '@/lib/useServerData'
import notfoundImage from '@/assets/not-found.png'

const { data: teams } = useTeams()
</script>

<template>
  <main class="admin-teams-container">
    <PageTitle icon="mdi:account-cog">チーム管理 (管理者)</PageTitle>

    <div v-if="teams !== undefined && teams.length === 0">
      <img :src="notfoundImage" alt="" />
      <div>現在チームがありません</div>
    </div>
    <div v-for="team in teams ?? []" :key="team.id">
      <AdminTeamManagement :teamId="team.id" />
    </div>
  </main>
</template>

<style scoped>
.admin-teams-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
