package testutil

import (
	"cmp"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func CompareTeams(t *testing.T, want, got []domain.Team) {
	t.Helper()

	wantCloned := slices.Clone(want)
	gotCloned := slices.Clone(got)

	slices.SortFunc(wantCloned, func(a, b domain.Team) int { return cmp.Compare(a.ID.String(), b.ID.String()) })
	slices.SortFunc(gotCloned, func(a, b domain.Team) int { return cmp.Compare(a.ID.String(), b.ID.String()) })

	assert.Len(t, got, len(wantCloned))
	for i := range wantCloned {
		CompareTeam(t, wantCloned[i], gotCloned[i])
	}
}

func CompareTeam(t *testing.T, want, got domain.Team) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID, "team.ID mismatch")
	assert.Equal(t, want.Name, got.Name, "team.Name mismatch")
	assert.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second, "team.CreatedAt mismatch")
	assert.ElementsMatch(t, want.GitHubIDs, got.GitHubIDs, "team.GitHubIDs mismatch")
	CompareUsers(t, want.Members, got.Members)
}
