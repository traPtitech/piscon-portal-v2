import { api } from '@/api'
import { useQuery } from '@tanstack/vue-query'

export const useUsers = () =>
  useQuery({
    queryKey: ['users'],
    queryFn: () => api.GET('/users').then((r) => r.data),
  })

export const useMe = () =>
  useQuery({
    queryKey: ['me'],
    queryFn: () => api.GET('/users/me').then((r) => r.data),
  })

export const useTeamBenches = (teamId: string) =>
  useQuery({
    queryKey: ['team-benches', teamId],
    queryFn: () =>
      api.GET('/teams/{teamId}/benchmarks', { params: { path: { teamId } } }).then((r) => r.data),
  })

export const useMyTeamBenches = () =>
  useQuery({
    queryKey: ['my-benches'],
    queryFn: async () => {
      const me = await api.GET('/users/me').then((r) => r.data)
      if (me?.teamId === undefined) return []
      const benches = await api
        .GET('/teams/{teamId}/benchmarks', { params: { path: { teamId: me?.teamId } } })
        .then((r) => r.data)
      return benches
    },
  })

export const useMyTeamInstances = () =>
  useQuery({
    queryKey: ['my-instances'],
    queryFn: async () => {
      const me = await api.GET('/users/me').then((r) => r.data)
      if (me?.teamId === undefined) return []
      const instances = await api
        .GET('/teams/{teamId}/instances', { params: { path: { teamId: me?.teamId } } })
        .then((r) => r.data)
      return instances
    },
  })

export const useTeamInstances = (teamId: string) =>
  useQuery({
    queryKey: ['team-instances', teamId],
    queryFn: () =>
      api.GET('/teams/{teamId}/instances', { params: { path: { teamId } } }).then((r) => r.data),
  })

export const useMyTeam = () =>
  useQuery({
    queryKey: ['my-team'],
    queryFn: async () => {
      const me = await api.GET('/users/me').then((r) => r.data)
      if (me?.teamId === undefined) return undefined
      const team = await api
        .GET('/teams/{teamId}', { params: { path: { teamId: me?.teamId } } })
        .then((r) => r.data)
      return team
    },
  })

export const useTeam = (teamId: string) =>
  useQuery({
    queryKey: ['team', teamId],
    queryFn: () => api.GET('/teams/{teamId}', { params: { path: { teamId } } }).then((r) => r.data),
  })
