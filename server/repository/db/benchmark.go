package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
)

func (r *Repository) FindBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error) {
	benchmark, err := models.Benchmarks.
		Query(
			models.PreloadBenchmarkInstance(),
			models.SelectWhere.Benchmarks.ID.EQ(id.String()),
		).
		One(ctx, r.executor(ctx))
	if err != nil {
		return domain.Benchmark{}, err
	}

	return toDomainBenchmark(benchmark)
}

func (r *Repository) CreateBenchmark(ctx context.Context, benchmark domain.Benchmark) error {
	_, err := models.Benchmarks.Insert(&models.BenchmarkSetter{
		ID:         omit.From(benchmark.ID.String()),
		InstanceID: omit.From(benchmark.Instance.ID.String()),
		TeamID:     omit.From(benchmark.TeamID.String()),
		UserID:     omit.From(benchmark.UserID.String()),
		Status:     omit.From(fromDomainBenchmarkStatus(benchmark.Status)),
		CreatedAt:  omit.From(benchmark.CreatedAt),
	}).Exec(ctx, r.executor(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBenchmarks(ctx context.Context) (domain.Benchmarks, error) {
	benchmarks, err := models.Benchmarks.Query(models.PreloadBenchmarkInstance()).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, err
	}

	res := make(domain.Benchmarks, 0, len(benchmarks))
	for _, b := range benchmarks {
		benchmark, err := toDomainBenchmark(b)
		if err != nil {
			return nil, err
		}
		res = append(res, benchmark)
	}

	return res, nil
}

func (r *Repository) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	benchmarkLogs, err := models.FindBenchmarkLog(ctx, r.executor(ctx), benchmarkID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BenchmarkLog{}, repository.ErrNotFound
		}
		return domain.BenchmarkLog{}, err
	}

	return toDomainBenchmarkLog(benchmarkLogs)
}

func fromDomainBenchmarkStatus(status domain.BenchmarkStatus) models.BenchmarksStatus {
	switch status {
	case domain.BenchmarkStatusWaiting:
		return models.BenchmarksStatusWaiting
	case domain.BenchmarkStatusRunning:
		return models.BenchmarksStatusRunning
	case domain.BenchmarkStatusFinished:
		return models.BenchmarksStatusFinished
	default:
		panic("unknown benchmark status")
	}
}

func toDomainBenchmarkStatus(status models.BenchmarksStatus) (domain.BenchmarkStatus, error) {
	switch status {
	case models.BenchmarksStatusWaiting:
		return domain.BenchmarkStatusWaiting, nil
	case models.BenchmarksStatusRunning:
		return domain.BenchmarkStatusRunning, nil
	case models.BenchmarksStatusFinished:
		return domain.BenchmarkStatusFinished, nil
	default:
		return "", errors.New("unknown benchmark status")
	}
}

func toDomainBenchmark(benchmark *models.Benchmark) (domain.Benchmark, error) {
	id, err := uuid.Parse(benchmark.ID)
	if err != nil {
		return domain.Benchmark{}, err
	}
	instance, err := toDomainInstance(benchmark.R.Instance)
	if err != nil {
		return domain.Benchmark{}, err
	}
	teamID, err := uuid.Parse(benchmark.TeamID)
	if err != nil {
		return domain.Benchmark{}, err
	}
	userID, err := uuid.Parse(benchmark.UserID)
	if err != nil {
		return domain.Benchmark{}, err
	}
	result, err := toDomainBenchmarkResult(benchmark.Result.Ptr())
	if err != nil {
		return domain.Benchmark{}, err
	}
	status, err := toDomainBenchmarkStatus(benchmark.Status)
	if err != nil {
		return domain.Benchmark{}, err
	}

	return domain.Benchmark{
		ID:         id,
		Instance:   instance,
		TeamID:     teamID,
		UserID:     userID,
		Status:     status,
		CreatedAt:  benchmark.CreatedAt,
		StartedAt:  benchmark.StartedAt.Ptr(),
		FinishedAt: benchmark.FinishedAt.Ptr(),
		Score:      benchmark.Score,
		Result:     result,
	}, nil
}

func toDomainBenchmarkLog(log *models.BenchmarkLog) (domain.BenchmarkLog, error) {
	return domain.BenchmarkLog{
		UserLog:  log.UserLog,
		AdminLog: log.AdminLog,
	}, nil
}

func toDomainBenchmarkResult(result *models.BenchmarksResult) (*domain.BenchmarkResult, error) {
	if result == nil {
		return nil, nil
	}
	switch *result {
	case models.BenchmarksResultPassed:
		return ptr.Of(domain.BenchmarkResultStatusPassed), nil
	case models.BenchmarksResultFailed:
		return ptr.Of(domain.BenchmarkResultStatusFailed), nil
	case models.BenchmarksResultError:
		return ptr.Of(domain.BenchmarkResultStatusError), nil
	default:
		return nil, errors.New("unknown benchmark result")
	}
}

func fromDomainBenchmarkResult(result *domain.BenchmarkResult) *models.BenchmarksResult {
	if result == nil {
		return nil
	}
	switch *result {
	case domain.BenchmarkResultStatusPassed:
		return ptr.Of(models.BenchmarksResultPassed)
	case domain.BenchmarkResultStatusFailed:
		return ptr.Of(models.BenchmarksResultFailed)
	case domain.BenchmarkResultStatusError:
		return ptr.Of(models.BenchmarksResultError)
	default:
		panic("unknown benchmark result")
	}
}
