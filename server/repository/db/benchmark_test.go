package db_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
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
