package db_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
	"github.com/traPtitech/piscon-portal-v2/server/utils/optional"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
)

func TestFindBenchmark(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	id := uuid.New()
	teamID := uuid.New()
	instanceID := uuid.New()
	userID := uuid.New()
	benchmark := domain.Benchmark{
		ID: id,
		Instance: domain.Instance{
			ID: instanceID,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "provider-instance-id",
				Status:             domain.InstanceStatusRunning,
			},
		},
		TeamID:     teamID,
		UserID:     userID,
		Status:     domain.BenchmarkStatusWaiting,
		CreatedAt:  time.Now(),
		StartedAt:  nil,
		FinishedAt: nil,
	}
	mustMakeInstance(t, db, benchmark.Instance)
	mustMakeBenchmark(t, db, benchmark)

	got, err := repo.FindBenchmark(t.Context(), id)
	assert.NoError(t, err)

	testutil.CompareBenchmark(t, benchmark, got)
}

func TestCreateBenchmark(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	id := uuid.New()
	teamID := uuid.New()
	instanceID := uuid.New()
	userID := uuid.New()
	instance := domain.Instance{
		ID: instanceID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
		},
	}
	benchmark := domain.Benchmark{
		ID:         id,
		Instance:   instance,
		TeamID:     teamID,
		UserID:     userID,
		Status:     domain.BenchmarkStatusWaiting,
		CreatedAt:  time.Now(),
		StartedAt:  nil,
		FinishedAt: nil,
	}

	mustMakeInstance(t, db, instance)

	err := repo.CreateBenchmark(t.Context(), benchmark)
	assert.NoError(t, err)

	got, err := repo.FindBenchmark(t.Context(), id)
	assert.NoError(t, err)

	testutil.CompareBenchmark(t, benchmark, got)
}

func TestGetAllBenchmarks(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	teamID := uuid.New()
	userID := uuid.New()
	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
		},
	}

	benchmarks := []domain.Benchmark{
		{
			ID:         uuid.New(),
			Instance:   instance,
			TeamID:     teamID,
			UserID:     userID,
			Status:     domain.BenchmarkStatusWaiting,
			CreatedAt:  time.Now(),
			StartedAt:  nil,
			FinishedAt: nil,
		},
		{
			ID:         uuid.New(),
			Instance:   instance,
			TeamID:     teamID,
			UserID:     userID,
			Status:     domain.BenchmarkStatusFinished,
			CreatedAt:  time.Now(),
			StartedAt:  ptr.Of(time.Now()),
			FinishedAt: ptr.Of(time.Now()),
			Score:      100,
			Result:     ptr.Of(domain.BenchmarkResultStatusPassed),
		},
	}
	mustMakeInstance(t, db, instance)
	for _, benchmark := range benchmarks {
		mustMakeBenchmark(t, db, benchmark)
	}

	got, err := repo.GetBenchmarks(t.Context(), repository.BenchmarkQuery{})
	assert.NoError(t, err)

	testutil.CompareBenchmarks(t, benchmarks, got)
}

func TestGetQueuedBenchmarks(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	teamID := uuid.New()
	userID := uuid.New()
	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
		},
	}

	benchmarks := []domain.Benchmark{
		{
			ID:         uuid.New(),
			Instance:   instance,
			TeamID:     teamID,
			UserID:     userID,
			Status:     domain.BenchmarkStatusWaiting,
			CreatedAt:  time.Now(),
			StartedAt:  nil,
			FinishedAt: nil,
		},
		{
			ID:         uuid.New(),
			Instance:   instance,
			TeamID:     teamID,
			UserID:     userID,
			Status:     domain.BenchmarkStatusFinished,
			CreatedAt:  time.Now(),
			StartedAt:  ptr.Of(time.Now()),
			FinishedAt: ptr.Of(time.Now()),
			Score:      100,
			Result:     ptr.Of(domain.BenchmarkResultStatusPassed),
		},
	}
	mustMakeInstance(t, db, instance)
	for _, benchmark := range benchmarks {
		mustMakeBenchmark(t, db, benchmark)
	}

	got, err := repo.GetBenchmarks(t.Context(), repository.BenchmarkQuery{
		StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusWaiting, domain.BenchmarkStatusRunning}),
	})
	assert.NoError(t, err)

	testutil.CompareBenchmarks(t, benchmarks[:1], got)
}

