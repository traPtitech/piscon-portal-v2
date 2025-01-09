package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	repo := setupRepository(t)

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
	team := domain.Team{
		ID:        uuid.New(),
		Name:      "team1",
		Members:   members,
		CreatedAt: time.Now(),
	}
	for _, member := range members {
		if err := repo.CreateUser(context.Background(), member); err != nil {
			require.Nil(t, err)
		}
	}

	err := repo.CreateTeam(context.Background(), team)
	if err != nil {
		require.Nil(t, err)
	}

	got, err := repo.FindTeam(context.Background(), team.ID)
	if err != nil {
		require.Nil(t, err)
	}
	require.Equal(t, team.ID, got.ID)
	for _, member := range got.Members {
		require.Equal(t, team.ID, member.TeamID.UUID)
	}
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	repo := setupRepository(t)

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
	team := domain.Team{
		ID:        uuid.New(),
		Name:      "team1",
		Members:   members,
		CreatedAt: time.Now(),
	}

	for _, member := range members {
		if err := repo.CreateUser(context.Background(), member); err != nil {
			require.Nil(t, err)
		}
	}

	err := repo.CreateTeam(context.Background(), team)
	if err != nil {
		require.Nil(t, err)
	}

	// change the team name and add a new member
	team.Name = "team2"
	team.AddMember(domain.User{
		ID:   uuid.New(),
		Name: "user3",
	})
	err = repo.UpdateTeam(context.Background(), team)
	if err != nil {
		require.Nil(t, err)
	}

	got, err := repo.FindTeam(context.Background(), team.ID)
	if err != nil {
		require.Nil(t, err)
	}
	require.Equal(t, team.ID, got.ID)
	for _, member := range got.Members {
		require.Equal(t, team.ID, member.TeamID.UUID)
	}
}
