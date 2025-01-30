import { apiBaseUrl } from '@/api'
import { http, HttpResponse } from 'msw'
import type { components, paths } from '@/api/openapi'
import { uuidv7 } from 'uuidv7'

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

const benchmarks: components['schemas']['BenchmarkAdminResult'][] = [
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
    log: '',
    adminLog: '',
    result: 'passed',
  },
  {
    id: '01943f68-7d22-7abb-8b13-0b727cd4597e',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'running',
    createdAt: '2025-01-01T02:00:00Z',
    startedAt: '2025-01-01T02:00:01Z',
    score: 0,
    log: '',
    adminLog: '',
  },
  {
    id: '01943f69-3aec-7702-8d1e-8642d9c5b47b',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'waiting',
    createdAt: '2025-01-01T03:00:00Z',
    log: '',
    adminLog: '',
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
    log: '',
    adminLog: '',
    result: 'passed',
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
    log: '',
    adminLog: '',
    result: 'passed',
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
    log: '',
    adminLog: '',
    result: 'passed',
  },
]

const docs = `# ドキュメント

## はじめに

これはドキュメントです。
同じ段落

違う段落

ここに説明が入る

## 使い方

1. これを使ってください。
2. あれを使ってください。
3. これを使ってください。

## 注意事項

- \`inline code\`
- **bold**
- *italic*
- [link](https://example.com)
- ==emphasized==
- ~deleted~

> 引用

## その他

\`\`\`js
const code = 'Hello, world!';
\`\`\`
`

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
  http.post(`${apiBaseUrl}/teams`, async (c) => {
    type Body = NonNullable<paths['/teams']['post']['requestBody']>['content']['application/json']
    const body = (await c.request.json()) as Body
    const newTeam: components['schemas']['Team'] = {
      id: uuidv7(),
      name: body.name,
      members: body.members,
      createdAt: new Date().toISOString(),
    }
    teams.push(newTeam)
    for (const member of body.members) {
      const user = users.find((u) => u.id === member)
      if (user !== undefined) user.teamId = newTeam.id
    }
    return HttpResponse.json(newTeam)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)$`), (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(team)
  }),
  http.patch(new RegExp(`${apiBaseUrl}/teams/([^/]+)$`), async (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })

    const body = (await c.request.json()) as { name: string; members: string[] }
    team.name = body.name
    const oldMembers = team.members
    team.members = body.members

    const removedMembers = oldMembers.filter((m) => !body.members.includes(m))
    for (const member of removedMembers) {
      const user = users.find((u) => u.id === member)
      if (user !== undefined) user.teamId = undefined
    }
    const addedMembers = body.members.filter((m) => !oldMembers.includes(m))
    for (const member of addedMembers) {
      const user = users.find((u) => u.id === member)
      if (user !== undefined) user.teamId = teamId
    }

    return HttpResponse.json(team)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances`), (c) => {
    const teamId = c.params[0] as string
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const res = instances.filter((i) => i.teamId === teamId)
    return HttpResponse.json(res)
  }),
  http.post(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances`), (c) => {
    const teamId = c.params[0] as string
    const newServerId = instances.reduce((max, i) => Math.max(max, i.serverId), 0) + 1
    const newInstance: components['schemas']['Instance'] = {
      id: uuidv7(),
      teamId,
      serverId: newServerId,
      privateIPAddress: `192.168.0.${newServerId}`,
      publicIPAddress: `203.0.113.${newServerId}`,
      status: 'building',
      createdAt: new Date().toISOString(),
    }
    instances.push(newInstance)
  }),
  http.delete(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances/([^/]+)`), (c) => {
    const teamId = c.params[0] as string
    const instanceId = c.params[1] as string
    const index = instances.findIndex((i) => i.teamId === teamId && i.id === instanceId)
    instances[index] = { ...instances[index], status: 'deleting' }
  }),
  http.patch(new RegExp(`${apiBaseUrl}/teams/([^/]+)/instances/([^/]+)`), async (c) => {
    const teamId = c.params[0] as string
    const instanceId = c.params[1] as string
    const index = instances.findIndex((i) => i.teamId === teamId && i.id === instanceId)
    type Body =
      paths['/teams/{teamId}/instances/{instanceId}']['patch']['requestBody']['content']['application/json']
    const body = (await c.request.json()) as Body
    const operation = body.operation as 'start' | 'stop'
    if (operation === 'start') {
      if (instances[index].status === 'stopped') {
        instances[index] = { ...instances[index], status: 'starting' }
        return HttpResponse.json({}, { status: 200 })
      }
    }
    if (operation === 'stop') {
      if (instances[index].status === 'running') {
        instances[index] = { ...instances[index], status: 'stopping' }
        return HttpResponse.json({}, { status: 200 })
      }
    }
    return HttpResponse.json({ message: 'Bad request' }, { status: 400 })
  }),
  http.get(`${apiBaseUrl}/instances`, () => HttpResponse.json(instances)),
  http.post(`${apiBaseUrl}/benchmarks`, async (c) => {
    type Body = NonNullable<
      paths['/benchmarks']['post']['requestBody']
    >['content']['application/json']
    const body = (await c.request.json()) as Body

    const instance = instances.find((i) => i.id === body.instanceId)
    if (instance === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })

    const me = users[0]

    const benchmark: components['schemas']['BenchmarkAdminResult'] = {
      id: uuidv7(),
      instanceId: instance.id,
      teamId: instance.teamId,
      userId: me.id,
      status: 'waiting',
      createdAt: new Date().toISOString(),
      log: '',
      adminLog: '',
    }

    benchmarks.push(benchmark)
    return HttpResponse.json(benchmark, { status: 201 })
  }),
  http.get(`${apiBaseUrl}/benchmarks`, () => HttpResponse.json(benchmarks)),
  http.get(`${apiBaseUrl}/benchmarks/queue`, () => {
    const res = benchmarks.filter((b) => b.status === 'waiting' || b.status === 'running')
    return HttpResponse.json(res)
  }),
  http.get(new RegExp(`${apiBaseUrl}/teams/([^/]+)/benchmarks$`), (c) => {
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
    const benchmark = benchmarks.find((b) => b.teamId === teamId && b.id === benchmarkId)
    if (benchmark === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const { adminLog: _, ...res } = benchmark
    return HttpResponse.json(res)
  }),
  http.get(new RegExp(`${apiBaseUrl}/benchmarks/([^/]+)`), (c) => {
    const benchmarkId = c.params[0] as string
    const benchmark = benchmarks.find((b) => b.id === benchmarkId)
    if (benchmark === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(benchmark)
  }),
  http.get(`${apiBaseUrl}/scores`, () => {
    const res: paths['/scores']['get']['responses']['200']['content']['application/json'] =
      teams.map((team) => ({
        teamId: team.id,
        scores: benchmarks
          .filter((b) => b.teamId === team.id)
          .filter((b) => b.status === 'finished')
          .map((b) => ({
            benchmarkId: b.id,
            teamId: team.id,
            score: b.score,
            createdAt: b.createdAt,
          })),
      }))

    return HttpResponse.json(res)
  }),
  http.get(`${apiBaseUrl}/scores/ranking`, () => {
    const sortFn = (
      a: components['schemas']['FinishedBenchmark'],
      b: components['schemas']['FinishedBenchmark'],
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
              .filter((b) => b.teamId === team.id)
              .filter((b) => b.status === 'finished')
              .sort(sortFn)[0],
        )
        .sort(sortFn)
        .map((b, i) => ({
          rank: i + 1,
          teamId: b.teamId,
          score: b.score,
          createdAt: b.createdAt,
        }))

    return HttpResponse.json(res)
  }),
  http.put(`${apiBaseUrl}/admins`, async (c) => {
    type Body = paths['/admins']['put']['requestBody']['content']['application/json']
    const body = (await c.request.json()) as Body

    for (const user of users) {
      user.isAdmin = body.includes(user.id)
    }

    return HttpResponse.json({})
  }),
  http.get(`${apiBaseUrl}/docs`, () => {
    const res: paths['/docs']['get']['responses']['200']['content']['application/json'] = {
      body: docs,
    }

    return HttpResponse.json(res)
  }),
  http.patch(`${apiBaseUrl}/docs`, () => {
    // TODO
  }),
]

// ベンチマークの状態を定期的に更新する
setInterval(() => {
  const waitingBenchmarks = benchmarks.filter((b) => b.status === 'waiting')
  const runningBenchmarks = benchmarks.filter((b) => b.status === 'running')

  // running のベンチマークにログを出力する
  for (const b of runningBenchmarks) {
    if (Math.random() < 0.1) {
      b.log += `${new Date().toISOString()} [INFO] Benchmark is running...\n`
    }
    if (Math.random() < 0.1) {
      b.adminLog += `${new Date().toISOString()} [ERROR] Admin log...\n`
    }
    if (Math.random() < 0.05) {
      b.score += Math.floor(Math.random() * (100 + b.score))
    }
  }

  // running のまま 60 秒経過したら finished にする
  for (const b of runningBenchmarks) {
    if (new Date(b.startedAt).getTime() + 60 * 1000 < Date.now()) {
      const index = benchmarks.findIndex((bb) => bb.id === b.id)
      const result = ['passed', 'failed', 'error'][Math.floor(Math.random() * 3)] as
        | 'passed'
        | 'failed'
        | 'error'
      benchmarks[index] = {
        ...b,
        status: 'finished',
        finishedAt: new Date().toISOString(),
        result,
      }
    }
  }

  // 実行中のベンチマークがなくなったら、waiting のベンチマークを1つ running にする
  if (runningBenchmarks.length === 0) {
    const waitingBenchmark = benchmarks.find((b) => b.status === 'waiting')
    if (waitingBenchmark !== undefined) {
      const index = benchmarks.findIndex((b) => b.id === waitingBenchmark.id)
      benchmarks[index] = {
        ...waitingBenchmark,
        status: 'running',
        startedAt: new Date().toISOString(),
        score: 0,
      }
    }
  }

  // waiting のベンチマークがなくなったら、新しいベンチマークを enqueue する
  const queueBenchmark = waitingBenchmarks.find((b) => b.teamId === teamIds['ikura-cp'])
  if (queueBenchmark === undefined) {
    const instanceId = [
      instanceIds['ikura-cp-1'],
      instanceIds['ikura-cp-2'],
      instanceIds['ikura-cp-3'],
    ][Math.floor(Math.random() * 3)]

    const userId = [userIds.cp20, userIds['ikura-hamu']][Math.floor(Math.random() * 2)]

    benchmarks.push({
      id: uuidv7(),
      instanceId,
      teamId: teamIds['ikura-cp'],
      userId,
      status: 'waiting',
      createdAt: new Date().toISOString(),
      log: '',
      adminLog: '',
    })
  }
}, 100)

// チームのインスタンスの状態を定期的に更新する
setInterval(() => {
  const buildingInstances = instances.filter((i) => i.status === 'building')
  const startingInstances = instances.filter((i) => i.status === 'starting')
  const stoppingInstances = instances.filter((i) => i.status === 'stopping')
  const deletingInstances = instances.filter((i) => i.status === 'deleting')

  // building -> running
  for (const i of buildingInstances) {
    if (Math.random() < 0.2) i.status = 'running'
  }

  // starting -> running
  for (const i of startingInstances) {
    if (Math.random() < 0.3) i.status = 'running'
  }

  // stopping -> stopped
  for (const i of stoppingInstances) {
    if (Math.random() < 0.3) i.status = 'stopped'
  }

  // deleting -> deleted
  for (const i of deletingInstances) {
    if (Math.random() < 0.3) i.status = 'deleted'
  }
}, 500)
