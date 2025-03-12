package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
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
		setup       func(mockRepo *mock.MockRepository)
		expectError bool
	}{
		{
			name: "success: valid",
			setup: func(mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Status: domain.InstanceStatusRunning,
					}, nil)
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context, repository.Repository) error) error {
						return f(ctx, mockRepo)
					})
				mockRepo.EXPECT().
					GetBenchmarks(gomock.Any(), gomock.Any()).
					Return([]domain.Benchmark{}, nil)
				mockRepo.EXPECT().CreateBenchmark(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name: "failure: instance is not running",
			setup: func(mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Status: domain.InstanceStatusStopped,
					}, nil)
			},
			expectError: true,
		},
		{
			name: "failure: instance not found",
			setup: func(mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{}, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name: "failure: user's teamID does not match instance's teamID",
			setup: func(mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: uuid.New(),
						Status: domain.InstanceStatusRunning,
					}, nil)
			},
			expectError: true,
		},
		{
			name: "failure: benchmark already queued",
			setup: func(mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().
					FindUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.User{
						ID:     userID,
						TeamID: uuid.NullUUID{Valid: true, UUID: teamID},
					}, nil)
				mockRepo.EXPECT().
					FindInstance(gomock.Any(), gomock.Eq(instanceID)).
					Return(domain.Instance{
						ID:     instanceID,
						TeamID: teamID,
						Status: domain.InstanceStatusRunning,
					}, nil)
				mockRepo.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(context.Context, repository.Repository) error) error {
						return f(ctx, mockRepo)
					})
				mockRepo.EXPECT().
					GetBenchmarks(gomock.Any(), gomock.Any()).
					Return([]domain.Benchmark{
						{ID: uuid.New(), Instance: domain.Instance{ID: instanceID}, Status: domain.BenchmarkStatusWaiting},
					}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			useCase := usecase.NewBenchmarkUseCase(mockRepo)

			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			_, err := useCase.CreateBenchmark(context.Background(), instanceID, userID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
