package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql/dialect"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Benchmark{}, repository.ErrNotFound
		}
		return domain.Benchmark{}, fmt.Errorf("find benchmark: %w", err)
	}

	return toDomainBenchmark(benchmark)
}

func (r *Repository) CreateBenchmark(ctx context.Context, benchmark domain.Benchmark) error {
	status, err := fromDomainBenchmarkStatus(benchmark.Status)
	if err != nil {
		return err
	}
	_, err = models.Benchmarks.Insert(&models.BenchmarkSetter{
		ID:         omit.From(benchmark.ID.String()),
		InstanceID: omit.From(benchmark.Instance.ID.String()),
		TeamID:     omit.From(benchmark.TeamID.String()),
		UserID:     omit.From(benchmark.UserID.String()),
		Status:     omit.From(status),
		CreatedAt:  omit.From(benchmark.CreatedAt),
	}).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("create benchmark: %w", err)
	}
	return nil
}

func (r *Repository) GetBenchmarks(ctx context.Context, query repository.BenchmarkQuery) ([]domain.Benchmark, error) {
	where := models.SelectWhere.Benchmarks

	mods := bob.Mods[*dialect.SelectQuery]{
		models.PreloadBenchmarkInstance(),
		sm.OrderBy(models.BenchmarkColumns.CreatedAt).Asc(),
	}
	if query.TeamID.IsSet() {
		teamID := query.TeamID.Get().String()
		mods = append(mods, where.TeamID.EQ(teamID))
	}
	if query.StatusIn.IsSet() {
		var statuses []models.BenchmarksStatus
		for _, status := range query.StatusIn.Get() {
			dbModelStatus, err := fromDomainBenchmarkStatus(status)
			if err != nil {
				return nil, err
			}
			statuses = append(statuses, dbModelStatus)
		}
		mods = append(mods, where.Status.In(statuses...))
	}

	benchmarks, err := models.Benchmarks.Query(mods...).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get benchmarks: %w", err)
	}

	res := make([]domain.Benchmark, 0, len(benchmarks))
	for _, b := range benchmarks {
		benchmark, err := toDomainBenchmark(b)
		if err != nil {
			return nil, err
		}
		res = append(res, benchmark)
	}

	return res, nil
}

func (r *Repository) GetOldestQueuedBenchmark(ctx context.Context) (domain.Benchmark, error) {
	statusWaiting := models.SelectWhere.Benchmarks.Status.EQ(models.BenchmarksStatusWaiting)
	orderByCreatedAtAsc := sm.OrderBy(models.BenchmarkColumns.CreatedAt).Asc()
	limit1 := sm.Limit(1)
	benchmark, err := models.Benchmarks.Query(
		models.PreloadBenchmarkInstance(),
		statusWaiting, orderByCreatedAtAsc, limit1).One(ctx, r.executor(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Benchmark{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("get benchmark: %w", err)
	}

	return toDomainBenchmark(benchmark)
}

func (r *Repository) UpdateBenchmark(ctx context.Context, id uuid.UUID, benchmark domain.Benchmark) error {
	status, err := fromDomainBenchmarkStatus(benchmark.Status)
	if err != nil {
		return err
	}
	result, err := fromDomainBenchmarkResult(benchmark.Result)
	if err != nil {
		return err
	}

	whereID := models.UpdateWhere.Benchmarks.ID.EQ(id.String())
	newBenchmark := models.BenchmarkSetter{
		ID:         omit.From(id.String()),
		InstanceID: omit.From(benchmark.Instance.ID.String()),
		TeamID:     omit.From(benchmark.TeamID.String()),
		UserID:     omit.From(benchmark.UserID.String()),
		Status:     omit.From(status),
		CreatedAt:  omit.From(benchmark.CreatedAt),
		StartedAt:  omitnull.FromPtr(benchmark.StartedAt),
		FinishedAt: omitnull.FromPtr(benchmark.FinishedAt),
		Score:      omit.From(benchmark.Score),
		Result:     omitnull.FromPtr(result),
	}

	_, err = models.Benchmarks.Update(whereID, newBenchmark.UpdateMod()).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("update benchmark: %w", err)
	}

	return nil
}

func (r *Repository) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	benchmarkLogs, err := models.FindBenchmarkLog(ctx, r.executor(ctx), benchmarkID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BenchmarkLog{}, repository.ErrNotFound
		}
		return domain.BenchmarkLog{}, fmt.Errorf("get benchmark log: %w", err)
	}

	return toDomainBenchmarkLog(benchmarkLogs)
}

func fromDomainBenchmarkStatus(status domain.BenchmarkStatus) (models.BenchmarksStatus, error) {
	switch status {
	case domain.BenchmarkStatusWaiting:
		return models.BenchmarksStatusWaiting, nil
	case domain.BenchmarkStatusRunning:
		return models.BenchmarksStatusRunning, nil
	case domain.BenchmarkStatusFinished:
		return models.BenchmarksStatusFinished, nil
	default:
		return "", errors.New("unknown benchmark status")
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
		return domain.Benchmark{}, fmt.Errorf("parse benchmark id: %w", err)
	}
	instance, err := toDomainInstance(benchmark.R.Instance)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("parse benchmark instance: %w", err)
	}
	teamID, err := uuid.Parse(benchmark.TeamID)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("parse benchmark team id: %w", err)
	}
	userID, err := uuid.Parse(benchmark.UserID)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("parse benchmark user id: %w", err)
	}
	result, err := toDomainBenchmarkResult(benchmark.Result.Ptr())
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("parse benchmark result: %w", err)
	}
	status, err := toDomainBenchmarkStatus(benchmark.Status)
	if err != nil {
		return domain.Benchmark{}, fmt.Errorf("parse benchmark status: %w", err)
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
		return nil, fmt.Errorf("unknown benchmark result: %v", *result)
	}
}

func fromDomainBenchmarkResult(result *domain.BenchmarkResult) (*models.BenchmarksResult, error) {
	if result == nil {
		return nil, nil
	}
	switch *result {
	case domain.BenchmarkResultStatusPassed:
		return ptr.Of(models.BenchmarksResultPassed), nil
	case domain.BenchmarkResultStatusFailed:
		return ptr.Of(models.BenchmarksResultFailed), nil
	case domain.BenchmarkResultStatusError:
		return ptr.Of(models.BenchmarksResultError), nil
	default:
		return nil, fmt.Errorf("unknown benchmark result: %v", *result)
	}
}
