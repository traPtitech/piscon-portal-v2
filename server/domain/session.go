package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/utils/random"
)

type Session struct {
	ID        string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func NewSessionID() string {
	return random.String(32)
}

func NewSession(id string, userID uuid.UUID, expiresAt time.Time) Session {
	return Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
}
