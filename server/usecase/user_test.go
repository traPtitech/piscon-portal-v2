package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/services/traq"
	traqmock "github.com/traPtitech/piscon-portal-v2/server/services/traq/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	mockSvc := traqmock.NewMockService(ctrl)
	usecase := usecase.NewUserUseCase(mockRepo, mockSvc)

	traqUsers := []traq.User{
		{ID: uuid.New(), Name: "user1"},
		{ID: uuid.New(), Name: "user2"},
		{ID: uuid.New(), Name: "user3"},
	}
	portalUsers := []domain.User{
		{ID: traqUsers[0].ID, Name: traqUsers[0].Name, TeamID: uuid.NullUUID{UUID: uuid.New(), Valid: true}},
	}

	mockSvc.EXPECT().GetUsers(gomock.Any()).Return(traqUsers, nil)
	mockRepo.EXPECT().GetUsers(gomock.Any()).Return(portalUsers, nil)

	got, err := usecase.GetUsers(context.Background())
	assert.NoError(t, err)

	want := []domain.User{
		{
			ID:     portalUsers[0].ID,
			Name:   portalUsers[0].Name,
			TeamID: portalUsers[0].TeamID,
		},
		{
			ID:   traqUsers[1].ID,
			Name: traqUsers[1].Name,
		},
		{
			ID:   traqUsers[2].ID,
			Name: traqUsers[2].Name,
		},
	}

	testutil.CompareUsers(t, want, got)
}
