package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type InstanceRepository interface {
	// FindInstance finds an instance by id. If the instance is not found, it returns [ErrNotFound].
	FindInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error)
}
