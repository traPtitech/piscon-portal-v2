package server_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/server"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
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
					PrivateIP: targetIP,
				},
			},
			job: &portalv1.BenchmarkJob{
				BenchmarkId: benchID.String(),
				TargetUrl:   targetIP,
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
