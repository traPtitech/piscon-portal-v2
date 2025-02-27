package usecase

import (
	"context"
	"errors"
	"fmt"

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
		return domain.Benchmark{}, err
	}

	return benchmark, nil
}

func (u *benchmarkUseCaseImpl) CreateBenchmark(ctx context.Context, instanceID uuid.UUID, userID uuid.UUID) (domain.Benchmark, error) {
	user, err := u.repo.FindUser(ctx, userID)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("find user: %v", err)
	}
	instance, err := u.repo.FindInstance(ctx, instanceID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Benchmark{}, NewUseCaseErrorFromMsg("instance not found")
		}
		return domain.Benchmark{}, fmt.Errorf("find instance: %v", err)
	}

	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		TeamID:   optional.From(user.TeamID.UUID),
		StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusWaiting, domain.BenchmarkStatusRunning}),
	})
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("get benchmarks: %v", err)
	}
	if len(benchmarks) > 0 {
		return domain.Benchmark{}, NewUseCaseErrorFromMsg("already exists benchmark")
	}

	benchmark, err := domain.NewBenchmark(instance, user)
	if err != nil {
		return domain.Benchmark{}, NewUseCaseError(err)
	}
	err = u.repo.CreateBenchmark(ctx, benchmark)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("create benchmark: %v", err)
	}

	return benchmark, nil
}

func (u *benchmarkUseCaseImpl) GetBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{})
	if err != nil {
		return nil, err
	}
	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetQueuedBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		StatusIn: optional.From([]domain.BenchmarkStatus{domain.BenchmarkStatusWaiting, domain.BenchmarkStatusRunning}),
	})
	if err != nil {
		return nil, err
	}

	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetTeamBenchmarks(ctx context.Context, teamID uuid.UUID) ([]domain.Benchmark, error) {
	benchmarks, err := u.repo.GetBenchmarks(ctx, repository.BenchmarkQuery{
		TeamID: optional.From(teamID),
	})
	if err != nil {
		return nil, err
	}

	return benchmarks, nil
}

func (u *benchmarkUseCaseImpl) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	log, err := u.repo.GetBenchmarkLog(ctx, benchmarkID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.BenchmarkLog{}, ErrNotFound
		}
		return domain.BenchmarkLog{}, err
	}

	return log, nil
}
