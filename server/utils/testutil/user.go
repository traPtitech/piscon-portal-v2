package testutil

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func CompareUsers(t *testing.T, want, got []domain.User) {
	t.Helper()

	wantCloned := slices.Clone(want)
	gotCloned := slices.Clone(got)

	slices.SortFunc(wantCloned, func(a, b domain.User) int { return cmp.Compare(a.ID.String(), b.ID.String()) })
	slices.SortFunc(gotCloned, func(a, b domain.User) int { return cmp.Compare(a.ID.String(), b.ID.String()) })

	assert.Len(t, got, len(want))
	for i := range wantCloned {
		CompareUser(t, wantCloned[i], gotCloned[i])
	}
}

func CompareUser(t *testing.T, want, got domain.User) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID, "user.ID mismatch")
	assert.Equal(t, want.Name, got.Name, "user.Name mismatch")
	assert.Equal(t, want.TeamID, got.TeamID, "user.TeamID mismatch")
	assert.Equal(t, want.IsAdmin, got.IsAdmin, "user.IsAdmin mismatch")
}
