package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/utils/optional"
)

type BenchmarkUseCase interface {
	GetBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error)
	CreateBenchmark(ctx context.Context, instanceID uuid.UUID, userID uuid.UUID) (domain.Benchmark, error)
	GetBenchmarks(ctx context.Context) ([]domain.Benchmark, error)
	GetQueuedBenchmarks(ctx context.Context) ([]domain.Benchmark, error)
	GetTeamBenchmarks(ctx context.Context, teamID uuid.UUID) ([]domain.Benchmark, error)
	// StartBenchmark
	// キューの先頭のベンチマークを取得し、ステータスを更新する。
	// 先頭のベンチマークを返すが、キューが空の場合はusecase.ErrNotFoundを返す。
	StartBenchmark(ctx context.Context) (domain.Benchmark, error)
	// SaveBenchmarkProgress
	// ベンチマーク実行中のログやスコアを更新して保存する。
	// 該当のベンチマークが存在しなかった場合、usecase.ErrNotFoundを返す。
	SaveBenchmarkProgress(ctx context.Context, benchmarkID uuid.UUID, benchLog domain.BenchmarkLog, score int64, startedAt time.Time) error
	// FinalizeBenchmark
	// ベンチマークを終了させ、結果を保存してステータスを更新する。
	// 該当のベンチマークが存在しなかった場合、usecase.ErrNotFoundを返す。
	FinalizeBenchmark(ctx context.Context, benchmarkID uuid.UUID, result domain.BenchmarkResult, finishedAt time.Time, errorMessage string) error

	GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error)
}

type benchmarkUseCaseImpl struct {
	repo repository.Repository
}

func NewBenchmarkUseCase(repo repository.Repository) BenchmarkUseCase {
	return &benchmarkUseCaseImpl{repo: repo}
}

func (u *benchmarkUseCaseImpl) GetBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error) {
	benchmark, err := u.repo.FindBenchmark(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Benchmark{}, ErrNotFound
		}
		return domain.Benchmark{}, fmt.Errorf("find benchmark: %v", err)
	}

	return benchmark, nil
}

func (u *benchmarkUseCaseImpl) CreateBenchmark(ctx context.Context, instanceID uuid.UUID, userID uuid.UUID) (domain.Benchmark, error) {
	var benchmark domain.Benchmark

	err := u.repo.Transaction(ctx, func(ctx context.Context) error {
		user, err := u.repo.FindUser(ctx, userID)
		if err != nil {
			return fmt.Errorf("find user: %v", err)
		}
		instance, err := u.repo.FindInstance(ctx, instanceID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return NewUseCaseErrorFromMsg("instance not found")
			}
			return fmt.Errorf("find instance: %v", err)
		}

		benchmark, err = domain.NewBenchmark(instance, user)
		if err != nil {
			return NewUseCaseError(err)
		}

		benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
			TeamID:   optional.From(user.TeamID.UUID),
			StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusWaiting, domain.BenchmarkStatusRunning}),
		})
		if err != nil {
			return fmt.Errorf("get benchmarks: %v", err)
		}
		if len(benchmarks) > 0 {
			return NewUseCaseErrorFromMsg("already exists benchmark")
		}

		err = u.repo.CreateBenchmark(ctx, benchmark)
		if err != nil {
			return fmt.Errorf("create benchmark: %v", err)
		}

		return nil
	})
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("transaction: %v", err)
	}

	return benchmark, nil
}

func (u *benchmarkUseCaseImpl) GetBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{})
	if err != nil {
		return nil, fmt.Errorf("get benchmarks: %v", err)
	}
	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetQueuedBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusWaiting, domain.BenchmarkStatusRunning}),
	})
	if err != nil {
		return nil, fmt.Errorf("get benchmarks: %v", err)
	}

	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetTeamBenchmarks(ctx context.Context, teamID uuid.UUID) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		TeamID: optional.From(teamID),
	})
	if err != nil {
		return nil, fmt.Errorf("get benchmarks: %v", err)
	}

	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	log, err := u.repo.GetBenchmarkLog(ctx, benchmarkID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.BenchmarkLog{}, ErrNotFound
		}
		return domain.BenchmarkLog{}, fmt.Errorf("get benchmark log: %v", err)
	}

	return log, nil
}

func (u *benchmarkUseCaseImpl) StartBenchmark(ctx context.Context) (domain.Benchmark, error) {
	var startedBenchmark domain.Benchmark
	err := u.repo.Transaction(ctx, func(ctx context.Context) error {
		bench, err := u.repo.GetOldestQueuedBenchmark(ctx)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get oldest queued benchmark: %v", err)
		}

		startedBenchmark = domain.Benchmark{
			ID:        bench.ID,
			Instance:  bench.Instance,
			TeamID:    bench.TeamID,
			UserID:    bench.UserID,
			Status:    domain.BenchmarkStatusRunning,
			CreatedAt: bench.CreatedAt,
		}
		err = u.repo.UpdateBenchmark(ctx, bench.ID, startedBenchmark)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("update benchmark: %v", err)
		}

		return nil
	})
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("transaction: %w", err)
	}

	return startedBenchmark, nil
}

func (u *benchmarkUseCaseImpl) SaveBenchmarkProgress(_ context.Context, _ uuid.UUID, _ domain.BenchmarkLog, _ int64, _ time.Time) error {
	// TODO
	return nil
}

func (u *benchmarkUseCaseImpl) FinalizeBenchmark(_ context.Context, _ uuid.UUID, _ domain.BenchmarkResult, _ time.Time, _ string) error {
	// TODO
	return nil
}
