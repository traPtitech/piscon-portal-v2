package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const MaxTeamMembers = 3

type Team struct {
	ID        uuid.UUID
	Name      string
	Members   []User
	GitHubIDs []string
	CreatedAt time.Time
}

func NewTeam(name string) Team {
	return Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (t *Team) SetMembers(users []User) error {
	if len(users) >= MaxTeamMembers {
		return errors.New("team is full")
	}

	for i, user := range users {
		if user.TeamID.Valid && user.TeamID.UUID != t.ID {
			return errors.New("user is already in another team")
		}
		users[i].TeamID = uuid.NullUUID{UUID: t.ID, Valid: true}
	}

	t.Members = users
	return nil
}
