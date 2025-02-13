package grpc_test

import (
	"context"
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

func TestSendProgress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		progress *domain.Progress
		SendErr  error
		err      error
	}{
		"正しく送信できる": {
			progress: domain.NewProgress("benchmark-id", "stdout", "stderr", 100, time.Now()),
		},
		"エラーが返ってくる": {
			progress: domain.NewProgress("benchmark-id", "stdout", "stderr", 100, time.Now()),
			SendErr:  assert.AnError,
			err:      assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			streamClient := mock.NewMockClientStreamingClient[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse](ctrl)

			streamClient.EXPECT().Send(gomock.Eq(&portalv1.SendBenchmarkProgressRequest{
				BenchmarkId: testCase.progress.GetBenchmarkID(),
				StartedAt:   timestamppb.New(testCase.progress.GetStartedAt()),
				Stdout:      testCase.progress.GetStdout(),
				Stderr:      testCase.progress.GetStderr(),
				Score:       int64(testCase.progress.GetScore()),
			})).Return(testCase.SendErr)

			client := grpc.NewProgressStreamClient(streamClient)

			ctx := context.Background()
			err := client.SendProgress(ctx, testCase.progress)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
