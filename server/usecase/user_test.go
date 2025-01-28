package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	usecase := usecase.NewUserUseCase(mockRepo)

	users := []domain.User{
		{ID: uuid.New(), Name: "user1", TeamID: uuid.NullUUID{UUID: uuid.New(), Valid: true}},
		{ID: uuid.New(), Name: "user2"},
		{ID: uuid.New(), Name: "user3"},
	}

	mockRepo.EXPECT().GetUsers(gomock.Any()).Return(users, nil)

	got, err := usecase.GetUsers(context.Background())
	assert.NoError(t, err)

	testutil.CompareUsers(t, users, got)
}
