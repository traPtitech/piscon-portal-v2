package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type UserRepository interface {
	// FindUser finds a user by id. If the user is not found, it returns [ErrNotFound].
	FindUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	// GetUsers returns all users.
	GetUsers(ctx context.Context) ([]domain.User, error)
	// CreateUser creates a user.
	CreateUser(ctx context.Context, user domain.User) error
}
