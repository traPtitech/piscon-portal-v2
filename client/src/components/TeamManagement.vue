<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import { useMe, useTeam, useUpdateTeam } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import MainButton from '@/components/MainButton.vue'
import ActionFormCard from '@/components/ActionFormCard.vue'
import { useUsers } from '@/lib/useUsers'
import ErrorMessage from '@/components/ErrorMessage.vue'
import { ref } from 'vue'

const { teamId } = defineProps<{ teamId: string }>()

const { data: me } = useMe()
const { data: team, error: teamError } = useTeam(teamId)
const { getUserById, getUserByName } = useUsers()
const { mutate: updateTeam } = useUpdateTeam()

const changeTeamName = (name: string) => {
  if (team.value === undefined) return false

  updateTeam({ teamId: team.value.id, name })
  return true
}

const removeMember = (memberId: string) => {
  if (team.value === undefined) return false
  if (!team.value.members.includes(memberId)) return false

  const newMembers = team.value.members.filter((member) => member !== memberId)
  updateTeam({ teamId: team.value.id, members: newMembers })
  return true
}

const addMember = (memberId: string) => {
  if (team.value === undefined) return false
  if (team.value.members.includes(memberId)) return false

  const newMembers = [...team.value.members, memberId]
  updateTeam({ teamId: team.value.id, members: newMembers })
  return true
}

const addNewMemberHandler = (newMemberName: string) => {
  const user = getUserByName(newMemberName)
  if (user === undefined) return false
  return addMember(user.id)
}

const githubIdValue = ref('')

const addGitHubId = () => {
  if (githubIdValue.value === '') return
  if (team.value?.githubIds?.includes(githubIdValue.value)) return

  updateTeam({
    teamId: teamId,
    githubIds: [...(team.value?.githubIds ?? []), githubIdValue.value],
  })

  githubIdValue.value = ''
}

const removeGitHubId = (id: string) => {
  if (team.value === undefined) return
  if (!team.value.githubIds?.includes(id)) return

  const newGithubIds = team.value.githubIds.filter((githubId) => githubId !== id)
  updateTeam({ teamId: team.value.id, githubIds: newGithubIds })
  githubIdValue.value = ''
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
          <div class="member-info">
            <UserAvatar :name="getUserById(member)?.name ?? ''" />
            <div>{{ getUserById(member)?.name }}</div>
          </div>
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
      <div class="github-id-management-card">
        <div>
          <div class="github-id-management-card-title">GitHub ID 管理</div>
          <div class="github-id-management-card-description">
            チームメンバーの GitHub ID を追加すると、GitHub
            に登録されている公開鍵が自動的にインスタンスに登録されます
          </div>
          <div class="github-id-management-card-description">
            ※設定が反映されるのは新しく作成されたインスタンスのみです
          </div>
        </div>
        <div>
          <div class="github-id-chip-container">
            <div v-for="id in team.githubIds" :key="id" class="github-id-chip">
              <Icon icon="mdi:github" width="20" height="20" />
              <a :href="`https://github.com/${id}`">
                {{ id }}
              </a>
              <button @click.prevent="removeGitHubId(id)">
                <Icon icon="mdi:close" width="20" height="20" />
              </button>
            </div>
          </div>
          <div
            v-if="team.githubIds === undefined || team.githubIds.length === 0"
            class="no-github-id-registered"
          >
            まだ GitHub ID が登録されていません
          </div>
          <form class="github-id-management-form" @submit.prevent="addGitHubId">
            <InputText
              v-model="githubIdValue"
              placeholder="GitHub ID (例: cp-20)"
              class="github-id-management-form-input"
            />
            <MainButton type="submit">
              <Icon icon="mdi:account-plus" width="20" height="20" />
              <span>追加</span>
            </MainButton>
          </form>
        </div>
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

@container (max-width: 480px) {
  .team-info {
    flex-direction: column;
    align-items: stretch;
    padding: 1rem;
  }

  .leave-team-button {
    width: 100%;
  }
}

.members-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(320px, 100%), 1fr));
  gap: 0.5rem;
}

.member-container {
  display: flex;
  align-items: center;
  font-weight: 600;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  gap: 0.25rem;
  border: 1px solid var(--ct-slate-300);
}

.member-info {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem;
}

.remove-member-button {
  margin-left: auto;
}

@container (max-width: 400px) {
  .member-container {
    flex-direction: column;
    align-items: stretch;
    padding: 1rem;
  }

  .remove-member-button {
    width: 100%;
  }
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

.github-id-management-card {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
  container-type: inline-size;
}

.github-id-chip-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.github-id-management-card-title {
  font-weight: 600;
  font-size: 1.2rem;
  margin-bottom: 0.25rem;
}

.github-id-management-card-description {
  font-size: 0.9rem;
  color: var(--ct-slate-500);
}

@container (max-width: 400px) {
  .github-id-management-card-description {
    font-size: 0.8rem;
  }
}

.github-id-chip {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  background-color: var(--ct-slate-100);
  color: inherit;
  text-decoration: none;
}

.github-id-chip a {
  color: inherit;
}

.github-id-chip button {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  width: 20px;
  height: 20px;
  margin-left: 0.25rem;
}

.no-github-id-registered {
  color: var(--ct-slate-500);
  font-weight: bold;
  text-align: center;
  margin-bottom: 1rem;
}

.github-id-management-form {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.github-id-management-form-input {
  flex: 1;
  font-size: 0.8rem;
}

@container (max-width: 320px) {
  .github-id-management-form {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
