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
	"github.com/traPtitech/piscon-portal-v2/server/utils/logpreview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip" // gRPCのgzip展開を登録する。
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
	job, err := bs.b.StartBenchmark(ctx)
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
			BenchmarkId:         job.Benchmark.ID.String(),
			TargetIpAddress:     *job.Benchmark.Instance.Infra.PrivateIP,
			InstanceIpAddresses: job.InstanceIPAddresses,
		},
	}, nil
}

func (bs *BenchmarkService) SendBenchmarkProgress(stream grpc.ClientStreamingServer[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]) error {
	ctx := stream.Context()
	var stdout, stderr logpreview.Buffer
	var streamID string
	var moreChunks bool
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if moreChunks {
				return status.Error(codes.InvalidArgument, "incomplete log batch")
			}
			break
		}
		if err != nil {
			return handleError("failed to receive progress", err)
		}

		if !req.LogsAreDelta {
			return status.Error(codes.FailedPrecondition, "runner must support delta log streaming; update runner")
		}

		benchID, err := uuid.Parse(req.BenchmarkId)
		if err != nil {
			return status.Error(codes.InvalidArgument, "invalid benchmark ID")
		}

		if streamID == "" {
			streamID = req.BenchmarkId
		} else if streamID != req.BenchmarkId {
			return status.Error(codes.InvalidArgument, "benchmark ID must not change within a stream")
		}
		stdout.WriteString(string(req.StdoutDelta))
		stderr.WriteString(string(req.StderrDelta))
		moreChunks = req.MoreLogChunks
		if moreChunks {
			continue
		}
		err = bs.b.SaveBenchmarkProgress(ctx, benchID, domain.BenchmarkLog{
			UserLog:  stdout.String(),
			AdminLog: stderr.String(),
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

	if err := stream.SendAndClose(&portalv1.SendBenchmarkProgressResponse{}); err != nil {
		return fmt.Errorf("send progress response: %w", err)
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
