package server_test

import (
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	mockportalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1/mock"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/server"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetBenchmarkJob(t *testing.T) {
	t.Parallel()

	benchID := uuid.New()
	targetIP := "0.0.0.0"

	testCases := map[string]struct {
		bench             domain.Benchmark
		StartBenchmarkErr error
		job               *portalv1.BenchmarkJob
		err               error
	}{
		"正常に取得できる": {
			bench: domain.Benchmark{
				ID: benchID,
				Instance: domain.Instance{
					Infra: domain.InfraInstance{
						PrivateIP: &targetIP,
					},
				},
			},
			job: &portalv1.BenchmarkJob{
				BenchmarkId:     benchID.String(),
				TargetIpAddress: targetIP,
			},
		},
		"StartBenchmarkがErrNotFound": {
			StartBenchmarkErr: usecase.ErrNotFound,
			job:               nil,
		},
		"StartBenchmarkがエラー": {
			StartBenchmarkErr: assert.AnError,
			err:               server.ExportedInternalError,
		},
	}

	ctrl := gomock.NewController(t)

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u := mock.NewMockUseCase(ctrl)
			s := server.NewBenchmarkService(u)

			u.EXPECT().StartBenchmark(gomock.Any()).Return(testCase.bench, testCase.StartBenchmarkErr)

			res, err := s.GetBenchmarkJob(t.Context(), &portalv1.GetBenchmarkJobRequest{})

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			if testCase.job != nil {
				assert.Equal(t, testCase.job, res.BenchmarkJob)
			} else {
				assert.Nil(t, res.BenchmarkJob)
			}
		})
	}
}

