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

export const useTeamBench = (teamId: string, benchmarkId: string) =>
  useQuery({
    queryKey: ['team-bench', teamId, benchmarkId],
    queryFn: () =>
      api
        .GET('/teams/{teamId}/benchmarks/{benchmarkId}', {
          params: { path: { teamId, benchmarkId } },
        })
        .then((r) => r.data),
  })

export const useTeamInstances = (teamId: string) =>
  useQuery({
    queryKey: ['team-instances', teamId],
    queryFn: () =>
      api.GET('/teams/{teamId}/instances', { params: { path: { teamId } } }).then((r) => r.data),
  })

export const useTeam = (teamId: string) =>
  useQuery({
    queryKey: ['team', teamId],
    queryFn: () => api.GET('/teams/{teamId}', { params: { path: { teamId } } }).then((r) => r.data),
  })
