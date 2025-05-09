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

func TestCreateInstance(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	repo := mock.NewMockRepository(ctrl)
	manager := usecase.NewMockInstanceManager(ctrl)
	usecase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

	teamID := uuid.New()

	repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	repo.EXPECT().GetTeamInstances(gomock.Any(), teamID).Times(1)
	repo.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Times(1)

	manager.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	instance, err := usecase.CreateInstance(t.Context(), teamID)
	assert.NoError(t, err)

	assert.Equal(t, teamID, instance.TeamID)
}

func TestCreateInstance_tooManyInstances(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	repo := mock.NewMockRepository(ctrl)
	manager := usecase.NewMockInstanceManager(ctrl)
	usecase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

	teamID := uuid.New()

	repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	repo.EXPECT().GetTeamInstances(gomock.Any(), teamID).Times(1).Return([]domain.Instance{
		{Index: 1}, {Index: 2}, {Index: 3},
	}, nil)

	_, err := usecase.CreateInstance(t.Context(), teamID)
	assert.Error(t, err)
}
