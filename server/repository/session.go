package repository

import (
	"context"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type SessionRepository interface {
	// FindSession finds a session by id. If the session is not found, it returns [ErrNotFound].
	FindSession(ctx context.Context, id string) (domain.Session, error)
	// CreateSession creates a session.
	CreateSession(ctx context.Context, session domain.Session) error
	// DeleteSession deletes a session by id.
	DeleteSession(ctx context.Context, id string) error
}
