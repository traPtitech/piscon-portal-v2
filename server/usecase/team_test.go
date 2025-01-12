package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := mock.NewMockRepository(ctrl)
	useCase := usecase.New(mockRepo)

	userID := uuid.New()

	tests := []struct {
		name        string
		input       usecase.CreateTeamInput
		expectError bool
		setup       func()
	}{
		{
			name: "valid",
			input: usecase.CreateTeamInput{
				Name:      "valid-test-team",
				MemberIDs: []uuid.UUID{userID},
				CreatorID: userID,
			},
			expectError: false,
			setup: func() {
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.User, error) {
						return domain.User{ID: id}, nil
					})
				mockRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "multiple members",
			input: usecase.CreateTeamInput{
				Name:      "multiple-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New()},
				CreatorID: userID,
			},
			expectError: false,
			setup: func() {
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.User, error) {
						return domain.User{ID: id}, nil
					}).Times(2)
				mockRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "more than 3 members team is not allowed",
			input: usecase.CreateTeamInput{
				Name:      "4-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New(), uuid.New(), uuid.New()},
				CreatorID: userID,
			},
			expectError: true,
			setup: func() {
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.User, error) {
						return domain.User{ID: id}, nil
					}).Times(4)
			},
		},
		{
			name: "team member not found",
			input: usecase.CreateTeamInput{
				Name:      "user-not-found-test-team",
				MemberIDs: []uuid.UUID{userID},
				CreatorID: userID,
			},
			expectError: true,
			setup: func() {
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.CreateTeam(context.Background(), tt.input)
			if err != nil && !tt.expectError {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && tt.expectError {
				t.Error("expected error, but got nil")
			}
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := mock.NewMockRepository(ctrl)
	useCase := usecase.New(mockRepo)

	userID := uuid.New()

	tests := []struct {
		name        string
		input       usecase.UpdateTeamInput
		expectError bool
		setup       func()
	}{
		{
			name: "valid",
			input: usecase.UpdateTeamInput{
				ID:        uuid.New(),
				Name:      "valid-test-team",
				MemberIDs: []uuid.UUID{userID},
			},
			expectError: false,
			setup: func() {
				mockRepo.EXPECT().FindTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.Team, error) {
						return domain.Team{
							ID:      id,
							Members: []domain.User{{ID: userID, TeamID: uuid.NullUUID{UUID: id, Valid: true}}},
						}, nil
					})
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.User, error) {
						return domain.User{ID: id}, nil
					})
				mockRepo.EXPECT().UpdateTeam(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "more than 3 members team is not allowed",
			input: usecase.UpdateTeamInput{
				ID:        uuid.New(),
				Name:      "4-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New(), uuid.New(), uuid.New()},
			},
			expectError: true,
			setup: func() {
				mockRepo.EXPECT().FindTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.Team, error) {
						return domain.Team{ID: id}, nil
					})
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.User, error) {
						return domain.User{ID: id}, nil
					}).Times(4)
			},
		},
		{
			name: "team member not found",
			input: usecase.UpdateTeamInput{
				ID:        uuid.New(),
				Name:      "user-not-found-test-team",
				MemberIDs: []uuid.UUID{userID},
			},
			expectError: true,
			setup: func() {
				mockRepo.EXPECT().FindTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (domain.Team, error) {
						return domain.Team{ID: id}, nil
					})
				mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := useCase.UpdateTeam(context.Background(), tt.input)
			if err != nil && !tt.expectError {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && tt.expectError {
				t.Error("expected error, but got nil")
			}
		})
	}
}
