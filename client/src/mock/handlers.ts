import { HttpResponse } from 'msw'
import { uuidv7 } from 'uuidv7'
import { createOpenApiHttp } from 'openapi-msw'
import type { components, paths } from '@/api/openapi'
import { apiBaseUrl } from '@/api'
import { users, teams, instances, benchmarks, docs, updateDocs } from './data'
import { startJobs } from '@/mock/jobs'

const http = createOpenApiHttp<paths>({
  baseUrl: apiBaseUrl,
})

startJobs()

export const handlers = [
  http.get(`/oauth2/code`, () => {
    // TODO
  }),
  http.get(`/oauth2/callback`, () => {
    // TODO
  }),
  http.post(`/oauth2/logout`, () => {
    // TODO
  }),
  http.get(`/users`, () => HttpResponse.json(users)),
  http.get(`/users/me`, () => HttpResponse.json(users[0])),
  http.get(`/teams`, () => HttpResponse.json(teams)),
  http.post(`/teams`, async (c) => {
    const body = await c.request.json()
    if (body === undefined) {
      return HttpResponse.json({ message: 'Bad request' }, { status: 400 })
    }
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
  http.get('/teams/{teamId}', (c) => {
    const teamId = c.params.teamId
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(team)
  }),
  http.patch('/teams/{teamId}', async (c) => {
    const teamId = c.params.teamId
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })

    const body = await c.request.json()
    if (body === undefined) {
      return HttpResponse.json({ message: 'Bad request' }, { status: 400 })
    }

    if (body.name !== undefined) {
      team.name = body.name
    }

    if (body.members !== undefined) {
      const oldMembers = team.members
      team.members = body.members

      const removedMembers = oldMembers.filter((m) => !body.members!.includes(m))
      for (const member of removedMembers) {
        const user = users.find((u) => u.id === member)
        if (user !== undefined) user.teamId = undefined
      }
      const addedMembers = body.members.filter((m) => !oldMembers.includes(m))
      for (const member of addedMembers) {
        const user = users.find((u) => u.id === member)
        if (user !== undefined) user.teamId = teamId
      }
    }

    return HttpResponse.json(team)
  }),
  http.get('/teams/{teamId}/instances', (c) => {
    const teamId = c.params.teamId
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const res = instances.filter((i) => i.teamId === teamId)
    return HttpResponse.json(res)
  }),
  http.post('/teams/{teamId}/instances', (c) => {
    const teamId = c.params.teamId
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
  http.delete('/teams/{teamId}/instances/{instanceId}', (c) => {
    const teamId = c.params.teamId
    const instanceId = c.params.instanceId
    const index = instances.findIndex((i) => i.teamId === teamId && i.id === instanceId)
    instances[index] = { ...instances[index], status: 'deleting' }
  }),
  http.patch('/teams/{teamId}/instances/{instanceId}', async (c) => {
    const teamId = c.params.teamId
    const instanceId = c.params.instanceId
    const index = instances.findIndex((i) => i.teamId === teamId && i.id === instanceId)
    const body = await c.request.json()
    const operation = body.operation
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
  http.get(`/instances`, () => HttpResponse.json(instances)),
  http.post(`/benchmarks`, async (c) => {
    const body = await c.request.json()
    if (body === undefined) {
      return HttpResponse.json({ message: 'Bad request' }, { status: 400 })
    }

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
  http.get(`/benchmarks`, () => HttpResponse.json(benchmarks)),
  http.get(`/benchmarks/queue`, () => {
    const res = benchmarks.filter((b) => b.status !== 'finished')
    return HttpResponse.json(res)
  }),
  http.get('/teams/{teamId}/benchmarks', (c) => {
    const teamId = c.params.teamId
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const res = benchmarks.filter((b) => b.teamId === teamId)
    return HttpResponse.json(res)
  }),
  http.get('/teams/{teamId}/benchmarks/{benchmarkId}', (c) => {
    const teamId = c.params.teamId
    const team = teams.find((t) => t.id === teamId)
    if (team === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const benchmarkId = c.params.benchmarkId
    const benchmark = benchmarks.find((b) => b.teamId === teamId && b.id === benchmarkId)
    if (benchmark === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    const { adminLog: _, ...res } = benchmark
    return HttpResponse.json(res)
  }),
  http.get('/benchmarks/{benchmarkId}', (c) => {
    const benchmarkId = c.params.benchmarkId
    const benchmark = benchmarks.find((b) => b.id === benchmarkId)
    if (benchmark === undefined) return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    return HttpResponse.json(benchmark)
  }),
  http.get(`/scores`, () => {
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
  http.get(`/scores/ranking`, () => {
    const sortFn = (
      a: components['schemas']['FinishedBenchmark'],
      b: components['schemas']['FinishedBenchmark'],
    ) => {
      // スコアの降順
      if (a.score < b.score) return 1
      if (a.score > b.score) return -1
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
              .slice(-1)[0],
        )
        .filter((b) => b !== undefined)
        .sort(sortFn)
        .map((b, i) => ({
          rank: i + 1,
          teamId: b.teamId,
          score: b.score,
          createdAt: b.createdAt,
        }))

    return HttpResponse.json(res)
  }),
  http.put(`/admins`, async (c) => {
    const body = await c.request.json()

    for (const user of users) {
      user.isAdmin = body.includes(user.id)
    }

    return HttpResponse.json({})
  }),
  http.get(`/docs`, () => {
    const res: paths['/docs']['get']['responses']['200']['content']['application/json'] = {
      body: docs,
    }

    return HttpResponse.json(res)
  }),
  http.patch(`/docs`, async (c) => {
    const body = await c.request.json()
    if (body === undefined) {
      return HttpResponse.json({ message: 'Bad request' }, { status: 400 })
    }

    updateDocs(body.body)
    return HttpResponse.json({ body: docs }, { status: 200 })
  }),
]
