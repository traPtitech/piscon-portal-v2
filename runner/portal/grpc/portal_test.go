package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/gen/portal/v1/mock"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	"go.uber.org/mock/gomock"
)

func TestGetJob(t *testing.T) {
	t.Parallel()

	type GetJobCall struct {
		jobRes *portalv1.GetBenchmarkJobResponse
		err    error
	}

	benchmarkJob := &portalv1.BenchmarkJob{
		BenchmarkId: "benchmark-id",
		TargetUrl:   "target-url",
	}
	job := domain.NewJob(benchmarkJob.BenchmarkId, benchmarkJob.TargetUrl)

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
