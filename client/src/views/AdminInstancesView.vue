<script setup lang="ts">
import InstanceCardList from '@/components/InstanceCardList.vue'
import MainSwitch from '@/components/MainSwitch.vue'
import PageTitle from '@/components/PageTitle.vue'
import { useAllInstances, useTeams } from '@/lib/useServerData'
import { useUsers } from '@/lib/useUsers'
import { ref, computed } from 'vue'

const { data: instances } = useAllInstances()
const { data: teams } = useTeams()
const { getUserById } = useUsers()

const showDeleted = ref(false)

const visibleInstances = computed(() =>
  instances.value?.filter((i) => showDeleted.value || i.status !== 'deleted'),
)
const instancesByTeams = computed(() =>
  teams.value?.map((team) => ({
    team,
    instances: visibleInstances.value?.filter((i) => i.teamId === team.id) || [],
  })),
)
</script>

<template>
  <main class="admin-instances-container">
    <PageTitle icon="mdi:database-cog">インスタンス (管理者)</PageTitle>

    <div class="admin-instances-header">
      <MainSwitch v-model="showDeleted">削除済みのインスタンスを表示</MainSwitch>
    </div>

    <div class="team-instances">
      <div v-for="team in instancesByTeams" :key="team.team.id" class="team-instance">
        <h2 class="team-instance-title" :id="team.team.id">
          <div>{{ team.team.name }}</div>
          <div class="team-members">
            <UserAvatar
              v-for="member in team.team.members"
              :key="member"
              :name="getUserById(member)?.name ?? ''"
              :title="getUserById(member)?.name ?? ''"
            />
          </div>
        </h2>
        <InstanceCardList :teamId="team.team.id" :instances="team.instances" />
      </div>
    </div>
  </main>
</template>

<style scoped>
.admin-instances-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.admin-instances-header {
  display: flex;
  justify-content: flex-end;
}

.team-instances {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.team-instance {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
}

.team-instance-title {
  font-size: 1.5rem;
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.team-members {
  display: flex;
  gap: 0.5rem;
}
</style>
