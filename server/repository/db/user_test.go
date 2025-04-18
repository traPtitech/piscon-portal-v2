package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
)

func TestFindUser(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	team := domain.NewTeam("team1")
	user := domain.User{
		ID:     uuid.New(),
		Name:   "user1",
		TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
	}

	mustMakeTeam(t, db, team)
	mustMakeUser(t, db, user)

	got, err := repo.FindUser(t.Context(), user.ID)
	assert.NoError(t, err)

	testutil.CompareUser(t, user, got)
}

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

	got, err := repo.GetUsers(t.Context())
	assert.NoError(t, err)

	testutil.CompareUsers(t, users, got)
}
