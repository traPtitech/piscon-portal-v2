<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import { useMe, useTeam, useUpdateTeam } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import MainButton from '@/components/MainButton.vue'
import ActionFormCard from '@/components/ActionFormCard.vue'
import { useUsers } from '@/lib/useUsers'
import ErrorMessage from '@/components/ErrorMessage.vue'

const { teamId } = defineProps<{ teamId: string }>()

const { data: me } = useMe()
const { data: team, error: teamError } = useTeam(teamId)
const { getUserById, getUserByName } = useUsers()
const { mutate: updateTeam } = useUpdateTeam()

const changeTeamName = (name: string) => {
  if (team.value === undefined) return false

  updateTeam({ teamId: team.value.id, name, members: team.value.members })
  return true
}

const removeMember = (memberId: string) => {
  if (team.value === undefined) return false
  if (!team.value.members.includes(memberId)) return false

  const newMembers = team.value.members.filter((member) => member !== memberId)
  updateTeam({ teamId: team.value.id, name: team.value.name, members: newMembers })
  return true
}

const addMember = (memberId: string) => {
  if (team.value === undefined) return false
  if (team.value.members.includes(memberId)) return false

  const newMembers = [...team.value.members, memberId]
  updateTeam({ teamId: team.value.id, name: team.value.name, members: newMembers })
  return true
}

const addNewMemberHandler = (newMemberName: string) => {
  const user = getUserByName(newMemberName)
  if (user === undefined) return false
  return addMember(user.id)
}
</script>

<template>
  <div>
    <ErrorMessage v-if="teamError" :error="teamError" />
    <div v-else-if="team !== undefined" class="team-management-container">
      <div class="team-info">
        <div class="team-name">チーム: {{ team.name }}</div>
        <MainButton
          v-if="me"
          @click="removeMember(me.id)"
          variant="destructive"
          class="leave-team-button"
        >
          <Icon icon="mdi:exit-run" width="20" height="20" />
          <span>チームを抜ける</span>
        </MainButton>
      </div>
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
        <ActionFormCard
          icon="mdi:account-plus"
          title="メンバー追加"
          inputPlaceholder="メンバー名 (例: cp20)"
          :action="addNewMemberHandler"
          actionIcon="mdi:account-plus"
          actionLabel="追加"
        />
        <ActionFormCard
          icon="mdi:rename"
          title="チーム名変更"
          inputPlaceholder="新しいチーム名"
          :action="changeTeamName"
          actionIcon="mdi:content-save"
          actionLabel="保存"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.team-management-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  container-type: inline-size;
}

.team-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.team-name {
  font-size: 1.5rem;
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
  gap: 0.5rem;
}

.add-member-form,
.team-name-form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
}

.add-member-form-header,
.team-name-form-header {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
}

.add-member-form-body,
.team-name-form-body {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.add-member-input,
.team-name-input {
  flex: 1;
  font-size: 0.8rem;
}

@container (max-width: 768px) {
  .team-management-forms {
    grid-template-columns: 1fr;
  }
}
</style>
