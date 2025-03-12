package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type BenchmarkStatus string

const (
	BenchmarkStatusWaiting  BenchmarkStatus = "waiting"
	BenchmarkStatusRunning  BenchmarkStatus = "running"
	BenchmarkStatusFinished BenchmarkStatus = "finished"
)

type Benchmark struct {
	ID         uuid.UUID
	Instance   Instance
	TeamID     uuid.UUID
	UserID     uuid.UUID
	Status     BenchmarkStatus
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Score      int64
	Result     *BenchmarkResult
}

func NewBenchmark(instance Instance, user User) (Benchmark, error) {
	if !user.TeamID.Valid {
		return Benchmark{}, errors.New("user is not in a team")
	}
	if instance.Status != InstanceStatusRunning {
		return Benchmark{}, errors.New("instance is not running")
	}
	if instance.TeamID != user.TeamID.UUID {
		return Benchmark{}, errors.New("teamID does not match")
	}

	return Benchmark{
		ID:        uuid.New(),
		Instance:  instance,
		TeamID:    user.TeamID.UUID,
		UserID:    user.ID,
		Status:    BenchmarkStatusWaiting,
		CreatedAt: time.Now(),
	}, nil
}

type BenchmarkResult string

const (
	BenchmarkResultStatusPassed BenchmarkResult = "passed"
	BenchmarkResultStatusFailed BenchmarkResult = "failed"
	BenchmarkResultStatusError  BenchmarkResult = "error"
)

type BenchmarkLog struct {
	UserLog  string
	AdminLog string
}
