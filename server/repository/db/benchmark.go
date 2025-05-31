package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/dialect"
	"github.com/stephenafamo/bob/dialect/mysql/fm"
	"github.com/stephenafamo/bob/dialect/mysql/im"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/wm"
	"github.com/stephenafamo/scan"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
)

func (r *Repository) FindBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error) {
	benchmark, err := models.Benchmarks.
		Query(
			models.Preload.Benchmark.Instance(),
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
		ID:         lo.ToPtr(benchmark.ID.String()),
		InstanceID: lo.ToPtr(benchmark.Instance.ID.String()),
		TeamID:     lo.ToPtr(benchmark.TeamID.String()),
		UserID:     lo.ToPtr(benchmark.UserID.String()),
		Status:     lo.ToPtr(status),
		CreatedAt:  lo.ToPtr(benchmark.CreatedAt),
	}).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("create benchmark: %w", err)
	}
	return nil
}

func (r *Repository) GetBenchmarks(ctx context.Context, query repository.BenchmarkQuery) ([]domain.Benchmark, error) {
	where := models.SelectWhere.Benchmarks

	mods := bob.Mods[*dialect.SelectQuery]{
		models.Preload.Benchmark.Instance(),
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
		models.Preload.Benchmark.Instance(),
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
		ID:           lo.ToPtr(id.String()),
		InstanceID:   lo.ToPtr(benchmark.Instance.ID.String()),
		TeamID:       lo.ToPtr(benchmark.TeamID.String()),
		UserID:       lo.ToPtr(benchmark.UserID.String()),
		Status:       lo.ToPtr(status),
		CreatedAt:    lo.ToPtr(benchmark.CreatedAt),
		StartedAt:    PtrToSQLNull(benchmark.StartedAt),
		FinishedAt:   PtrToSQLNull(benchmark.FinishedAt),
		Score:        lo.ToPtr(benchmark.Score),
		Result:       PtrToSQLNull(result),
		ErrorMessage: PtrToSQLNull(benchmark.ErrorMes),
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

func (r *Repository) UpdateBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID, log domain.BenchmarkLog) error {
	_, err := models.BenchmarkLogs.Insert(
		&models.BenchmarkLogSetter{
			BenchmarkID: lo.ToPtr(benchmarkID.String()),
			UserLog:     lo.ToPtr(log.UserLog),
			AdminLog:    lo.ToPtr(log.AdminLog),
		},
		im.OnDuplicateKeyUpdate(
			im.UpdateWithValues(models.ColumnNames.BenchmarkLogs.UserLog, models.ColumnNames.BenchmarkLogs.AdminLog),
		),
	).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("update benchmark log: %w", err)
	}

	return nil
}

func (r *Repository) GetRanking(ctx context.Context, query repository.RankingQuery) (ranking []domain.Score, err error) {
	var maxColumn string
	switch query.OrderBy {
	case domain.RankingOrderByLatestScore:
		maxColumn = "created_at"
	case domain.RankingOrderByHighestScore:
		maxColumn = "score"
	default:
		return nil, fmt.Errorf("unknown ranking order by: %v", query.OrderBy)
	}

	q := mysql.Select(
		sm.From(
			models.Benchmarks.Query(
				sm.Columns(
					models.Benchmarks.Columns(),
					mysql.F("RANK")(
						fm.Over(
							wm.PartitionBy(models.BenchmarkColumns.TeamID),
							wm.OrderBy(maxColumn).Desc()),
					).As("rank_in_team"),
				),
				models.SelectWhere.Benchmarks.Status.EQ(models.BenchmarksStatusFinished),
			),
		).As("rank_team"),
		sm.Where(mysql.Quote("rank_team", "rank_in_team").EQ(mysql.Arg(1))),
		sm.OrderBy(mysql.Quote("rank_team", "score")).Desc(),
		sm.OrderBy(mysql.Quote("rank_team", "created_at")).Asc(),
	)

	type BenchmarkWithRank struct {
		models.Benchmark
		TeamRank int64 `db:"rank_in_team"`
	}

	benchmarks, err := bob.All(ctx, r.executor(ctx), q, scan.StructMapper[BenchmarkWithRank]())
	if err != nil {
		return nil, fmt.Errorf("scan ranking: %w", err)
	}

	ranking = make([]domain.Score, 0, len(benchmarks))
	for _, benchmark := range benchmarks {
		domainScore, err := toDomainScore(benchmark.Benchmark)
		if err != nil {
			return nil, fmt.Errorf("to domain score: %w", err)
		}
		ranking = append(ranking, domainScore)
	}

	return ranking, nil
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
	result, err := toDomainBenchmarkResult(benchmark.Result)
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
		StartedAt:  SQLNullToPtr(benchmark.StartedAt),
		FinishedAt: SQLNullToPtr(benchmark.FinishedAt),
		Score:      benchmark.Score,
		Result:     result,
		ErrorMes:   SQLNullToPtr(benchmark.ErrorMessage),
	}, nil
}

func toDomainBenchmarkLog(log *models.BenchmarkLog) (domain.BenchmarkLog, error) {
	return domain.BenchmarkLog{
		UserLog:  log.UserLog,
		AdminLog: log.AdminLog,
	}, nil
}

func toDomainBenchmarkResult(result sql.Null[models.BenchmarksResult]) (*domain.BenchmarkResult, error) {
	if !result.Valid {
		return nil, nil
	}
	switch result.V {
	case models.BenchmarksResultPassed:
		return ptr.Of(domain.BenchmarkResultStatusPassed), nil
	case models.BenchmarksResultFailed:
		return ptr.Of(domain.BenchmarkResultStatusFailed), nil
	case models.BenchmarksResultError:
		return ptr.Of(domain.BenchmarkResultStatusError), nil
	default:
		return nil, fmt.Errorf("unknown benchmark result: %v", result.V)
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

func toDomainScore(benchmark models.Benchmark) (domain.Score, error) {
	id, err := uuid.Parse(benchmark.ID)
	if err != nil {
		return domain.Score{}, fmt.Errorf("parse benchmark id: %w", err)
	}
	teamID, err := uuid.Parse(benchmark.TeamID)
	if err != nil {
		return domain.Score{}, fmt.Errorf("parse benchmark team id: %w", err)
	}

	return domain.Score{
		BenchmarkID: id,
		TeamID:      teamID,
		Score:       benchmark.Score,
		CreatedAt:   benchmark.CreatedAt,
	}, nil
}
