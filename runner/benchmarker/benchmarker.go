package benchmarker

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true

import (
	"context"
	"io"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type Benchmarker interface {
	// Start starts the benchmark job, but does not wait for it to complete.
	// It returns the stdout and stderr [io.ReadCloser] of the job.
	// When the job finishes, the io.ReadCloser will be closed and Read method will return io.EOF.
	// This mechanism is like [os/exec.Cmd.Start].
	// If ctx is canceled, the job must be stopped.
	//
	// After a successful call to Start the [Benchmarker.Wait] method must be called in order to release associated system resources.
	Start(ctx context.Context, job *domain.Job) (out Outputs, startedAt time.Time, err error)
	// Wait waits for the job to complete and returns an error if the job failed.
	Wait(ctx context.Context) (domain.Result, time.Time, error)
	// CalculateScore calculates the score of the benchmark job based on stdout and stderr content.
	CalculateScore(ctx context.Context, allStdout, allStderr string) (int, error)
}

type Outputs struct {
	Stdout io.Reader
	Stderr io.Reader
}
