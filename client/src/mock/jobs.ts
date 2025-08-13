import {
  benchmarks,
  dummyInstances,
  dummyTeams,
  instances,
  userIds,
  benchmarkReadyingAt,
} from '@/mock/data'
import { uuidv7 } from 'uuidv7'

const BENCHMARKER_CONCURRENCY = 3

export const startJobs = () => {
  // ベンチマークの状態を定期的に更新する
  setInterval(() => {
    const waitingBenchmarks = benchmarks.filter((b) => b.status === 'waiting')
    const runningBenchmarks = benchmarks.filter((b) => b.status === 'running')
    const readyingBenchmarks = benchmarks.filter((b) => b.status === 'readying')

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

    // running のまま60秒経過したら finished にする
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

    // 実行中のベンチマークが BENCHMARKER_CONCURRENCY 未満なら、waiting のベンチマークを readying にする
    if (runningBenchmarks.length + readyingBenchmarks.length < BENCHMARKER_CONCURRENCY) {
      const waitingBenchmark = benchmarks.find((b) => b.status === 'waiting')
      if (waitingBenchmark !== undefined) {
        const index = benchmarks.findIndex((b) => b.id === waitingBenchmark.id)
        benchmarks[index] = {
          ...waitingBenchmark,
          status: 'readying',
        }
        benchmarkReadyingAt.push({ id: waitingBenchmark.id, readyingAt: new Date() })
      }
    }

    // readyingになって3秒以上たったベンチマークがあったら、runningにする
    for (const b of readyingBenchmarks) {
      const readyingAt = benchmarkReadyingAt.find((bb) => bb.id === b.id)?.readyingAt
      if (readyingAt && readyingAt?.getTime() + 3000 > Date.now()) {
        continue
      }
      const index = benchmarks.findIndex((bb) => bb.id === b.id)
      benchmarks[index] = {
        ...b,
        status: 'running',
        startedAt: new Date().toISOString(),
        score: 0,
      }
    }

    // waiting のベンチマークがなくなったら、新しいベンチマークを enqueue する
    for (const [i, dummyTeam] of dummyTeams.entries()) {
      const queueBenchmark = waitingBenchmarks.find((b) => b.teamId === dummyTeam.id)
      if (queueBenchmark === undefined && Math.random() < 0.01) {
        benchmarks.push({
          id: uuidv7(),
          instanceId: dummyInstances[i].id,
          teamId: dummyTeam.id,
          userId: userIds.cp20,
          status: 'waiting',
          createdAt: new Date().toISOString(),
          log: '',
          adminLog: '',
        })
      }
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
}
