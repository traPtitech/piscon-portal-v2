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
			ID:     instanceID,
			Status: domain.InstanceStatusRunning,
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
		ID:     instanceID,
		Status: domain.InstanceStatusRunning,
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
		Status: domain.InstanceStatusRunning,
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
		Status: domain.InstanceStatusRunning,
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
			ID:     uuid.New(),
			Status: domain.InstanceStatusRunning,
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
		ID:             uuid.New(),
		Status:         domain.InstanceStatusRunning,
		PublicIP:       "192.168.1.1",
		PrivateIP:      "10.0.0.1",
		TeamID:         uuid.New(),
		InstanceNumber: 1,
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
			// テーブル全体を見るようなテストなので、t.Paralle()はできない

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
