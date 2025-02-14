package grpc

import (
	"context"
	"fmt"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProgressStreamClient struct {
	cl grpc.ClientStreamingClient[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]
}

func NewProgressStreamClient(cl grpc.ClientStreamingClient[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]) *ProgressStreamClient {
	return &ProgressStreamClient{cl: cl}
}

var _ portal.ProgressStreamClient = (*ProgressStreamClient)(nil)

func (c *ProgressStreamClient) SendProgress(_ context.Context, progress *domain.Progress) error {
	req := &portalv1.SendBenchmarkProgressRequest{
		BenchmarkId: progress.GetBenchmarkID(),
		StartedAt:   timestamppb.New(progress.GetStartedAt()),
		Stdout:      progress.GetStdout(),
		Stderr:      progress.GetStderr(),
		Score:       int64(progress.GetScore()),
	}

	return c.cl.Send(req)
}

func (c *ProgressStreamClient) Close() error {
	_, err := c.cl.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close progress stream: %w", err)
	}

	return nil
}
