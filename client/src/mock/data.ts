import type { components, paths } from '@/api/openapi'
import { uuidv7 } from 'uuidv7'

export const userIds = {
  cp20: '01943f2e-0dca-7599-9f79-901de49660fd',
  'ikura-hamu': '01943f39-4304-7cf0-b420-2bb5cfde5f1a',
  pirosiki: '01943f39-8171-7db8-96d8-e672a72d14fd',
}

export const teamIds = {
  'ikura-cp': '01943f2e-4e49-782c-9eb8-6a496b509bc8',
  piropiro: '01943f39-9502-76b2-8ace-a1cb3931fcf3',
}

const oneMinute = 60 * 1000
const justBefore = new Date(Date.now() - 5 * oneMinute)

export const users: paths['/users']['get']['responses']['200']['content']['application/json'] = [
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

export const teams: paths['/teams']['get']['responses']['200']['content']['application/json'] = [
  {
    id: teamIds['ikura-cp'],
    name: 'ikura-cp',
    members: [userIds.cp20, userIds['ikura-hamu']],
    createdAt: justBefore.toISOString(),
    githubIds: ['cp-20', 'ikura-hamu'],
  },
  {
    id: teamIds.piropiro,
    name: 'piropiro',
    members: [userIds.pirosiki],
    createdAt: justBefore.toISOString(),
    githubIds: ['pirosiki197'],
  },
]

export const instanceIds = {
  'ikura-cp-1': '01943f39-b74e-7d2d-8a16-9d408bf92579',
  'ikura-cp-2': '01943f4b-bd05-72bd-b603-633e49b94edb',
  'ikura-cp-3': '01943f4c-21a3-7c95-ba5a-63c7ea02b9f0',
  'piropiro-1': '01943f4c-3662-7cac-9674-2b0a39da6725',
}

export const instances: paths['/teams/{teamId}/instances']['get']['responses']['200']['content']['application/json'] =
  [
    {
      id: instanceIds['ikura-cp-1'],
      teamId: teamIds['ikura-cp'],
      serverId: 1,
      privateIPAddress: '192.168.0.1',
      publicIPAddress: '203.0.113.1',
      status: 'running',
      createdAt: justBefore.toISOString(),
    },
    {
      id: instanceIds['ikura-cp-2'],
      teamId: teamIds['ikura-cp'],
      serverId: 2,
      privateIPAddress: '192.168.0.2',
      publicIPAddress: '203.0.113.2',
      status: 'building',
      createdAt: justBefore.toISOString(),
    },
    {
      id: instanceIds['ikura-cp-3'],
      teamId: teamIds['ikura-cp'],
      serverId: 3,
      privateIPAddress: '192.168.0.3',
      publicIPAddress: '203.0.113.3',
      status: 'stopped',
      createdAt: justBefore.toISOString(),
    },
    {
      id: instanceIds['piropiro-1'],
      teamId: teamIds['piropiro'],
      serverId: 1,
      privateIPAddress: '192.168.0.4',
      publicIPAddress: '203.0.113.4',
      status: 'running',
      createdAt: justBefore.toISOString(),
    },
  ]

export const benchmarks: components['schemas']['BenchmarkAdminResult'][] = [
  {
    id: '01943f67-d9ed-7bbb-81eb-20f81391ffea',
    instanceId: instanceIds['ikura-cp-1'],
    teamId: teamIds['ikura-cp'],
    userId: userIds.cp20,
    status: 'finished',
    createdAt: justBefore.toISOString(),
    startedAt: justBefore.toISOString(),
    finishedAt: new Date(justBefore.getTime() + oneMinute).toISOString(),
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
    createdAt: new Date(justBefore.getTime() + 2 * oneMinute).toISOString(),
    startedAt: new Date(justBefore.getTime() + 2 * oneMinute).toISOString(),
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
    createdAt: new Date(justBefore.getTime() + 3 * oneMinute).toISOString(),
    log: '',
    adminLog: '',
  },
  {
    id: '01943f6b-7276-79a9-9bd3-69d1d5a9cb3d',
    instanceId: instanceIds['piropiro-1'],
    teamId: teamIds['piropiro'],
    userId: userIds.pirosiki,
    status: 'finished',
    createdAt: new Date(justBefore.getTime()).toISOString(),
    startedAt: new Date(justBefore.getTime()).toISOString(),
    finishedAt: new Date(justBefore.getTime() + oneMinute).toISOString(),
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
    createdAt: new Date(justBefore.getTime() + 2 * oneMinute).toISOString(),
    startedAt: new Date(justBefore.getTime() + 2 * oneMinute).toISOString(),
    finishedAt: new Date(justBefore.getTime() + 3 * oneMinute).toISOString(),
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
    createdAt: new Date(justBefore.getTime() + 3 * oneMinute).toISOString(),
    startedAt: new Date(justBefore.getTime() + 3 * oneMinute).toISOString(),
    finishedAt: new Date(justBefore.getTime() + 4 * oneMinute).toISOString(),
    score: 1000,
    log: '',
    adminLog: '',
    result: 'passed',
  },
]

export const docs = `# ドキュメント

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

const DUMMY_TEAM_COUNT = 10

type Team = components['schemas']['Team']
const dummyTeams: Team[] = [...Array(DUMMY_TEAM_COUNT).keys()].map((i) => ({
  id: uuidv7(),
  name: `ダミーチーム${i + 1}`,
  members: [],
  createdAt: justBefore.toISOString(),
  githubIds: [],
}))

teams.push(...dummyTeams)

type Instance = components['schemas']['Instance']
const dummyInstances: Instance[] = [...Array(DUMMY_TEAM_COUNT).keys()].map((i) => ({
  id: uuidv7(),
  teamId: dummyTeams[i].id,
  serverId: 1,
  privateIPAddress: `192.168.0.${i + 10}`,
  publicIPAddress: `203.0.113.${i + 10}`,
  status: 'running',
  createdAt: justBefore.toISOString(),
}))

instances.push(...dummyInstances)

type Benchmark = components['schemas']['BenchmarkAdminResult']
const dummyBenchmarks: Benchmark[] = [...Array(DUMMY_TEAM_COUNT).keys()].map((i) => ({
  id: uuidv7(),
  instanceId: dummyInstances[i].id,
  teamId: dummyTeams[i].id,
  userId: userIds.cp20,
  status: 'finished',
  createdAt: justBefore.toISOString(),
  startedAt: justBefore.toISOString(),
  finishedAt: new Date(justBefore.getTime() + oneMinute).toISOString(),
  score: 1000 * 2 ** i,
  log: '',
  adminLog: '',
  result: 'passed',
}))

benchmarks.push(...dummyBenchmarks)
