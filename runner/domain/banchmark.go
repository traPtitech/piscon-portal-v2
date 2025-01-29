package domain

import "time"

type Job struct {
	id        string
	targetURL string
}

func NewJob(id string, targetURL string) *Job {
	return &Job{
		id:        id,
		targetURL: targetURL,
	}
}

func (j *Job) GetID() string {
	return j.id
}

type Result int

const (
	ResultPassed Result = iota
	ResultFailed
	ResultError
)

type Progress struct {
	benchmarkID string
	stdout      string
	stderr      string
	score       int
	startedAt   time.Time
}

func NewProgress(benchmarkID string, stdout string, stderr string, score int, startedAt time.Time) *Progress {
	return &Progress{
		benchmarkID: benchmarkID,
		stdout:      stdout,
		stderr:      stderr,
		score:       score,
		startedAt:   startedAt,
	}
}

func (p *Progress) GetBenchmarkID() string {
	return p.benchmarkID
}

func (p *Progress) GetStdout() string {
	return p.stdout
}

func (p *Progress) GetStderr() string {
	return p.stderr
}

func (p *Progress) GetScore() int {
	return p.score
}

func (p *Progress) GetStartedAt() time.Time {
	return p.startedAt
}
