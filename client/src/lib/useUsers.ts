import { useUsersData } from '@/lib/useServerData'
import { computed } from 'vue'

export const useUsers = () => {
  const { data: users, ...rest } = useUsersData()

  const userMapById = computed(() => new Map(users.value?.map((user) => [user.id, user])))
  const getUserById = (id: string | undefined) => userMapById.value.get(id ?? '')

  const userMapByName = computed(() => new Map(users.value?.map((user) => [user.name, user])))
  const getUserByName = (name: string | undefined) => userMapByName.value.get(name ?? '')

  return { getUserById, getUserByName, users, ...rest }
}
