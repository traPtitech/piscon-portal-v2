import { apiBaseUrl } from '@/api'
import { http, HttpResponse } from 'msw'
import type { components, paths } from '@/api/openapi'

const userIds = {
  cp20: '01943f2e-0dca-7599-9f79-901de49660fd',
  'ikura-hamu': '01943f39-4304-7cf0-b420-2bb5cfde5f1a',
  pirosiki: '01943f39-8171-7db8-96d8-e672a72d14fd',
}

const teamIds = {
  'ikura-cp': '01943f2e-4e49-782c-9eb8-6a496b509bc8',
  piropiro: '01943f39-9502-76b2-8ace-a1cb3931fcf3',
}

const users: paths['/users']['get']['responses']['200']['content']['application/json'] = [
  {
    id: userIds.cp20,
    isAdmin: true,
    name: 'cp20',
    teamId: teamIds['ikura-cp'],
  },
  {
    id: userIds['ikura-hamu'],
    isAdmin: false,
    name: 'ikura-hamu',
    teamId: teamIds['ikura-cp'],
  },
  {
    id: userIds.pirosiki,
    isAdmin: false,
    name: 'pirosiki',
    teamId: teamIds.piropiro,
  },
]

const teams: paths['/teams']['get']['responses']['200']['content']['application/json'] = [
  {
    id: teamIds['ikura-cp'],
    name: 'ikura-cp',
    members: [userIds.cp20, userIds['ikura-hamu']],
    createdAt: '2025-01-01T00:00:00Z',
    githubIds: ['cp-20', 'ikura-hamu'],
  },
  {
    id: teamIds.piropiro,
    name: 'piropiro',
    members: [userIds.pirosiki],
    createdAt: '2025-01-02T00:00:00Z',
    githubIds: ['pirosiki197'],
  },
]

const instanceIds = {
  'ikura-cp-1': '01943f39-b74e-7d2d-8a16-9d408bf92579',
  'ikura-cp-2': '01943f4b-bd05-72bd-b603-633e49b94edb',
  'ikura-cp-3': '01943f4c-21a3-7c95-ba5a-63c7ea02b9f0',
  'piropiro-1': '01943f4c-3662-7cac-9674-2b0a39da6725',
}

const instances: paths['/teams/{teamId}/instances']['get']['responses']['200']['content']['application/json'] =
  [
    {
      id: instanceIds['ikura-cp-1'],
      teamId: teamIds['ikura-cp'],
      serverId: 1,
      privateIPAddress: '192.168.0.1',
      publicIPAddress: '203.0.113.1',
      status: 'running',
      createdAt: '2025-01-01T00:00:00Z',
    },
    {
      id: instanceIds['ikura-cp-2'],
      teamId: teamIds['ikura-cp'],
      serverId: 2,
      privateIPAddress: '192.168.0.2',
      publicIPAddress: '203.0.113.2',
      status: 'building',
      createdAt: '2025-01-01T00:00:00Z',
    },
    {
      id: instanceIds['ikura-cp-3'],
      teamId: teamIds['ikura-cp'],
      serverId: 3,
      privateIPAddress: '192.168.0.3',
      publicIPAddress: '203.0.113.3',
      status: 'stopped',
      createdAt: '2025-01-01T00:00:00Z',
    },
    {
      id: instanceIds['piropiro-1'],
      teamId: teamIds['piropiro'],
      serverId: 1,
      privateIPAddress: '192.168.0.4',
      publicIPAddress: '203.0.113.4',
      status: 'running',
      createdAt: '2025-01-02T00:00:00Z',
    },
  ]

