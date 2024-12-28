package domain

import (
	"time"

	"github.com/traPtitech/piscon-portal-v2/server/utils/random"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

func NewSessionID() string {
	return random.String(32)
}

func NewSession(id, userID string, expiresAt time.Time) Session {
	return Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
}
