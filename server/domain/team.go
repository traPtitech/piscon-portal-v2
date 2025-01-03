package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        string
	Name      string
	Members   []User
	CreatedAt time.Time
}

func NewTeam(name string) Team {
	return Team{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func (t *Team) AddMember(user User) error {
	if slices.ContainsFunc(t.Members, func(u User) bool { return u.ID == user.ID }) {
		return nil
	}
	if user.TeamID != nil && *user.TeamID != t.ID {
		return errors.New("user is already in another team")
	}
	if len(t.Members) >= 3 {
		return errors.New("team is full")
	}
	t.Members = append(t.Members, user)
	return nil
}
