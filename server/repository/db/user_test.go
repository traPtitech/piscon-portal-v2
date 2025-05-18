package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
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

func TestGetUsersByIDs(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)
	team := domain.NewTeam("team1")
	userID1, userID2, userID3 := uuid.New(), uuid.New(), uuid.New()
	user1 := domain.User{
		ID:     userID1,
		Name:   "user1",
		TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
	}
	user2 := domain.User{
		ID:     userID2,
		Name:   "user2",
		TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
	}
	user3 := domain.User{
		ID:     userID3,
		Name:   "user3",
		TeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
	}
	users := []domain.User{user1, user2, user3}

	mustMakeTeam(t, db, team)
	for _, user := range users {
		mustMakeUser(t, db, user)
	}

	testCases := map[string]struct {
		ids   []uuid.UUID
		users []domain.User
	}{
		"idsが空": {},
		"1つ取得できる": {
			ids:   []uuid.UUID{userID1},
			users: []domain.User{user1},
		},
		"2つ取得できる": {
			ids:   []uuid.UUID{userID1, userID2},
			users: []domain.User{user1, user2},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			users, err := repo.GetUsersByIDs(t.Context(), testCase.ids)

			expectedUsersMap := make(map[uuid.UUID]domain.User, len(testCase.users))
			for _, user := range testCase.users {
				expectedUsersMap[user.ID] = user
			}

			assert.NoError(t, err)
			assert.Len(t, users, len(testCase.users))
			for _, user := range users {
				expectedUser, ok := expectedUsersMap[user.ID]
				assert.True(t, ok)
				assert.Equal(t, expectedUser.ID, user.ID)
				assert.Equal(t, expectedUser.Name, user.Name)
				assert.Equal(t, expectedUser.TeamID, user.TeamID)
				assert.Equal(t, expectedUser.IsAdmin, user.IsAdmin)
			}
		})
	}
}

func TestGetAdmins(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	adminUser := domain.User{
		ID:      uuid.New(),
		Name:    "adminUser",
		IsAdmin: true,
	}
	normalUser := domain.User{
		ID:   uuid.New(),
		Name: "normalUser",
	}
	mustMakeUser(t, db, adminUser)
	mustMakeUser(t, db, normalUser)

	got, err := repo.GetAdmins(t.Context())
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, adminUser.ID, got[0].ID)
	assert.Equal(t, adminUser.Name, got[0].Name)
	assert.Equal(t, adminUser.IsAdmin, got[0].IsAdmin)
}

func TestAddAdmins(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	adminUser := domain.User{
		ID:      uuid.New(),
		Name:    "adminUser",
		IsAdmin: true,
	}
	normalUser := domain.User{
		ID:   uuid.New(),
		Name: "normalUser",
	}
	mustMakeUser(t, db, adminUser)
	mustMakeUser(t, db, normalUser)

	err := repo.AddAdmins(t.Context(), []uuid.UUID{adminUser.ID, normalUser.ID})
	assert.NoError(t, err)

	resultAdmins, err := models.Users.Query(
		sm.Where(models.UserColumns.ID.In(mysql.Arg(adminUser.ID, normalUser.ID))),
	).All(t.Context(), db)
	assert.NoError(t, err)
	assert.Len(t, resultAdmins, 2)
	for _, admin := range resultAdmins {
		assert.True(t, admin.IsAdmin)
	}
}

func TestDeleteAdmins(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	adminUser := domain.User{
		ID:      uuid.New(),
		Name:    "adminUser",
		IsAdmin: true,
	}
	normalUser := domain.User{
		ID:   uuid.New(),
		Name: "normalUser",
	}
	mustMakeUser(t, db, adminUser)
	mustMakeUser(t, db, normalUser)

	err := repo.DeleteAdmins(t.Context(), []uuid.UUID{adminUser.ID, normalUser.ID})
	assert.NoError(t, err)

	resultAdmins, err := models.Users.Query(
		sm.Where(models.UserColumns.ID.In(mysql.Arg(adminUser.ID, normalUser.ID))),
	).All(t.Context(), db)
	assert.NoError(t, err)
	assert.Len(t, resultAdmins, 2)
	for _, admin := range resultAdmins {
		assert.False(t, admin.IsAdmin)
	}
}