func TestGetBenchmarkLog(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	benchmarkID := uuid.New()
	benchmark := domain.Benchmark{
		ID:     benchmarkID,
		Status: domain.BenchmarkStatusFinished,
		Instance: domain.Instance{
			ID: uuid.New(),
			Infra: domain.InfraInstance{
				ProviderInstanceID: "provider-instance-id",
				Status:             domain.InstanceStatusRunning,
			},
		},
		CreatedAt: time.Now(),
	}
	benchmarkLog := domain.BenchmarkLog{
		UserLog:  "user log",
		AdminLog: "admin log",
	}
	mustMakeInstance(t, db, benchmark.Instance)
	mustMakeBenchmark(t, db, benchmark)
	mustMakeBenchmarkLog(t, db, benchmarkID, benchmarkLog)

	got, err := repo.GetBenchmarkLog(t.Context(), benchmarkID)
	assert.NoError(t, err)

	testutil.CompareBenchmarkLog(t, benchmarkLog, got)
}

func TestGetOldestQueuedBenchmark(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: uuid.New(),
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
			PublicIP:           ptr.Of("192.168.1.1"),
			PrivateIP:          ptr.Of("10.0.0.1"),
		},
	}
	waitingBench := domain.Benchmark{
		ID:        uuid.New(),
		Instance:  instance,
		TeamID:    uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.BenchmarkStatusWaiting,
		CreatedAt: time.Now(),
	}
	waitingBench2 := domain.Benchmark{
		ID:        uuid.New(),
		Instance:  instance,
		TeamID:    uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.BenchmarkStatusWaiting,
		CreatedAt: time.Now().Add(-time.Hour),
	}
	finishedBenchmark := domain.Benchmark{
		ID:         uuid.New(),
		Instance:   instance,
		TeamID:     uuid.New(),
		UserID:     uuid.New(),
		Status:     domain.BenchmarkStatusFinished,
		CreatedAt:  time.Now().Add(-2 * time.Hour),
		StartedAt:  ptr.Of(time.Now().Add(-time.Hour)),
		FinishedAt: ptr.Of(time.Now().Add(-time.Hour)),
	}

	testCases := map[string]struct {
		benchmarks []domain.Benchmark
		expected   domain.Benchmark
		err        error
	}{
		"1個しかベンチマークが無い時に正常に取得できる": {
			benchmarks: []domain.Benchmark{waitingBench},
			expected:   waitingBench,
		},
		"waitingが2個あっても古い方を取得できる": {
			benchmarks: []domain.Benchmark{waitingBench, waitingBench2},
			expected:   waitingBench2,
		},
		"waitingでないベンチマークがあっても正しく取得できる": {
			benchmarks: []domain.Benchmark{waitingBench, finishedBenchmark},
			expected:   waitingBench,
		},
		"キューが空なのでErrNotFound": {
			err: repository.ErrNotFound,
		},
		"waitingが無いのでErrNotFound": {
			benchmarks: []domain.Benchmark{finishedBenchmark},
			err:        repository.ErrNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			instances := make([]domain.Instance, 0, len(testCase.benchmarks))
			for _, bench := range testCase.benchmarks {
				instances = append(instances, bench.Instance)
			}
			// instanceの重複除去
			instances = slices.CompactFunc(instances, func(a, b domain.Instance) bool {
				return a.ID == b.ID
			})
			for _, instance := range instances {
				mustMakeInstance(t, db, instance)
			}
			for _, bench := range testCase.benchmarks {
				mustMakeBenchmark(t, db, bench)
			}

			// 他のテストケースに影響を与えないために削除
			t.Cleanup(func() {
				ctx := context.Background()
				if len(testCase.benchmarks) != 0 {
					benchmarkIDs := make([]string, 0, len(testCase.benchmarks))
					for _, bench := range testCase.benchmarks {
						benchmarkIDs = append(benchmarkIDs, bench.ID.String())
					}
					_, err := models.Benchmarks.Delete(models.DeleteWhere.Benchmarks.ID.In(benchmarkIDs...)).Exec(ctx, db)
					require.NoError(t, err)
				}
				if len(instances) != 0 {
					instanceIDs := make([]string, 0, len(instances))
					for _, instance := range instances {
						instanceIDs = append(instanceIDs, instance.ID.String())
					}
					_, err := models.Instances.Delete(models.DeleteWhere.Instances.ID.In(instanceIDs...)).Exec(ctx, db)
					require.NoError(t, err)
				}
			})

			got, err := repo.GetOldestQueuedBenchmark(t.Context())

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			testutil.CompareBenchmark(t, testCase.expected, got)
		})
	}
}

