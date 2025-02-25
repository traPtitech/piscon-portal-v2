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
		name        string
		user        domain.User
		instance    domain.Instance
		expectError bool
	}{
		{
			name: "Valid",
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
			name: "InstanceNotRunning",
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
			name: "InvalidInstance",
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
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			useCase := usecase.NewBenchmarkUseCase(mockRepo)

			mockRepo.EXPECT().
				FindUser(gomock.Any(), gomock.Eq(userID)).
				DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.User, error) {
					return tt.user, nil
				})
			mockRepo.EXPECT().
				FindInstance(gomock.Any(), gomock.Eq(instanceID)).
				DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.Instance, error) {
					return tt.instance, nil
				})
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
