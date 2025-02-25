package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
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
		return Benchmark{}, errors.New("teamID is not match")
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

func (b *Benchmark) Start() {
	b.Status = BenchmarkStatusRunning
	b.StartedAt = ptr.Of(time.Now())
}

func (b *Benchmark) Finish(score int64, result BenchmarkResult) {
	b.Status = BenchmarkStatusFinished
	b.Score = score
	b.FinishedAt = ptr.Of(time.Now())
	b.Result = &result
}

func (b *Benchmark) IsFinished() bool {
	return b.Status == BenchmarkStatusFinished
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

type Benchmarks []Benchmark

func (b Benchmarks) Filter(f func(b Benchmark) bool) Benchmarks {
	var res Benchmarks
	for _, b := range b {
		if f(b) {
			res = append(res, b)
		}
	}
	return res
}
