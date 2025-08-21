package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/gen/portal/v1/mock"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetJob(t *testing.T) {
	t.Parallel()

	type GetJobCall struct {
		jobRes *portalv1.GetBenchmarkJobResponse
		err    error
	}

	benchmarkJob := &portalv1.BenchmarkJob{
		BenchmarkId:     "benchmark-id",
		TargetIpAddress: "target-ip",
	}
	job := domain.NewJob(benchmarkJob.BenchmarkId, benchmarkJob.TargetIpAddress)

	testCases := map[string]struct {
		GetJobCalls []GetJobCall
		job         *domain.Job
		err         error
	}{
		"1回のリクエストで取得できる": {
			GetJobCalls: []GetJobCall{
				{
					jobRes: &portalv1.GetBenchmarkJobResponse{
						BenchmarkJob: benchmarkJob,
					},
				},
			},
			job: job,
		},
		"2回目で取得できる": {
			GetJobCalls: []GetJobCall{
				{
					jobRes: &portalv1.GetBenchmarkJobResponse{},
				},
				{
					jobRes: &portalv1.GetBenchmarkJobResponse{
						BenchmarkJob: benchmarkJob,
					},
				},
			},
			job: job,
		},
		"1回目でエラー": {
			GetJobCalls: []GetJobCall{
				{
					err: assert.AnError,
				},
			},
			err: assert.AnError,
		},
		"2回目でエラー": {
			GetJobCalls: []GetJobCall{
				{},
				{
					err: assert.AnError,
				},
			},
			err: assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mock.NewMockBenchmarkServiceClient(ctrl)

			calls := make([]any, 0, len(testCase.GetJobCalls))
			for _, call := range testCase.GetJobCalls {
				call := client.EXPECT().GetBenchmarkJob(gomock.Any(), &portalv1.GetBenchmarkJobRequest{}).Return(call.jobRes, call.err)
				calls = append(calls, call)
			}
			gomock.InOrder(calls...)

			portal := grpc.NewPortal(client, 0)

			ctx := context.Background()
			job, err := portal.GetJob(ctx)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.job, job)
		})
	}

}

func TestMakeProgressStreamClient(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		MakeProgressStreamClientErr error
		expectedErr                 error
	}{
		"success": {},
		"error": {
			MakeProgressStreamClientErr: assert.AnError,
			expectedErr:                 assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mock.NewMockBenchmarkServiceClient(ctrl)

			client.EXPECT().SendBenchmarkProgress(gomock.Any()).Return(nil, testCase.MakeProgressStreamClientErr)

			portal := grpc.NewPortal(client, 0)

			ctx := context.Background()
			streamClient, err := portal.MakeProgressStreamClient(ctx)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, streamClient)
			}
		})
	}
}
func TestPostJobFinished(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		execPostJobFinished bool
		jobID               string
		finishedAt          time.Time
		result              domain.Result
		runnerErr           error
		PostJobFinishedErr  error
		isError             bool
		expectedErr         error
	}{
		"結果がpassed": {
			jobID:               "job-id",
			finishedAt:          time.Now(),
			result:              domain.ResultPassed,
			runnerErr:           nil,
			execPostJobFinished: true,
		},
		"結果がfailed": {
			jobID:               "job-id",
			finishedAt:          time.Now(),
			result:              domain.ResultFailed,
			runnerErr:           nil,
			execPostJobFinished: true,
		},
		"結果がerror": {
			jobID:               "job-id",
			finishedAt:          time.Now(),
			result:              domain.ResultError,
			runnerErr:           errors.New("error"),
			execPostJobFinished: true,
		},
		"結果が不正なのでエラー": {
			jobID:               "job-id",
			finishedAt:          time.Now(),
			result:              100,
			runnerErr:           nil,
			isError:             true,
			execPostJobFinished: false,
		},
		"PostJobFinishedでエラー": {
			jobID:               "job-id",
			finishedAt:          time.Now(),
			result:              domain.ResultPassed,
			runnerErr:           nil,
			execPostJobFinished: true,
			isError:             true,
			PostJobFinishedErr:  assert.AnError,
			expectedErr:         assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mock.NewMockBenchmarkServiceClient(ctrl)

			if testCase.execPostJobFinished {
				var result portalv1.BenchmarkResult
				switch testCase.result {
				case domain.ResultPassed:
					result = portalv1.BenchmarkResult_BENCHMARK_RESULT_PASSED
				case domain.ResultFailed:
					result = portalv1.BenchmarkResult_BENCHMARK_RESULT_FAILED
				case domain.ResultError:
					result = portalv1.BenchmarkResult_BENCHMARK_RESULT_ERROR
				default:
					t.Fatalf("unknown result: %v", testCase.result)
				}
				var runnerError *string
				if testCase.runnerErr != nil {
					errStr := testCase.runnerErr.Error()
					runnerError = &errStr
				}
				client.EXPECT().PostJobFinished(gomock.Any(), &portalv1.PostJobFinishedRequest{
					BenchmarkId: testCase.jobID,
					FinishedAt:  timestamppb.New(testCase.finishedAt),
					Result:      result,
					RunnerError: runnerError,
				}).Return(nil, testCase.expectedErr)
			}

			portal := grpc.NewPortal(client, 0)

			ctx := context.Background()
			err := portal.PostJobFinished(ctx, testCase.jobID, testCase.finishedAt, testCase.result, testCase.runnerErr)

			if testCase.isError {
				if testCase.expectedErr != nil {
					assert.ErrorIs(t, err, testCase.expectedErr)
				} else {
					assert.Error(t, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
