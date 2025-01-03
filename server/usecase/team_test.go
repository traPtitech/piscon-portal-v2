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
		DoAndReturn(func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().CreateTeam(gomock.Any(), gomock.Any()).AnyTimes()

	userID := uuid.NewString()

	tests := []struct {
		name    string
		input   usecase.CreateTeamInput
		wantErr bool
	}{
		{
			name: "valid",
			input: usecase.CreateTeamInput{
				Name:      "valid-test-team",
				MemberIDs: []string{userID},
				CreatorID: userID,
			},
			wantErr: false,
		},
		{
			name: "multiple members",
			input: usecase.CreateTeamInput{
				Name:      "multiple-members-test-team",
				MemberIDs: []string{userID, uuid.NewString()},
				CreatorID: userID,
			},
			wantErr: false,
		},
		{
			name: "4 members",
			input: usecase.CreateTeamInput{
				Name:      "4-members-test-team",
				MemberIDs: []string{userID, uuid.NewString(), uuid.NewString(), uuid.NewString()},
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
		DoAndReturn(func(ctx context.Context, id string) (domain.Team, error) {
			return domain.Team{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{ID: id}, nil
		}).
		AnyTimes()
	mockRepo.EXPECT().UpdateTeam(gomock.Any(), gomock.Any()).AnyTimes()

	userID := uuid.NewString()

	tests := []struct {
		name    string
		input   usecase.UpdateTeamInput
		wantErr bool
	}{
		{
			name: "valid",
			input: usecase.UpdateTeamInput{
				ID:        uuid.NewString(),
				Name:      "valid-test-team",
				MemberIDs: []string{userID},
			},
			wantErr: false,
		},
		{
			name: "4 members",
			input: usecase.UpdateTeamInput{
				ID:        uuid.NewString(),
				Name:      "4-members-test-team",
				MemberIDs: []string{userID, uuid.NewString(), uuid.NewString(), uuid.NewString()},
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
