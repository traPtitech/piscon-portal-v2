package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/utils/optional"
)

type BenchmarkRepository interface {
	// CreateInstance creates an instance.
	CreateBenchmark(ctx context.Context, benchmark domain.Benchmark) error
	// GetBenchmarks returns all benchmarks. If query is set, it filters the benchmarks.
	// The returned benchmarks are sorted by CreatedAt in ascending order.
	GetBenchmarks(ctx context.Context, query BenchmarkQuery) ([]domain.Benchmark, error)
	// FindBenchmark finds a benchmark by id. If the benchmark is not found, it returns [ErrNotFound].
	FindBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error)
	// GetBenchmarkLog returns a benchmark log.
	GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error)
	// GetOldestQueuedBenchmark returns the oldest queued benchmark.
	// If there are no queued benchmarks, it returns [ErrNotFound].
	GetOldestQueuedBenchmark(ctx context.Context) (domain.Benchmark, error)
	// UpdateBenchmark updates a benchmark record.
	UpdateBenchmark(ctx context.Context, id uuid.UUID, benchmark domain.Benchmark) error
	// UpdateBenchmarkLog updates a benchmark log. If not exists, it creates a new one.
	UpdateBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID, log domain.BenchmarkLog) error
	// GetRanking returns the ranking of score. It does not contains teams without score.
	GetRanking(ctx context.Context, query RankingQuery) ([]domain.Score, error)
}

type BenchmarkQuery struct {
	TeamID   optional.Of[uuid.UUID]
	StatusIn optional.Of[[]domain.BenchmarkStatus]
}

type RankingQuery struct {
	// OrderBy は、スコアの並び順を指定する。
	// ここで指定した値が同じ場合は、古いものが優先される。
	OrderBy domain.RankingOrderBy
}