func TestSendBenchmarkProgress(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type (
		RecvResult struct {
			req  *portalv1.SendBenchmarkProgressRequest
			err  error
			save bool // 保存するもののかどうか
		}
		SaveBenchmarkProgressResult struct {
			err error
		}
	)

	benchID := uuid.New()
	startedAt := time.Now()
	req1 := &portalv1.SendBenchmarkProgressRequest{
		BenchmarkId: benchID.String(),
		Stdout:      "stdout",
		Stderr:      "stderr",
		Score:       100,
		StartedAt:   timestamppb.New(startedAt),
	}
	req2 := &portalv1.SendBenchmarkProgressRequest{
		BenchmarkId: benchID.String(),
		Stdout:      "stdout 2",
		Stderr:      "stderr 2",
		Score:       200,
		StartedAt:   timestamppb.New(startedAt),
	}

	testCases := map[string]struct {
		recvResults                 []RecvResult
		SaveBenchmarkProgressResult []SaveBenchmarkProgressResult
		isError                     bool
		err                         error
		errCode                     codes.Code
	}{
		"1回受信できる": {
			recvResults: []RecvResult{
				{req: req1, err: nil, save: true},
				{req: nil, err: io.EOF, save: false},
			},
			SaveBenchmarkProgressResult: []SaveBenchmarkProgressResult{
				{err: nil},
			},
		},
		"2回受信できる": {
			recvResults: []RecvResult{
				{req: req1, err: nil, save: true},
				{req: req2, err: nil, save: true},
				{req: nil, err: io.EOF, save: false},
			},
			SaveBenchmarkProgressResult: []SaveBenchmarkProgressResult{
				{err: nil}, {err: nil},
			},
		},
		"Recvがエラー": {
			recvResults: []RecvResult{{req: nil, err: assert.AnError, save: false}},
			isError:     true,
			err:         server.ExportedInternalError,
		},
		"idがUUIDでないのでエラー": {
			recvResults: []RecvResult{
				{req: &portalv1.SendBenchmarkProgressRequest{
					BenchmarkId: "invalid id",
					Stdout:      "stdout",
					Stderr:      "stderr",
					Score:       100,
					StartedAt:   timestamppb.New(startedAt),
				}, err: nil, save: false},
			},
			isError: true,
			errCode: codes.InvalidArgument,
		},
		"SaveBenchmarkProgressがUseCaseError": {
			recvResults: []RecvResult{{req: req1, err: nil, save: true}},
			SaveBenchmarkProgressResult: []SaveBenchmarkProgressResult{
				{err: usecase.NewUseCaseErrorFromMsg("error")},
			},
			isError: true,
			errCode: codes.InvalidArgument,
		},
		"SaveBenchmarkProgressがErrNotFound": {
			recvResults:                 []RecvResult{{req: req1, err: nil, save: true}},
			SaveBenchmarkProgressResult: []SaveBenchmarkProgressResult{{err: usecase.ErrNotFound}},
			isError:                     true,
			errCode:                     codes.NotFound,
		},
		"SaveBenchmarkProgressがエラー": {
			recvResults:                 []RecvResult{{req: req1, err: nil, save: true}},
			SaveBenchmarkProgressResult: []SaveBenchmarkProgressResult{{err: assert.AnError}},
			isError:                     true,
			err:                         server.ExportedInternalError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := mockportalv1.NewMockClientStreamingServer[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse](ctrl)
			u := mock.NewMockUseCase(ctrl)

			client.EXPECT().Context().Return(t.Context()).AnyTimes()

			for i, recvResult := range testCase.recvResults {
				client.EXPECT().Recv().Return(recvResult.req, recvResult.err)

				if recvResult.save {
					benchID, err := uuid.Parse(recvResult.req.BenchmarkId)
					require.NoError(t, err)

					saveResult := testCase.SaveBenchmarkProgressResult[i]

					benchLog := domain.BenchmarkLog{
						UserLog:  recvResult.req.Stdout,
						AdminLog: recvResult.req.Stderr,
					}
					u.EXPECT().
						SaveBenchmarkProgress(gomock.Any(), benchID, benchLog, recvResult.req.Score, gomock.Cond(
							func(t time.Time) bool {
								return t.Sub(recvResult.req.StartedAt.AsTime()).Abs() < time.Second // timestamppb.Newすると、情報量が減って比較できなくなる
							},
						)).
						Return(saveResult.err)
				}
			}

			s := server.NewBenchmarkService(u)

			err := s.SendBenchmarkProgress(client)
			if testCase.isError {
				assert.Error(t, err)
				if testCase.err != nil {
					assert.ErrorIs(t, err, testCase.err)
				}
				if testCase.errCode != codes.Code(0) {
					assert.Equal(t, testCase.errCode, status.Code(err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestPostJobFinished(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	benchID := uuid.New()
	req := &portalv1.PostJobFinishedRequest{
		BenchmarkId: benchID.String(),
		Result:      portalv1.BenchmarkResult_BENCHMARK_RESULT_PASSED,
		FinishedAt:  timestamppb.New(time.Now()),
	}

	testCases := map[string]struct {
		req                      *portalv1.PostJobFinishedRequest
		executeFinalizeBenchmark bool
		FinalizeBenchmarkErr     error
		isError                  bool
		err                      error
		errCode                  codes.Code
	}{
		"uuidが無効": {
			req: &portalv1.PostJobFinishedRequest{
				BenchmarkId: "invalid id",
			},
			isError: true,
			errCode: codes.InvalidArgument,
		},
		"無効なresult": {
			req: &portalv1.PostJobFinishedRequest{
				BenchmarkId: benchID.String(),
				Result:      100,
			},
			isError: true,
			errCode: codes.InvalidArgument,
		},
		"FinalizeBenchmarkがUseCaseError": {
			req:                      req,
			executeFinalizeBenchmark: true,
			FinalizeBenchmarkErr:     usecase.NewUseCaseErrorFromMsg("error"),
			isError:                  true,
			errCode:                  codes.InvalidArgument,
		},
		"FinalizeBenchmarkがErrNotFound": {
			req:                      req,
			executeFinalizeBenchmark: true,
			FinalizeBenchmarkErr:     usecase.ErrNotFound,
			isError:                  true,
			errCode:                  codes.NotFound,
		},
		"FinalizeBenchmarkがエラー": {
			req:                      req,
			executeFinalizeBenchmark: true,
			FinalizeBenchmarkErr:     assert.AnError,
			isError:                  true,
			errCode:                  codes.Internal,
		},
		"正常に終了": {
			req:                      req,
			executeFinalizeBenchmark: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u := mock.NewMockUseCase(ctrl)
			s := server.NewBenchmarkService(u)

			if testCase.executeFinalizeBenchmark {
				var result domain.BenchmarkResult
				switch testCase.req.Result {
				case portalv1.BenchmarkResult_BENCHMARK_RESULT_PASSED:
					result = domain.BenchmarkResultStatusPassed
				case portalv1.BenchmarkResult_BENCHMARK_RESULT_FAILED:
					result = domain.BenchmarkResultStatusFailed
				case portalv1.BenchmarkResult_BENCHMARK_RESULT_ERROR:
					result = domain.BenchmarkResultStatusError
				default:
					t.Fatalf("invalid benchmark result: %v", testCase.req.Result)
				}
				u.EXPECT().
					FinalizeBenchmark(gomock.Any(), benchID, result, testCase.req.FinishedAt.AsTime(), testCase.req.RunnerError).
					Return(testCase.FinalizeBenchmarkErr)
			}

			res, err := s.PostJobFinished(t.Context(), testCase.req)
			if testCase.isError {
				assert.Error(t, err)
				assert.Nil(t, res)
				if testCase.err != nil {
					assert.ErrorIs(t, err, testCase.err)
				}
				if testCase.errCode != codes.Code(0) {
					assert.Equal(t, testCase.errCode, status.Code(err))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}