func TestUpdateBenchmark(t *testing.T) {
	t.Parallel()

	repo, testDB := setupRepository(t)

	benchID := uuid.New()
	benchID2 := uuid.New()
	benchID3 := uuid.New()
	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: uuid.New(),
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
			PrivateIP:          ptr.Of("0.0.0.0"),
			PublicIP:           ptr.Of("0.0.0.0"),
		},
	}

	createdAt := time.Now()
	startedAt := time.Now().Add(time.Minute)
	finishedAt := time.Now().Add(time.Hour)

	testCases := map[string]struct {
		id          uuid.UUID
		beforeBench domain.Benchmark
		afterBench  domain.Benchmark
		err         error
	}{
		"startedAtとScoreを更新できる": {
			id: benchID,
			beforeBench: domain.Benchmark{
				ID:        benchID,
				Instance:  instance,
				TeamID:    uuid.New(),
				UserID:    uuid.New(),
				Status:    domain.BenchmarkStatusWaiting,
				CreatedAt: createdAt,
			},
			afterBench: domain.Benchmark{
				ID:        benchID,
				Instance:  instance,
				TeamID:    uuid.New(),
				UserID:    uuid.New(),
				Status:    domain.BenchmarkStatusRunning,
				CreatedAt: time.Now(),
				StartedAt: &startedAt,
				Score:     100,
			},
		},
		"finishedAtとResultを更新できる": {
			id: benchID2,
			beforeBench: domain.Benchmark{
				ID:        benchID2,
				Instance:  instance,
				TeamID:    uuid.New(),
				UserID:    uuid.New(),
				Status:    domain.BenchmarkStatusRunning,
				CreatedAt: time.Now(),
				StartedAt: &startedAt,
				Score:     100,
			},
			afterBench: domain.Benchmark{
				ID:         benchID2,
				Instance:   instance,
				TeamID:     uuid.New(),
				UserID:     uuid.New(),
				Status:     domain.BenchmarkStatusRunning,
				CreatedAt:  time.Now(),
				StartedAt:  &startedAt,
				FinishedAt: &finishedAt,
				Score:      300,
				Result:     ptr.Of(domain.BenchmarkResultStatusPassed),
			},
		},
		"statusを更新できる": {
			id: benchID3,
			beforeBench: domain.Benchmark{
				ID:        benchID3,
				Instance:  instance,
				TeamID:    uuid.New(),
				UserID:    uuid.New(),
				Status:    domain.BenchmarkStatusWaiting,
				CreatedAt: createdAt,
			},
			afterBench: domain.Benchmark{
				ID:        benchID3,
				Instance:  instance,
				TeamID:    uuid.New(),
				UserID:    uuid.New(),
				Status:    domain.BenchmarkStatusReadying,
				CreatedAt: createdAt,
			},
		},
	}

	mustMakeInstance(t, testDB, instance)
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {

			mustMakeBenchmark(t, testDB, testCase.beforeBench)

			err := repo.UpdateBenchmark(t.Context(), testCase.beforeBench.ID, testCase.afterBench)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			bench, err := models.Benchmarks.Query(models.SelectWhere.Benchmarks.ID.EQ(testCase.beforeBench.ID.String())).One(t.Context(), testDB)
			assert.NoError(t, err)

			assert.Equal(t, testCase.afterBench.ID.String(), bench.ID)
			assert.Equal(t, testCase.afterBench.Instance.ID.String(), bench.InstanceID)
			assert.Equal(t, testCase.afterBench.TeamID.String(), bench.TeamID)
			assert.Equal(t, testCase.afterBench.UserID.String(), bench.UserID)
			assert.Equal(t, testCase.afterBench.Score, bench.Score)
			status, err := db.FromDomainBenchmarkStatus(testCase.afterBench.Status)
			require.NoError(t, err)
			assert.Equal(t, status, bench.Status)
			result, err := db.FromDomainBenchmarkResult(testCase.afterBench.Result)
			require.NoError(t, err)
			assert.Equal(t, result, db.SQLNullToPtr(bench.Result))
			assert.WithinDuration(t, testCase.afterBench.CreatedAt, bench.CreatedAt, time.Second)
			if testCase.afterBench.StartedAt != nil {
				assert.WithinDuration(t, *testCase.afterBench.StartedAt, bench.StartedAt.V, time.Second)
			} else {
				assert.Nil(t, db.SQLNullToPtr(bench.StartedAt))
			}
			if testCase.afterBench.FinishedAt != nil {
				assert.WithinDuration(t, *testCase.afterBench.FinishedAt, bench.FinishedAt.V, time.Second)
			} else {
				assert.Nil(t, db.SQLNullToPtr(bench.FinishedAt))
			}
		})
	}
}

