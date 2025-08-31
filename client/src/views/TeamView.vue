<script setup lang="ts">
import PageTitle from '@/components/PageTitle.vue'
import TeamManagement from '@/components/TeamManagement.vue'
import { useCreateTeam, useMe } from '@/lib/useServerData'
import notfoundImage from '@/assets/not-found.png'
import ActionFormCard from '@/components/ActionFormCard.vue'

const { data: me } = useMe()
const { mutate: createTeam } = useCreateTeam()

const createTeamHandler = (teamName: string) => {
  if (me.value === undefined) return false
  createTeam({ name: teamName, members: [me.value.id] })
  return true
}
</script>

<template>
  <main class="team-container">
    <PageTitle icon="mdi:account-group">チーム管理</PageTitle>
    <div v-if="me?.teamId !== undefined">
      <TeamManagement :teamId="me.teamId" />
    </div>
    <div v-if="me !== undefined && me?.teamId === undefined" class="no-team-container">
      <div class="no-team-display">
        <img :src="notfoundImage" alt="" width="192" height="192" />
        <div>現在チームに所属していません</div>
      </div>
      <ActionFormCard
        icon="mdi:account-multiple-plus"
        title="新しいチームを作成"
        inputPlaceholder="チーム名"
        :action="createTeamHandler"
        actionIcon="mdi:account-multiple-plus"
        actionLabel="作成"
      />
    </div>
  </main>
</template>

<style scoped>
.team-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.no-team-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.no-team-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-weight: 600;
}

.create-team-form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
}

.create-team-form-header {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
}

.create-team-form-body {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.create-team-input {
  flex: 1;
  font-size: 0.8rem;
}
</style>
