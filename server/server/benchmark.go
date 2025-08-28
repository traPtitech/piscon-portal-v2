package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			BenchmarkId:     bench.ID.String(),
			TargetIpAddress: *bench.Instance.Infra.PrivateIP,
		},
	}, nil
}

func (bs *BenchmarkService) SendBenchmarkProgress(stream grpc.ClientStreamingServer[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]) error {
	ctx := stream.Context()
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return handleError("failed to receive progress", err)
		}

		benchID, err := uuid.Parse(req.BenchmarkId)
		if err != nil {
			return status.Error(codes.InvalidArgument, "invalid benchmark ID")
		}

		err = bs.b.SaveBenchmarkProgress(ctx, benchID, domain.BenchmarkLog{
			UserLog:  req.Stdout,
			AdminLog: req.Stderr,
		}, req.Score, req.StartedAt.AsTime())
		if usecase.IsUseCaseError(err) {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid benchmark status: %v", err))
		}
		if errors.Is(err, usecase.ErrNotFound) {
			return status.Error(codes.NotFound, "benchmark not found")
		}
		if err != nil {
			return handleError("failed to save benchmark progress", err)
		}
	}

	return nil
}

func (bs *BenchmarkService) PostJobFinished(ctx context.Context, req *portalv1.PostJobFinishedRequest) (*portalv1.PostJobFinishedResponse, error) {
	benchID, err := uuid.Parse(req.BenchmarkId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid benchmark ID")
	}

	var result domain.BenchmarkResult
	switch req.Result {
	case portalv1.BenchmarkResult_BENCHMARK_RESULT_PASSED:
		result = domain.BenchmarkResultStatusPassed
	case portalv1.BenchmarkResult_BENCHMARK_RESULT_FAILED:
		result = domain.BenchmarkResultStatusFailed
	case portalv1.BenchmarkResult_BENCHMARK_RESULT_ERROR:
		result = domain.BenchmarkResultStatusError
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid benchmark result")
	}

	err = bs.b.FinalizeBenchmark(ctx, benchID, result, req.FinishedAt.AsTime(), req.RunnerError)
	if usecase.IsUseCaseError(err) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid benchmark status: %v", err))
	}
	if errors.Is(err, usecase.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "benchmark not found")
	}
	if err != nil {
		return nil, handleError("failed to finalize benchmark", err)
	}

	return &portalv1.PostJobFinishedResponse{}, nil
}
