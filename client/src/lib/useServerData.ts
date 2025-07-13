import { api } from '@/api'
import router from '@/router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'

export const useUsersData = () =>
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

export const useCreateTeamInstance = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string }) =>
      api.POST('/teams/{teamId}/instances', {
        params: { path: params },
      }),
    onSuccess: (_, params) => {
      void client.invalidateQueries({ queryKey: ['team-instances', params.teamId] })
      void client.invalidateQueries({ queryKey: ['instances'] })
    },
  })
}

export const useStartTeamInstance = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string; instanceId: string }) =>
      api.PATCH('/teams/{teamId}/instances/{instanceId}', {
        params: { path: params },
        body: { operation: 'start' },
      }),
    onSuccess: (_, params) => {
      void client.invalidateQueries({ queryKey: ['team-instances', params.teamId] })
      void client.invalidateQueries({ queryKey: ['instances'] })
    },
  })
}

export const useStopTeamInstance = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string; instanceId: string }) =>
      api.PATCH('/teams/{teamId}/instances/{instanceId}', {
        params: { path: params },
        body: { operation: 'stop' },
      }),
    onSuccess: (_, params) => {
      void client.invalidateQueries({ queryKey: ['team-instances', params.teamId] })
      void client.invalidateQueries({ queryKey: ['instances'] })
    },
  })
}

export const useDeleteTeamInstance = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string; instanceId: string }) =>
      api.DELETE('/teams/{teamId}/instances/{instanceId}', {
        params: { path: params },
      }),
    onSuccess: (_, params) => {
      void client.invalidateQueries({ queryKey: ['team-instances', params.teamId] })
      void client.invalidateQueries({ queryKey: ['instances'] })
    },
  })
}

export const useCreateTeam = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { name: string; members: string[] }) =>
      api.POST('/teams', {
        body: params,
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: ['teams'] })
      void client.invalidateQueries({ queryKey: ['me'] })
    },
  })
}

export const useUpdateTeam = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string; name: string; members: string[] }) =>
      api.PATCH('/teams/{teamId}', {
        params: { path: { teamId: params.teamId } },
        body: params,
      }),
    onSuccess: (_, params) => {
      void client.invalidateQueries({ queryKey: ['team', params.teamId] })
      void client.invalidateQueries({ queryKey: ['teams'] })
      void client.invalidateQueries({ queryKey: ['me'] })
    },
  })
}

export const useEnqueueBenchmark = (options?: { redirect?: boolean }) => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { teamId: string; instanceId: string }) =>
      api.POST('/benchmarks', {
        body: { instanceId: params.instanceId },
      }),
    onSuccess: (res, params) => {
      void client.invalidateQueries({ queryKey: ['team-benches', params.teamId] })
      void client.invalidateQueries({ queryKey: ['benches'] })
      if (options?.redirect && res.data !== undefined) {
        void router.push(`/benches/${res.data.id}`)
      }
    },
  })
}

export const useAllInstances = () =>
  useQuery({
    queryKey: ['instances'],
    queryFn: () => api.GET('/instances').then((r) => r.data),
  })

export const useTeamsData = () =>
  useQuery({
    queryKey: ['teams'],
    queryFn: () => api.GET('/teams').then((r) => r.data),
  })

export const useAllBenches = () =>
  useQuery({
    queryKey: ['benches'],
    queryFn: () => api.GET('/benchmarks').then((r) => r.data),
  })

export const useBench = (benchmarkId: string) =>
  useQuery({
    queryKey: ['bench', benchmarkId],
    queryFn: () =>
      api
        .GET('/benchmarks/{benchmarkId}', { params: { path: { benchmarkId } } })
        .then((r) => r.data),
  })

export const useDocs = () =>
  useQuery({
    queryKey: ['docs'],
    queryFn: () => api.GET('/docs').then((r) => r.data),
  })

export const useUpdateDocs = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: { body: string }) =>
      api.PATCH('/docs', {
        body: params,
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: ['docs'] })
    },
  })
}

export const useUpdateAdmins = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (params: string[]) =>
      api.PUT('/admins', {
        body: params,
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: ['admins'] })
      void client.invalidateQueries({ queryKey: ['users'] })
    },
  })
}

export const useScores = () =>
  useQuery({
    queryKey: ['scores'],
    queryFn: () => api.GET('/scores').then((r) => r.data),
  })

export const useRanking = (orderBy: 'latest' | 'highest') =>
  useQuery({
    queryKey: ['ranking', orderBy],
    queryFn: () =>
      api.GET('/scores/ranking', { params: { query: { orderBy } } }).then((r) => r.data),
  })

export const useBenchmarkQueue = () =>
  useQuery({
    queryKey: ['benchmark-queue'],
    queryFn: () => api.GET('/benchmarks/queue').then((r) => r.data),
  })
