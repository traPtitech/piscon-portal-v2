package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestCreateBenchmark(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	teamID := uuid.New()
	instanceID := uuid.New()

	tests := []struct {
		name             string
		user             domain.User
		instance         domain.Instance
		queuedBenchmarks []domain.Benchmark
		expectError      bool
	}{
		{
			name: "success: valid",
			user: domain.User{
				ID:     userID,
				TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
			},
			instance: domain.Instance{
				ID:     instanceID,
				TeamID: teamID,
				Status: domain.InstanceStatusRunning,
			},
			expectError: false,
		},
		{
			name: "failure: instance is not running",
			user: domain.User{ID: userID,
				TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
			},
			instance: domain.Instance{
				ID:     instanceID,
				TeamID: teamID, Status: domain.InstanceStatusStopped,
			},
			expectError: true,
		},
		{
			name: "failure: user's teamID does not match instance's teamID",
			user: domain.User{
				ID:     userID,
				TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
			},
			instance: domain.Instance{
				ID:     instanceID,
				TeamID: uuid.New(), // invalid teamID
				Status: domain.InstanceStatusRunning,
			},
			expectError: true,
		},
		{
			name: "failure: benchmark already queued",
			user: domain.User{
				ID:     userID,
				TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
			},
			instance: domain.Instance{
				ID:     instanceID,
				TeamID: teamID,
				Status: domain.InstanceStatusRunning,
			},
			queuedBenchmarks: []domain.Benchmark{
				{ID: uuid.New(), Instance: domain.Instance{ID: instanceID}, Status: domain.BenchmarkStatusWaiting},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			mockRepo := mock.NewMockRepository(ctrl)
			useCase := usecase.NewBenchmarkUseCase(mockRepo)

			mockRepo.EXPECT().
				FindUser(gomock.Any(), gomock.Eq(userID)).
				Return(tt.user, nil)
			mockRepo.EXPECT().
				FindInstance(gomock.Any(), gomock.Eq(instanceID)).
				Return(tt.instance, nil)
			mockRepo.EXPECT().
				GetBenchmarks(gomock.Any(), gomock.Any()).
				Return(tt.queuedBenchmarks, nil)
			mockRepo.EXPECT().CreateBenchmark(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			_, err := useCase.CreateBenchmark(context.Background(), instanceID, userID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
