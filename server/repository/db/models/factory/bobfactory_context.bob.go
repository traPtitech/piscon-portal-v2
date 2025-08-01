// Code generated by BobGen mysql v0.38.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"context"

	models "github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

type contextKey string

var (
	// Table context

	benchmarkLogCtx = newContextual[*models.BenchmarkLog]("benchmarkLog")
	benchmarkCtx    = newContextual[*models.Benchmark]("benchmark")
	documentCtx     = newContextual[*models.Document]("document")
	instanceCtx     = newContextual[*models.Instance]("instance")
	sessionCtx      = newContextual[*models.Session]("session")
	teamCtx         = newContextual[*models.Team]("team")
	userCtx         = newContextual[*models.User]("user")

	// Relationship Contexts for benchmark_logs
	benchmarkLogWithParentsCascadingCtx = newContextual[bool]("benchmarkLogWithParentsCascading")
	benchmarkLogRelBenchmarkCtx         = newContextual[bool]("benchmark_logs.benchmarks.benchmark_logs_ibfk_1")

	// Relationship Contexts for benchmarks
	benchmarkWithParentsCascadingCtx = newContextual[bool]("benchmarkWithParentsCascading")
	benchmarkRelBenchmarkLogCtx      = newContextual[bool]("benchmark_logs.benchmarks.benchmark_logs_ibfk_1")
	benchmarkRelInstanceCtx          = newContextual[bool]("benchmarks.instances.benchmarks_ibfk_1")

	// Relationship Contexts for documents
	documentWithParentsCascadingCtx = newContextual[bool]("documentWithParentsCascading")

	// Relationship Contexts for instances
	instanceWithParentsCascadingCtx = newContextual[bool]("instanceWithParentsCascading")
	instanceRelBenchmarksCtx        = newContextual[bool]("benchmarks.instances.benchmarks_ibfk_1")

	// Relationship Contexts for sessions
	sessionWithParentsCascadingCtx = newContextual[bool]("sessionWithParentsCascading")
	sessionRelUserCtx              = newContextual[bool]("sessions.users.sessions_ibfk_1")

	// Relationship Contexts for teams
	teamWithParentsCascadingCtx = newContextual[bool]("teamWithParentsCascading")
	teamRelUsersCtx             = newContextual[bool]("teams.users.users_ibfk_1")

	// Relationship Contexts for users
	userWithParentsCascadingCtx = newContextual[bool]("userWithParentsCascading")
	userRelSessionsCtx          = newContextual[bool]("sessions.users.sessions_ibfk_1")
	userRelTeamCtx              = newContextual[bool]("teams.users.users_ibfk_1")
)

// Contextual is a convienience wrapper around context.WithValue and context.Value
type contextual[V any] struct {
	key contextKey
}

func newContextual[V any](key string) contextual[V] {
	return contextual[V]{key: contextKey(key)}
}

func (k contextual[V]) WithValue(ctx context.Context, val V) context.Context {
	return context.WithValue(ctx, k.key, val)
}

func (k contextual[V]) Value(ctx context.Context) (V, bool) {
	v, ok := ctx.Value(k.key).(V)
	return v, ok
}
