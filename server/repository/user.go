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
	// GetUsersByIDs returns users by ids.
	// If the user is not found, it returns an empty slice.
	// If the ids is empty, it returns an empty slice.
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.User, error)
	// GetAdmins returns all admin users.
	GetAdmins(ctx context.Context) ([]domain.User, error)
	// AddAdmins updates admin users.
	// If userIDs contains a id that does not exist, it is ignored.
	// If userIDs is empty, it does nothing.
	// If userIDs contains a id that is already an admin, it is ignored.
	AddAdmins(ctx context.Context, userIDs []uuid.UUID) error
	// DeleteAdmins deletes admin users.
	// If userIDs contains a id that does not exist, it is ignored.
	// If userIDs is empty, it does nothing.
	// If userIDs contains a id that is not an admin, it is ignored.
	DeleteAdmins(ctx context.Context, userIDs []uuid.UUID) error
}
