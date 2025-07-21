package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
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
	).Return(nil)
	repo.EXPECT().UpdateInstance(gomock.Any(), gomock.Any()).Return(nil)

	err := usecase.DeleteInstance(t.Context(), instanceID)
	assert.NoError(t, err)
}

func TestDeleteInstance_instanceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := repomock.NewMockRepository(ctrl)
	manager := instancemock.NewMockManager(ctrl)
	instanceUsecase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

	instanceID := uuid.New()

	repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	repo.EXPECT().FindInstance(gomock.Any(), instanceID).Return(domain.Instance{}, repository.ErrNotFound)

	err := instanceUsecase.DeleteInstance(t.Context(), instanceID)
	assert.ErrorIs(t, err, usecase.ErrNotFound)
}

func TestUpdateInstance(t *testing.T) {
	ctrl := gomock.NewController(t)

	type fields struct {
		instance domain.Instance
		findErr  error
	}
	type args struct {
		op domain.InstanceOperation
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		expectManager func(manager *instancemock.MockManager, instance domain.Instance)
		expectErr     error
		expectErrAs   bool
	}{
		{
			name: "success (start)",
			fields: fields{
				instance: func() domain.Instance {
					id := uuid.New()
					providerID := uuid.New().String()
					return domain.Instance{
						ID: id,
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerID,
							Status:             domain.InstanceStatusStopped,
						},
					}
				}(),
			},
			args: args{
				op: domain.InstanceOperationStart,
			},
			expectManager: func(manager *instancemock.MockManager, instance domain.Instance) {
				manager.EXPECT().Start(gomock.Any(), gomock.Cond(func(infra domain.InfraInstance) bool {
					return infra.ProviderInstanceID == instance.Infra.ProviderInstanceID
				})).Return(nil)
			},
			expectErr:   nil,
			expectErrAs: false,
		},
		{
			name: "success (stop)",
			fields: fields{
				instance: func() domain.Instance {
					id := uuid.New()
					providerID := uuid.New().String()
					return domain.Instance{
						ID: id,
						Infra: domain.InfraInstance{
							ProviderInstanceID: providerID,
							Status:             domain.InstanceStatusRunning,
						},
					}
				}(),
			},
			args: args{
				op: domain.InstanceOperationStop,
			},
			expectManager: func(manager *instancemock.MockManager, instance domain.Instance) {
				manager.EXPECT().Stop(gomock.Any(), gomock.Cond(func(infra domain.InfraInstance) bool {
					return infra.ProviderInstanceID == instance.Infra.ProviderInstanceID
				})).Return(nil)
			},
			expectErr:   nil,
			expectErrAs: false,
		},
		{
			name: "instance not found",
			fields: fields{
				instance: domain.Instance{},
				findErr:  repository.ErrNotFound,
			},
			args: args{
				op: domain.InstanceOperationStart,
			},
			expectManager: nil,
			expectErr:     usecase.ErrNotFound,
			expectErrAs:   false,
		},
		{
			name: "invalid operation",
			fields: fields{
				instance: domain.Instance{},
			},
			args: args{
				op: 10000,
			},
			expectManager: nil,
			expectErr:     nil,
			expectErrAs:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repomock.NewMockRepository(ctrl)
			manager := instancemock.NewMockManager(ctrl)
			instanceUsecase := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

			instanceID := tt.fields.instance.ID
			if instanceID == uuid.Nil {
				instanceID = uuid.New()
			}

			repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
				return f(ctx)
			})
			repo.EXPECT().FindInstance(gomock.Any(), instanceID).Return(tt.fields.instance, tt.fields.findErr)

			if tt.expectManager != nil {
				tt.expectManager(manager, tt.fields.instance)
				repo.EXPECT().UpdateInstance(gomock.Any(), gomock.Any()).Return(nil)
			}

			err := instanceUsecase.UpdateInstance(t.Context(), instanceID, tt.args.op)
			if tt.expectErrAs {
				var ucErr usecase.UseCaseError
				assert.ErrorAs(t, err, &ucErr)
			} else if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
