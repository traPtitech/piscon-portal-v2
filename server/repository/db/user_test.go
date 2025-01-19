package db_test

import (
	"cmp"
	"context"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	team := domain.NewTeam("team1")
	users := []domain.User{
		{
			ID:     uuid.New(),
			Name:   "user1",
			TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
		},
		{
			ID:     uuid.New(),
			Name:   "user2",
			TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
		},
	}

	mustMakeTeam(t, db, team)
	for _, user := range users {
		mustMakeUser(t, db, user)
	}

	got, err := repo.GetUsers(context.Background())
	require.NoError(t, err)

	compareUsers(t, users, got)
}

func compareUsers(t *testing.T, want, got []domain.User) {
	t.Helper()

	// sort by ID
	slices.SortFunc(want, func(a, b domain.User) int { return cmp.Compare(a.ID.String(), b.ID.String()) })
	slices.SortFunc(got, func(a, b domain.User) int { return cmp.Compare(a.ID.String(), b.ID.String()) })

	require.Len(t, got, len(want))
	for i := range want {
		compareUser(t, want[i], got[i])
	}
}

func compareUser(t *testing.T, want, got domain.User) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.TeamID, got.TeamID)
}
