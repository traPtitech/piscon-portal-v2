import { useTeamsData } from '@/lib/useServerData'
import { useUsers } from '@/lib/useUsers'

export const useTeams = () => {
  const { data: teams, ...rest } = useTeamsData()
  const { getUserById } = useUsers()

  const getTeamName = (teamId: string) => {
    return teams.value?.find((team) => team.id === teamId)?.name ?? 'Unknown Team'
  }

  const getTeamMembers = (teamId: string) => {
    const memberIds = teams.value?.find((team) => team.id === teamId)?.members ?? []
    return memberIds.map((id) => getUserById(id)).filter((user) => user !== undefined)
  }

  return { getTeamName, getTeamMembers, teams, ...rest }
}
