package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
)

func TestGetTeam(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	teamID := uuid.New()
	members := []domain.User{
		{
			ID:     uuid.New(),
			Name:   "user1",
			TeamID: uuid.NullUUID{UUID: teamID, Valid: true},
		},
		{
			ID:     uuid.New(),
			Name:   "user2",
			TeamID: uuid.NullUUID{UUID: teamID, Valid: true},
		},
	}
	team := domain.Team{
		ID:        teamID,
		Name:      "team1",
		Members:   members,
		CreatedAt: time.Now(),
	}
	mustMakeTeam(t, db, team)
	for _, member := range members {
		mustMakeUser(t, db, member)
	}

	got, err := repo.FindTeam(context.Background(), team.ID)
	require.NoError(t, err)

	testutil.CompareTeam(t, team, got)
}

func TestGetTeams(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	teamID1 := uuid.New()
	teamID2 := uuid.New()

	members1 := []domain.User{
		{
			ID:     uuid.New(),
			Name:   "user1",
			TeamID: uuid.NullUUID{UUID: teamID1, Valid: true},
		},
		{
			ID:     uuid.New(),
			Name:   "user2",
			TeamID: uuid.NullUUID{UUID: teamID1, Valid: true},
		},
	}
	members2 := []domain.User{
		{
			ID:     uuid.New(),
			Name:   "user3",
			TeamID: uuid.NullUUID{UUID: teamID2, Valid: true},
		},
		{
			ID:     uuid.New(),
			Name:   "user4",
			TeamID: uuid.NullUUID{UUID: teamID2, Valid: true},
		},
	}
	teams := []domain.Team{
		{
			ID:        teamID1,
			Name:      "team1",
			Members:   members1,
			CreatedAt: time.Now(),
		},
		{
			ID:        teamID2,
			Name:      "team2",
			Members:   members2,
			CreatedAt: time.Now(),
		},
	}
	for _, team := range teams {
		mustMakeTeam(t, db, team)
		for _, member := range team.Members {
			mustMakeUser(t, db, member)
		}
	}

	got, err := repo.GetTeams(context.Background())
	require.NoError(t, err)

	testutil.CompareTeams(t, teams, got)

}

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	members := []domain.User{
		{
			ID:   uuid.New(),
			Name: "user1",
		},
		{
			ID:   uuid.New(),
			Name: "user2",
		},
	}
	for _, member := range members {
		mustMakeUser(t, db, member)
	}

	team := domain.Team{
		ID:        uuid.New(),
		Name:      "team1",
		Members:   members,
		CreatedAt: time.Now(),
	}
	for i := range team.Members {
		team.Members[i].TeamID = uuid.NullUUID{UUID: team.ID, Valid: true}
	}

	err := repo.CreateTeam(context.Background(), team)
	if !assert.NoError(t, err) {
		return
	}

	got, err := repo.FindTeam(context.Background(), team.ID)
	require.NoError(t, err)

	testutil.CompareTeam(t, team, got)
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	team := domain.Team{
		ID:        uuid.New(),
		Name:      "team1",
		Members:   nil,
		CreatedAt: time.Now(),
	}
	newMember := domain.User{
		ID:   uuid.New(),
		Name: "user2",
	}
	mustMakeUser(t, db, newMember)
	mustMakeTeam(t, db, team)

	// change the team name and add a new member
	team.Name = "team2"
	require.NoError(t, team.AddMember(newMember))

	err := repo.UpdateTeam(context.Background(), team)
	assert.NoError(t, err)

	got, err := repo.FindTeam(context.Background(), team.ID)
	require.NoError(t, err)

	testutil.CompareTeam(t, team, got)
}