const benchmarks: paths['/benchmarks']['get']['responses']['200']['content']['application/json'] = [
  {
    id: '01943f67-d9ed-7bbb-81eb-20f81391ffea',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'finished',
    createdAt: '2025-01-01T01:00:00Z',
    startedAt: '2025-01-01T01:00:01Z',
    finishedAt: '2025-01-01T01:01:00Z',
    score: 2000,
  },
  {
    id: '01943f68-7d22-7abb-8b13-0b727cd4597e',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'running',
    createdAt: '2025-01-01T02:00:00Z',
    startedAt: '2025-01-01T02:00:01Z',
  },
  {
    id: '01943f69-3aec-7702-8d1e-8642d9c5b47b',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'waiting',
    createdAt: '2025-01-01T03:00:00Z',
  },
  {
    id: '01943f6b-7276-79a9-9bd3-69d1d5a9cb3d',
    instanceId: instanceIds['piropiro-1'],
    teamId: teamIds['piropiro'],
    userId: userIds.pirosiki,
    status: 'finished',
    createdAt: '2025-01-02T01:00:00Z',
    startedAt: '2025-01-02T01:00:01Z',
    finishedAt: '2025-01-02T01:01:00Z',
    score: 100,
  },
  {
    id: '01943f6e-69dd-7167-84b3-478cf9c3253d',
    instanceId: instanceIds['piropiro-1'],
    teamId: teamIds['piropiro'],
    userId: userIds.pirosiki,
    status: 'finished',
    createdAt: '2025-01-02T02:00:00Z',
    startedAt: '2025-01-02T02:00:01Z',
    finishedAt: '2025-01-02T02:01:00Z',
    score: 1000,
  },
  {
    id: '01943f6e-8b29-79af-8430-7b06ae9307e5',
    instanceId: instanceIds['piropiro-1'],
    teamId: teamIds['piropiro'],
    userId: userIds.pirosiki,
    status: 'finished',
    createdAt: '2025-01-02T03:00:00Z',
    startedAt: '2025-01-02T03:00:01Z',
    finishedAt: '2025-01-02T03:01:00Z',
    score: 1000,
  },
]

export const handlers = [
  http.get(`${apiBaseUrl}/oauth2/code`, () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/oauth2/callback`, () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/oauth2/logout`, () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/users`, () => HttpResponse.json(users)),
  http.get(`${apiBaseUrl}/users/me`, () => HttpResponse.json(users[0])),
  http.get(`${apiBaseUrl}/teams`, () => HttpResponse.json(teams)),
  http.post(`${apiBaseUrl}/teams`, () => {
    // TODO
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)$`), (c) => {
    // TODO
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(team)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances`), (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const res = instances.filter((i) => i.teamId === teamId)
    return HttpResponse.json(res)
  }),
  http.post(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances`), () => {
    // TODO
  }),
  http.delete(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances/([^/]+)`), () => {
    // TODO
  }),
  http.patch(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances/([^/]+)`), () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/instances`, () => HttpResponse.json(instances)),
  http.post(`${apiBaseUrl}/benchmarks`, () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/benchmarks`, () => HttpResponse.json(benchmarks)),
  http.get(`${apiBaseUrl}/benchmarks/queue`, () => {
    const res = benchmarks.filter((b) => b.status === 'waiting' || b.status === 'running')
    return HttpResponse.json(res)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)/benchmarks`), (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const res = benchmarks.filter((b) => b.teamId === teamId)
    return HttpResponse.json(res)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)/benchmarks/([^/]+)`), (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const benchmarkId = c.params[1] as string
    const benchmark = benchmarks.find((b) => b.id === benchmarkId)
    if (benchmark === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(benchmark)
  }),
  http.patch(new RegExp(`${apiBaseUrl}/benchmarks/([^/]+)`), () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/scores`, () => {
    const res: paths['/scores']['get']['responses']['200']['content']['application/json'] =
      teams.map((team) => ({
        teamId: team.id,
        scores: benchmarks
          .filter((b) => b.teamId === team.id && b.status === 'finished')
          .map((b) => ({
            benchmarkId: b.id,
            teamId: team.id,
            score: b.score!,
            createdAt: b.createdAt,
          })),
      }))

    return HttpResponse.json(res)
  }),
  http.get(`${apiBaseUrl}/scores/ranking`, () => {
    const sortFn = (
      a: components['schemas']['Benchmark'],
      b: components['schemas']['Benchmark'],
    ) => {
      // スコアの降順
      if (a.score! < b.score!) return 1
      if (a.score! > b.score!) return -1
      // スコアが同じなら、作成日時の昇順
      if (new Date(a.createdAt) > new Date(b.createdAt)) return 1
      return -1
    }
    const res: paths['/scores/ranking']['get']['responses']['200']['content']['application/json'] =
      teams
        .map(
          (team) =>
            benchmarks
              .filter((b) => b.teamId === team.id && b.status === 'finished')
              .sort(sortFn)[0],
        )
        .sort(sortFn)
        .map((b, i) => ({
          rank: i + 1,
          teamId: b.teamId,
          score: b.score!,
          createdAt: b.createdAt,
        }))

    return HttpResponse.json(res)
  }),
  http.put(`${apiBaseUrl}/admins`, () => {
    // TODO
  }),
  http.get(`${apiBaseUrl}/docs`, () => {
    const res: paths['/docs']['get']['responses']['200']['content']['application/json'] = {
      body: 'This is a document.',
    }

    return HttpResponse.json(res)
  }),
  http.patch(`${apiBaseUrl}/docs`, () => {
    // TODO
  }),
]
