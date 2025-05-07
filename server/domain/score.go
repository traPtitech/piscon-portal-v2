package domain

import (
	"time"

	"github.com/google/uuid"
)

type Score struct {
	BenchmarkID uuid.UUID
	TeamID      uuid.UUID
	Score       int64
	CreatedAt   time.Time
}
