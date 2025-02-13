package grpc

import (
	"context"
	"fmt"
	"time"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
)

type Portal struct {
	cl              portalv1.BenchmarkServiceClient
	pollingInterval time.Duration
}

func NewPortal(cl portalv1.BenchmarkServiceClient, pollingInterval time.Duration) *Portal {
	return &Portal{
		cl:              cl,
		pollingInterval: pollingInterval,
	}
}

var _ portal.Portal = (*Portal)(nil)

func (p *Portal) GetJob(ctx context.Context) (*domain.Job, error) {
	for {
		jobRes, err := p.cl.GetBenchmarkJob(ctx, &portalv1.GetBenchmarkJobRequest{})
		if err != nil {
			return nil, fmt.Errorf("get benchmark job: %w", err)
		}

		job := jobRes.GetBenchmarkJob()
		if job != nil {
			return domain.NewJob(job.GetBenchmarkId(), job.GetTargetUrl()), nil
		}

		time.Sleep(p.pollingInterval)
	}
}

func (p *Portal) MakeProgressStreamClient(ctx context.Context) (portal.ProgressStreamClient, error) {
	streamClient, err := p.cl.SendBenchmarkProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("send benchmark progress: %w", err)
	}

	return &ProgressStreamClient{cl: streamClient}, nil
}

func (p *Portal) PostJobFinished(ctx context.Context, jobID string, finishedAt time.Time, result domain.Result, runnerErr error) error {
	return nil
}
