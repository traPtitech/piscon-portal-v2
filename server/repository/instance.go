package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type InstanceRepository interface {
	CreateInstance(ctx context.Context, instance domain.Instance) error
	// FindInstance finds an instance by id. If the instance is not found, it returns [ErrNotFound].
	FindInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error)
	// GetTeamInstances retrieves all instances for a given team ID.
	// Deleted instances are not included.
	GetTeamInstances(ctx context.Context, teamID uuid.UUID) ([]domain.Instance, error)
	// GetAllInstances retrieves all instances.
	// Deleted instances are not included.
	GetAllInstances(ctx context.Context) ([]domain.Instance, error)
	// DeleteInstance marks an instance as deleted by its ID.
	// If the instance is not found or already deleted, it returns [ErrNotFound].
	DeleteInstance(ctx context.Context, id uuid.UUID) error
}
