package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type Repository interface {
	// Transaction starts a transaction and calls f with the transaction.
	// If f returns an error, the transaction is rolled back and the error is returned.
	Transaction(ctx context.Context, f func(ctx context.Context, r Repository) error) error

	// FindUser finds a user by id. If the user is not found, it returns [ErrNotFound].
	FindUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	// GetUsers returns all users.
	GetUsers(ctx context.Context) ([]domain.User, error)
	// CreateUser creates a user.
	CreateUser(ctx context.Context, user domain.User) error

	// FindSession finds a session by id. If the session is not found, it returns [ErrNotFound].
	FindSession(ctx context.Context, id string) (domain.Session, error)
	// CreateSession creates a session.
	CreateSession(ctx context.Context, session domain.Session) error
	// DeleteSession deletes a session by id.
	DeleteSession(ctx context.Context, id string) error

	// FindTeam finds a team by id. If the team is not found, it returns [ErrNotFound].
	FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error)
	// GetTeams returns all teams.
	GetTeams(ctx context.Context) ([]domain.Team, error)
	// CreateTeam creates a team.
	CreateTeam(ctx context.Context, team domain.Team) error
	// UpdateTeam updates a team.
	UpdateTeam(ctx context.Context, team domain.Team) error
}

var (
	ErrNotFound = errors.New("not found")
)
