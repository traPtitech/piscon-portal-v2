<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import { useMe, useTeam, useUpdateTeam } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import MainButton from '@/components/MainButton.vue'
import { useUsers } from '@/lib/useUsers'
import { ref } from 'vue'

const { teamId } = defineProps<{ teamId: string }>()

const { data: me } = useMe()
const { data: team } = useTeam(teamId)
const { getUserById, getUserByName } = useUsers()
const { mutate: updateTeam } = useUpdateTeam()

const changeTeamName = (name: string) => {
  if (team.value === undefined) return
  updateTeam({ teamId: team.value.id, name, members: team.value.members })
}

const removeMember = (memberId: string) => {
  if (team.value === undefined) return
  if (!team.value.members.includes(memberId)) return

  const newMembers = team.value.members.filter((member) => member !== memberId)
  updateTeam({ teamId: team.value.id, name: team.value.name, members: newMembers })
}

const addMember = (memberId: string) => {
  if (team.value === undefined) return
  if (team.value.members.includes(memberId)) return

  const newMembers = [...team.value.members, memberId]
  updateTeam({ teamId: team.value.id, name: team.value.name, members: newMembers })
}

const newMemberName = ref('')
const addNewMemberHandler = () => {
  if (newMemberName.value === '') return
  const user = getUserByName(newMemberName.value)
  if (user === undefined) return
  addMember(user.id)
  newMemberName.value = ''
}

const newTeamName = ref('')
const changeTeamNameHandler = () => {
  if (newTeamName.value === '') return
  changeTeamName(newTeamName.value)
  newTeamName.value = ''
}
</script>

<template>
  <div class="team-management-container">
    <template v-if="team !== undefined">
      <div class="team-name">チーム: {{ team.name }}</div>
      <div class="members-list">
        <div v-for="member in team.members" :key="member" class="member-container">
          <UserAvatar :name="getUserById(member)?.name ?? ''" />
          <div>{{ getUserById(member)?.name }}</div>
          <MainButton
            @click="removeMember(member)"
            :disabled="me?.id === member"
            class="remove-member-button"
            variant="destructive"
          >
            <Icon icon="mdi:account-remove" width="20" height="20" />
            <span>削除</span>
          </MainButton>
        </div>
      </div>
      <div class="team-management-forms">
        <div class="team-management-form">
          <div class="team-management-form-header">
            <Icon icon="mdi:account-plus" width="20" height="20" />
            <span>メンバー追加</span>
          </div>
          <form class="team-management-form-body" @submit.prevent="addNewMemberHandler">
            <InputText
              v-model="newMemberName"
              placeholder="メンバー名 (例: cp20)"
              class="team-management-form-input"
            />
            <MainButton type="submit">
              <Icon icon="mdi:account-plus" width="20" height="20" />
              <span>追加</span>
            </MainButton>
          </form>
        </div>
        <div class="team-management-form">
          <div class="team-management-form-header">
            <Icon icon="mdi:rename" width="20" height="20" />
            <span>チーム名変更</span>
          </div>
          <form class="team-management-form-body" @submit.prevent="changeTeamNameHandler">
            <InputText
              v-model="newMemberName"
              placeholder="新しいチーム名"
              class="team-management-form-input"
            />
            <MainButton type="submit">
              <Icon icon="mdi:content-save" width="20" height="20" />
              <span>保存</span>
            </MainButton>
          </form>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.team-management-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  container-type: inline-size;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
}

.team-name {
  font-size: 1.2rem;
  font-weight: bold;
}

.leave-team-button {
  margin-left: auto;
}

.members-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(320px, 100%), 1fr));
  gap: 0.5rem;
}

.member-container {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: 1px solid var(--ct-slate-300);
}

.remove-member-button {
  margin-left: auto;
}

.team-management-forms {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.team-management-form {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.team-management-form-header {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
}

.team-management-form-body {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.team-management-form-input {
  flex: 1;
  font-size: 0.8rem;
}

@container (max-width: 768px) {
  .team-management-forms {
    grid-template-columns: 1fr;
  }
}
</style>
