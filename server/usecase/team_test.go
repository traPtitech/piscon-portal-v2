package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := mock.NewMockRepository(ctrl)
	teamUseCase := usecase.NewTeamUseCase(mockRepo)

	mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.User, error) {
			return domain.User{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any()).AnyTimes()

	userID := uuid.New()

	tests := []struct {
		name    string
		input   usecase.CreateTeamInput
		wantErr bool
	}{
		{
			name: "valid",
			input: usecase.CreateTeamInput{
				Name:      "valid-test-team",
				MemberIDs: []uuid.UUID{userID},
				CreatorID: userID,
			},
			wantErr: false,
		},
		{
			name: "multiple members",
			input: usecase.CreateTeamInput{
				Name:      "multiple-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New()},
				CreatorID: userID,
			},
			wantErr: false,
		},
		{
			name: "4 members",
			input: usecase.CreateTeamInput{
				Name:      "4-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New(), uuid.New(), uuid.New()},
				CreatorID: userID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := teamUseCase.CreateTeam(context.TODO(), tt.input)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && tt.wantErr {
				t.Error("expected error, but got nil")
			}
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := mock.NewMockRepository(ctrl)
	teamUseCase := usecase.NewTeamUseCase(mockRepo)

	mockRepo.EXPECT().FindTeam(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.Team, error) {
			return domain.Team{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, id uuid.UUID) (domain.User, error) {
			return domain.User{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().UpdateTeam(gomock.Any(), gomock.Any()).AnyTimes()

	userID := uuid.New()

	tests := []struct {
		name    string
		input   usecase.UpdateTeamInput
		wantErr bool
	}{
		{
			name: "valid",
			input: usecase.UpdateTeamInput{
				ID:        uuid.New(),
				Name:      "valid-test-team",
				MemberIDs: []uuid.UUID{userID},
			},
			wantErr: false,
		},
		{
			name: "4 members",
			input: usecase.UpdateTeamInput{
				ID:        uuid.New(),
				Name:      "4-members-test-team",
				MemberIDs: []uuid.UUID{userID, uuid.New(), uuid.New(), uuid.New()},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := teamUseCase.UpdateTeam(context.TODO(), tt.input)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && tt.wantErr {
				t.Error("expected error, but got nil")
			}
		})
	}
}
