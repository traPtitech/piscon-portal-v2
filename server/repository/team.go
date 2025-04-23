package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type TeamRepository interface {
	// FindTeam finds a team by id. If the team is not found, it returns [ErrNotFound].
	FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error)
	// GetTeams returns all teams.
	GetTeams(ctx context.Context) ([]domain.Team, error)
	// CreateTeam creates a team.
	CreateTeam(ctx context.Context, team domain.Team) error
	// UpdateTeam updates a team.
	UpdateTeam(ctx context.Context, team domain.Team) error
}
