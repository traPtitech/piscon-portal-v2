package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

const MaxTeamMembers = 3

type Team struct {
	ID        uuid.UUID
	Name      string
	Members   []User
	CreatedAt time.Time
}

func NewTeam(name string) Team {
	return Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (t *Team) AddMember(user User) error {
	if slices.ContainsFunc(t.Members, func(u User) bool { return u.ID == user.ID }) {
		return nil
	}
	if user.TeamID.Valid && user.TeamID.UUID != t.ID {
		return errors.New("user is already in another team")
	}
	if len(t.Members) >= MaxTeamMembers {
		return errors.New("team is full")
	}
	t.Members = append(t.Members, user)
	return nil
}
