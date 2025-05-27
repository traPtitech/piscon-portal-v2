package domain

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID
	Body      string
	CreatedAt time.Time
}
