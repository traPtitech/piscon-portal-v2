package db_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func mustMakeUser(t *testing.T, executor bob.Executor, user domain.User) {
	t.Helper()
	_, err := models.Users.Insert(&models.UserSetter{
		ID:      lo.ToPtr(user.ID.String()),
		Name:    lo.ToPtr(user.Name),
		TeamID:  lo.Ternary(user.TeamID.Valid, db.ToSQLNull(user.TeamID.UUID.String()), nil),
		IsAdmin: lo.ToPtr(user.IsAdmin),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeTeam(t *testing.T, executor bob.Executor, team domain.Team) {
	t.Helper()
	_, err := models.Teams.Insert(&models.TeamSetter{
		ID:        lo.ToPtr(team.ID.String()),
		Name:      lo.ToPtr(team.Name),
		CreatedAt: lo.ToPtr(team.CreatedAt),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeInstance(t *testing.T, executor bob.Executor, instance domain.Instance) {
	t.Helper()
	status, err := db.FromDomainInstanceStatus(instance.Infra.Status)
	require.NoError(t, err)
	_, err = models.Instances.Insert(&models.InstanceSetter{
		ID:                 lo.ToPtr(instance.ID.String()),
		ProviderInstanceID: lo.ToPtr(instance.Infra.ProviderInstanceID),
		TeamID:             lo.ToPtr(instance.TeamID.String()),
		InstanceNumber:     lo.ToPtr(int32(instance.Index)),
		Status:             lo.ToPtr(status),
		PrivateIP:          db.ToSQLNull(instance.Infra.PrivateIP),
		PublicIP:           db.ToSQLNull(instance.Infra.PublicIP),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeBenchmark(t *testing.T, executor bob.Executor, benchmark domain.Benchmark) {
	t.Helper()
	status, err := db.FromDomainBenchmarkStatus(benchmark.Status)
	require.NoError(t, err)
	result, err := db.FromDomainBenchmarkResult(benchmark.Result)
	require.NoError(t, err)
	_, err = models.Benchmarks.Insert(&models.BenchmarkSetter{
		ID:         lo.ToPtr(benchmark.ID.String()),
		InstanceID: lo.ToPtr(benchmark.Instance.ID.String()),
		TeamID:     lo.ToPtr(benchmark.TeamID.String()),
		UserID:     lo.ToPtr(benchmark.UserID.String()),
		Status:     lo.ToPtr(status),
		CreatedAt:  lo.ToPtr(benchmark.CreatedAt),
		StartedAt:  db.PtrToSQLNull(benchmark.StartedAt),
		FinishedAt: db.PtrToSQLNull(benchmark.FinishedAt),
		Score:      lo.ToPtr(benchmark.Score),
		Result:     db.PtrToSQLNull(result),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeBenchmarkLog(t *testing.T, executor bob.Executor, benchmarkID uuid.UUID, log domain.BenchmarkLog) {
	t.Helper()
	_, err := models.BenchmarkLogs.Insert(&models.BenchmarkLogSetter{
		BenchmarkID: lo.ToPtr(benchmarkID.String()),
		UserLog:     lo.ToPtr(log.UserLog),
		AdminLog:    lo.ToPtr(log.AdminLog),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeDocument(t *testing.T, executor bob.Executor, document domain.Document) {
	t.Helper()
	_, err := models.Documents.Insert(&models.DocumentSetter{
		ID:        lo.ToPtr(document.ID.String()),
		Body:      lo.ToPtr(document.Body),
		CreatedAt: lo.ToPtr(document.CreatedAt),
	}).Exec(t.Context(), executor)
	require.NoError(t, err)
}
