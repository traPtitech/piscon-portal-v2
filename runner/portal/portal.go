package portal

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true

import (
	"context"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type Portal interface {
	// GetJob gets the next benchmark job.
	// If there is no job, it should block until a job is available.
	GetJob(ctx context.Context) (*domain.Job, error)
	// MakeProgressStreamClient creates a ProgressStreamClient.
	MakeProgressStreamClient(ctx context.Context) (ProgressStreamClient, error)
	// PostJobFinished posts the result of the job.
	PostJobFinished(ctx context.Context, jobID string, finishedAt time.Time, result domain.Result, runnerErr error) error
}

type ProgressStreamClient interface {
	// SendProgress sends the progress of the job.
	SendProgress(ctx context.Context, progress *domain.Progress) error
	// Close closes the client.
	Close() error
}
