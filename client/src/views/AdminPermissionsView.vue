<script setup lang="ts">
import ActionFormCard from '@/components/ActionFormCard.vue'
import PageTitle from '@/components/PageTitle.vue'
import { useUsers } from '@/lib/useUsers'
import { computed } from 'vue'
import { useUpdateAdmins } from '@/lib/useServerData'
import UserAvatar from '@/components/UserAvatar.vue'
import MainButton from '@/components/MainButton.vue'
import { Icon } from '@iconify/vue'

const { users, getUserByName } = useUsers()
const { mutate: updateAdmins } = useUpdateAdmins()

const admins = computed(() => users.value?.filter((u) => u.isAdmin))

const addAdminHandler = (value: string) => {
  if (value === '') return
  const newUser = getUserByName(value)
  if (!newUser) return
  const newAdmins = [...(admins.value?.map((u) => u.id) ?? []), newUser.id]
  updateAdmins(newAdmins)
}

const removeAdminHandler = (userId: string) => {
  if (admins.value === undefined) return
  const newAdmins = admins.value.filter((u) => u.id !== userId).map((u) => u.id)
  updateAdmins(newAdmins)
}
</script>

<template>
  <main class="admin-permissions-container">
    <PageTitle icon="mdi:account-lock">権限管理 (管理者)</PageTitle>

    <div class="admin-containers">
      <div v-for="admin in admins" :key="admin.id" class="admin-container">
        <UserAvatar :name="admin.name" />
        <div>{{ admin.name }}</div>
        <MainButton
          @click="removeAdminHandler(admin.id)"
          variant="destructive"
          class="admin-delete-button"
        >
          <Icon icon="mdi:account-remove" width="20" height="20" />
          <span>削除</span>
        </MainButton>
      </div>
    </div>

    <ActionFormCard
      icon="mdi:account-plus"
      title="新しい管理者を追加"
      inputPlaceholder="ユーザー名 (例: cp20)"
      :action="addAdminHandler"
      actionIcon="mdi:account-plus"
      actionLabel="追加"
    />
  </main>
</template>

<style scoped>
.admin-permissions-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.admin-containers {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(320px, 100%), 1fr));
  gap: 0.5rem;
}

.admin-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 0.5rem 1rem;
  font-weight: 600;
}

.admin-delete-button {
  margin-left: auto;
}
</style>
