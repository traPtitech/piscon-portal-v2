package db_test

import (
	"context"
	"testing"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
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
		ID:     omit.From(user.ID.String()),
		Name:   omit.From(user.Name),
		TeamID: lo.Ternary(user.TeamID.Valid, omitnull.From(user.TeamID.UUID.String()), omitnull.Val[string]{}),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeTeam(t *testing.T, executor bob.Executor, team domain.Team) {
	t.Helper()
	_, err := models.Teams.Insert(&models.TeamSetter{
		ID:        omit.From(team.ID.String()),
		Name:      omit.From(team.Name),
		CreatedAt: omit.From(team.CreatedAt),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeInstance(t *testing.T, executor bob.Executor, instance domain.Instance) {
	t.Helper()
	status, err := db.FromDomainInstanceStatus(instance.Status)
	require.NoError(t, err)
	_, err = models.Instances.Insert(&models.InstanceSetter{
		ID:             omit.From(instance.ID.String()),
		TeamID:         omit.From(instance.TeamID.String()),
		InstanceNumber: omit.From(int32(instance.InstanceNumber)),
		Status:         omit.From(status),
		PrivateIP:      omitnull.From(instance.PrivateIP),
		PublicIP:       omitnull.From(instance.PublicIP),
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
		ID:         omit.From(benchmark.ID.String()),
		InstanceID: omit.From(benchmark.Instance.ID.String()),
		TeamID:     omit.From(benchmark.TeamID.String()),
		UserID:     omit.From(benchmark.UserID.String()),
		Status:     omit.From(status),
		CreatedAt:  omit.From(benchmark.CreatedAt),
		StartedAt:  omitnull.FromPtr(benchmark.StartedAt),
		FinishedAt: omitnull.FromPtr(benchmark.FinishedAt),
		Score:      omit.From(benchmark.Score),
		Result:     omitnull.FromPtr(result),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}

func mustMakeBenchmarkLog(t *testing.T, executor bob.Executor, benchmarkID uuid.UUID, log domain.BenchmarkLog) {
	t.Helper()
	_, err := models.BenchmarkLogs.Insert(&models.BenchmarkLogSetter{
		BenchmarkID: omit.From(benchmarkID.String()),
		UserLog:     omit.From(log.UserLog),
		AdminLog:    omit.From(log.AdminLog),
	}).Exec(context.Background(), executor)
	require.NoError(t, err)
}
