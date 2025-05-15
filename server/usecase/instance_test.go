package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	instancemock "github.com/traPtitech/piscon-portal-v2/server/services/instance/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestCreateInstance(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	repo := repomock.NewMockRepository(ctrl)
	manager := instancemock.NewMockManager(ctrl)
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

	repo := repomock.NewMockRepository(ctrl)
	manager := instancemock.NewMockManager(ctrl)
	instanceUseCase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

	teamID := uuid.New()

	repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	repo.EXPECT().GetTeamInstances(gomock.Any(), teamID).Times(1).Return([]domain.Instance{
		{Index: 1}, {Index: 2}, {Index: 3},
	}, nil)

	_, err := instanceUseCase.CreateInstance(t.Context(), teamID)
	assert.ErrorIs(t, err, usecase.NewUseCaseError(domain.ErrInstanceLimitExceeded))
}

func TestDeleteInstance(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := repomock.NewMockRepository(ctrl)
	manager := instancemock.NewMockManager(ctrl)
	usecase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

	instanceID := uuid.New()
	providerID := uuid.New()
	instance := domain.Instance{
		ID: instanceID,
		Infra: domain.InfraInstance{
			ProviderInstanceID: providerID.String(),
			Status:             domain.InstanceStatusDeleted,
		},
	}

	repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	repo.EXPECT().FindInstance(gomock.Any(), instanceID).Return(instance, nil)
	manager.EXPECT().Delete(
		gomock.Any(),
		gomock.Cond(func(instance domain.InfraInstance) bool {
			return instance.ProviderInstanceID == providerID.String()
		}),
	).Return(domain.InfraInstance{}, nil)
	repo.EXPECT().UpdateInstance(gomock.Any(), gomock.Any()).Return(nil)

	err := usecase.DeleteInstance(t.Context(), instanceID)
	assert.NoError(t, err)
}
