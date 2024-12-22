package domain

import (
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

func NewSession(id, userID string, expiresAt time.Time) Session {
	return Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
}