func TestUpdateBenchmarkLog(t *testing.T) {
	t.Parallel()

	repo, testDB := setupRepository(t)

	benchID := uuid.New()
	benchID2 := uuid.New()

	benchLog := domain.BenchmarkLog{
		UserLog:  "user log",
		AdminLog: "admin log",
	}

	instanceID := uuid.New()

	mustMakeInstance(t, testDB, domain.Instance{
		ID: instanceID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
		},
	})
	mustMakeBenchmark(t, testDB, domain.Benchmark{
		ID:        benchID,
		Status:    domain.BenchmarkStatusRunning,
		CreatedAt: time.Now(),
		Instance: domain.Instance{
			ID: instanceID,
		},
	})
	mustMakeBenchmark(t, testDB, domain.Benchmark{
		ID:        benchID2,
		Status:    domain.BenchmarkStatusRunning,
		CreatedAt: time.Now(),
		Instance: domain.Instance{
			ID: instanceID,
		},
	})
	mustMakeBenchmarkLog(t, testDB, benchID, benchLog)

	testCases := map[string]struct {
		benchID  uuid.UUID
		benchLog domain.BenchmarkLog
	}{
		"すでにあるlogを更新する": {
			benchID: benchID,
			benchLog: domain.BenchmarkLog{
				UserLog:  "updated user log",
				AdminLog: "updated admin log",
			},
		},
		"新しいlogを追加する": {
			benchID: benchID2,
			benchLog: domain.BenchmarkLog{
				UserLog:  "new user log",
				AdminLog: "new admin log",
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := repo.UpdateBenchmarkLog(t.Context(), testCase.benchID, testCase.benchLog)
			assert.NoError(t, err)

			got, err := models.FindBenchmarkLog(t.Context(), testDB, testCase.benchID.String())
			assert.NoError(t, err)

			assert.Equal(t, testCase.benchLog.UserLog, got.UserLog)
			assert.Equal(t, testCase.benchLog.AdminLog, got.AdminLog)
		})
	}
}

func TestGetRanking(t *testing.T) {
	t.Parallel()

	repo, testDB := setupRepository(t)

	team1 := uuid.New()
	team2 := uuid.New()
	team3 := uuid.New()
	score0 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team1,
		Score:       300,
		CreatedAt:   time.Now(),
	}
	score1 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team1,
		Score:       200,
		CreatedAt:   time.Now(),
	}
	score2 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team1,
		Score:       100,
		CreatedAt:   time.Now().Add(-time.Hour),
	}
	score3 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team2,
		Score:       100,
		CreatedAt:   time.Now().Add(-time.Hour),
	}
	score4 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team2,
		Score:       150,
		CreatedAt:   time.Now().Add(-time.Hour * 2),
	}
	score5 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team3,
		Score:       100,
		CreatedAt:   time.Now().Add(-time.Hour * 2),
	}
	score6 := domain.Score{
		BenchmarkID: uuid.New(),
		TeamID:      team3,
		Score:       150,
		CreatedAt:   time.Now().Add(-time.Hour * 3),
	}
	scores := []domain.Score{score1, score2, score3, score4, score5, score6}

	instanceID := uuid.New()

	mustMakeInstance(t, testDB, domain.Instance{
		ID: instanceID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "provider-instance-id",
			Status:             domain.InstanceStatusRunning,
		},
	})
	mustMakeTeam(t, testDB, domain.Team{ID: team1, CreatedAt: time.Now()})
	mustMakeTeam(t, testDB, domain.Team{ID: team2, CreatedAt: time.Now()})
	mustMakeTeam(t, testDB, domain.Team{ID: team3, CreatedAt: time.Now()})
	for _, score := range scores {
		mustMakeBenchmark(t, testDB, domain.Benchmark{
			ID:        score.BenchmarkID,
			TeamID:    score.TeamID,
			Score:     score.Score,
			CreatedAt: score.CreatedAt,
			Status:    domain.BenchmarkStatusFinished,
			Instance:  domain.Instance{ID: instanceID},
		})
	}
	mustMakeBenchmark(t, testDB, domain.Benchmark{
		ID:        score0.BenchmarkID,
		TeamID:    score0.TeamID,
		Score:     score0.Score,
		CreatedAt: score0.CreatedAt,
		Status:    domain.BenchmarkStatusRunning,
		Instance:  domain.Instance{ID: instanceID},
	})

	testCases := map[string]struct {
		query    repository.RankingQuery
		expected []domain.Score
	}{
		"最新のスコア高い順": {
			query:    repository.RankingQuery{OrderBy: domain.RankingOrderByLatestScore},
			expected: []domain.Score{score1, score5, score3},
		},
		"最高点が高い順": {
			query:    repository.RankingQuery{OrderBy: domain.RankingOrderByHighestScore},
			expected: []domain.Score{score1, score6, score4},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := repo.GetRanking(t.Context(), testCase.query)

			assert.NoError(t, err)

			assert.Len(t, got, len(testCase.expected))
			for i, score := range got {
				assert.Equal(t, testCase.expected[i].BenchmarkID, score.BenchmarkID)
				assert.Equal(t, testCase.expected[i].TeamID, score.TeamID)
				assert.Equal(t, testCase.expected[i].Score, score.Score)
				assert.WithinDuration(t, testCase.expected[i].CreatedAt, score.CreatedAt, time.Second)
			}
		})
	}

}
