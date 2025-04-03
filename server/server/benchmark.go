package server

import (
	"context"
	"errors"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"google.golang.org/grpc"
)

type BenchmarkService struct {
	b usecase.BenchmarkUseCase
	portalv1.UnimplementedBenchmarkServiceServer
}

func NewBenchmarkService(u usecase.BenchmarkUseCase) *BenchmarkService {
	return &BenchmarkService{
		b: u,
	}
}

var _ portalv1.BenchmarkServiceServer = (*BenchmarkService)(nil)

func (bs *BenchmarkService) GetBenchmarkJob(ctx context.Context, _ *portalv1.GetBenchmarkJobRequest) (*portalv1.GetBenchmarkJobResponse, error) {
	bench, err := bs.b.StartBenchmark(ctx)
	if errors.Is(err, usecase.ErrNotFound) {
		return &portalv1.GetBenchmarkJobResponse{
			BenchmarkJob: nil,
		}, nil
	}
	if err != nil {
		return nil, handleError("failed to start benchmark", err)
	}

	return &portalv1.GetBenchmarkJobResponse{
		BenchmarkJob: &portalv1.BenchmarkJob{
			BenchmarkId: bench.ID.String(),
			TargetUrl:   bench.Instance.PrivateIP,
		},
	}, nil
}

func (bs *BenchmarkService) SendBenchmarkProgress(grpc.ClientStreamingServer[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]) error {
	// TODO
	return nil
}
func (bs *BenchmarkService) PostJobFinished(context.Context, *portalv1.PostJobFinishedRequest) (*portalv1.PostJobFinishedResponse, error) {
	// TODO
	return nil, nil
}
