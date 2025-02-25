package testutil

import (
	"cmp"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func CompareBenchmarks(t *testing.T, want, got []domain.Benchmark) {
	t.Helper()

	wantCloned := slices.Clone(want)
	gotCloned := slices.Clone(got)

	slices.SortFunc(wantCloned, func(a, b domain.Benchmark) int { return cmp.Compare(a.ID.String(), b.ID.String()) })
	slices.SortFunc(gotCloned, func(a, b domain.Benchmark) int { return cmp.Compare(a.ID.String(), b.ID.String()) })

	assert.Len(t, got, len(want), "benchmark length mismatch")
	for i := range want {
		CompareBenchmark(t, wantCloned[i], gotCloned[i])
	}
}

func CompareBenchmark(t *testing.T, want, got domain.Benchmark) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID, "benchmark.ID mismatch")
	assert.Equal(t, want.TeamID, got.TeamID, "benchmark.TeamID mismatch")
	assert.Equal(t, want.UserID, got.UserID, "benchmark.UserID mismatch")
	assert.Equal(t, want.Status, got.Status, "benchmark.Status mismatch")
	assert.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second, "benchmark.CreatedAt mismatch")
	if want.StartedAt == nil {
		assert.Nil(t, got.StartedAt, "benchmark.StartedAt mismatch")
	} else {
		assert.NotNil(t, got.StartedAt, "benchmark.StartedAt mismatch")
		assert.WithinDuration(t, *want.StartedAt, *got.StartedAt, time.Second, "benchmark.StartedAt mismatch")
	}
	if want.FinishedAt == nil {
		assert.Nil(t, got.FinishedAt, "benchmark.FinishedAt mismatch")
	} else {
		assert.NotNil(t, got.FinishedAt, "benchmark.FinishedAt mismatch")
		assert.WithinDuration(t, *want.FinishedAt, *got.FinishedAt, time.Second, "benchmark.EndedAt mismatch")
	}
	assert.Equal(t, want.Score, got.Score, "benchmark.Score mismatch")
	assert.Equal(t, want.Result, got.Result, "benchmark.Result mismatch")
}

func CompareBenchmarkLog(t *testing.T, want, got domain.BenchmarkLog) {
	t.Helper()

	assert.Equal(t, want.UserLog, got.UserLog, "benchmarklog.UserLog mismatch")
	assert.Equal(t, want.AdminLog, got.AdminLog, "benchmarklog.AdminLog mismatch")
}
