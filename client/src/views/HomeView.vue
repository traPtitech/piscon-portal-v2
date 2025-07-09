<script setup lang="ts">
import { computed } from 'vue'
import PageTitle from '@/components/PageTitle.vue'
import StatusCard from '@/components/StatusCard.vue'
import RankingList from '@/components/RankingList.vue'
import NavigationCard from '@/components/NavigationCard.vue'
import ScoreChart from '@/components/ScoreChart.vue'
import { useMe, useTeamsData, useRanking, useScores } from '@/lib/useServerData'

const { data: me } = useMe()
const { data: teams } = useTeamsData()
const { data: ranking } = useRanking('latest')
const { data: scores } = useScores()

const teamCount = computed(() => teams.value?.length ?? 0)

const myTeam = computed(() => {
  if (!me.value?.teamId || !teams.value) return null
  return teams.value.find((team) => team.id === me.value?.teamId) || null
})

const myTeamRanking = computed(() => {
  if (!myTeam.value || !ranking.value) return null
  return ranking.value.find((item) => item.teamId === myTeam.value?.id) || null
})

const chartScores = computed(() => {
  if (!scores.value) return []

  return scores.value.flatMap((teamScore) =>
    teamScore.scores.map((score) => ({
      teamId: teamScore.teamId,
      score: score.score,
      createdAt: score.createdAt,
    })),
  )
})
</script>

<template>
  <main class="home-container">
    <PageTitle icon="mdi:home">ホーム</PageTitle>

    <section class="overview-section">
      <StatusCard icon="mdi:account-group" title="参加チーム数" :value="teamCount" />
      <StatusCard
        v-if="myTeamRanking"
        icon="mdi:trophy"
        title="チーム順位"
        :value="`${myTeamRanking.rank}位`"
        variant="highlight"
      />
      <StatusCard
        v-if="myTeamRanking"
        icon="mdi:chart-line"
        title="チームスコア"
        :value="myTeamRanking.score"
        variant="primary"
      />
    </section>

    <section class="content-section">
      <div class="content-section-row">
        <div class="ranking-container">
          <h2 class="section-title">スコアランキング</h2>
          <RankingList
            v-if="ranking && ranking.length > 0"
            :ranking="ranking"
            :highlight-team-id="myTeam?.id"
          />
          <div v-else class="empty-state">まだベンチマーク結果がありません</div>
        </div>
        <div class="queue-container">
          <h2 class="section-title">ベンチマークキュー</h2>
          <BenchmarkQueue class="benchmark-queue" />
        </div>
      </div>

      <div class="chart-container">
        <h2 class="section-title">全体スコア推移</h2>
        <div class="chart-wrapper">
          <ScoreChart v-if="chartScores.length > 0" :scores="chartScores" />
          <div v-else class="empty-state">まだベンチマーク結果がありません</div>
        </div>
      </div>
    </section>

    <section class="navigation-section">
      <div class="navigation-grid">
        <NavigationCard
          v-if="myTeam"
          icon="mdi:server"
          title="インスタンス管理"
          description="インスタンスの起動・停止・削除"
          to="/instances"
        />

        <NavigationCard
          icon="mdi:thunder"
          title="ベンチマーク"
          description="ベンチマークの実行と結果確認"
          to="/benches"
        />

        <NavigationCard
          icon="mdi:account-group"
          title="チーム管理"
          description="メンバー管理など"
          to="/team"
        />

        <NavigationCard
          icon="mdi:file-document"
          title="ドキュメント"
          description="ルール・注意点など"
          to="/docs"
        />
      </div>
    </section>
  </main>
</template>

<style scoped>
.home-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  container-type: inline-size;
}

.overview-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.content-section {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.content-section-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  justify-content: space-between;
  align-items: stretch;
}

@container (max-width: 720px) {
  .content-section-row {
    grid-template-columns: 1fr;
  }
}

.ranking-container,
.queue-container,
.chart-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.benchmark-queue {
  max-height: 446px;
  overflow-y: auto;
}

.chart-wrapper {
  height: 400px;
  padding: 1rem;
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  background: var(--ct-white);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--ct-slate-900);
  margin: 0;
}

.navigation-section {
  margin-top: 1rem;
}

.navigation-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: var(--ct-slate-500);
  background: var(--ct-slate-50);
  border-radius: 6px;
  border: 1px solid var(--ct-slate-200);
}
</style>
