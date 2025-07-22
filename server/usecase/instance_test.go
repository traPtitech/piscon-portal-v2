package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	instancemock "github.com/traPtitech/piscon-portal-v2/server/services/instance/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
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
	repo.EXPECT().DeleteInstance(gomock.Any(), instanceID).Return(nil)

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

func TestGetInstance(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	instanceID := uuid.New()
	infraInstanceID := uuid.New().String()
	instance := domain.Instance{
		ID:     instanceID,
		TeamID: uuid.New(),
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: infraInstanceID,
		},
		CreatedAt: time.Now(),
	}
	fullInstance := instance
	fullInstance.Infra = domain.InfraInstance{
		ProviderInstanceID: infraInstanceID,
		Status:             domain.InstanceStatusRunning,
		PrivateIP:          ptr.Of("private ip"),
		PublicIP:           ptr.Of("public ip"),
	}

	testCases := map[string]struct {
		instanceID      uuid.UUID
		instance        domain.Instance
		FindInstanceErr error
		infraInstance   domain.InfraInstance
		executeGet      bool
		GetErr          error
		result          domain.Instance
		expectErr       error
	}{
		"正しく取得できる": {
			instanceID:    instanceID,
			instance:      instance,
			executeGet:    true,
			infraInstance: fullInstance.Infra,
			result:        fullInstance,
		},
		"FindInstanceでErrNotFoundなのでErrNotFound": {
			instanceID:      instanceID,
			FindInstanceErr: repository.ErrNotFound,
			expectErr:       usecase.ErrNotFound,
		},
		"FindInstanceがエラーなのでエラー": {
			instanceID:      instanceID,
			FindInstanceErr: assert.AnError,
			expectErr:       assert.AnError,
		},
		"Getでエラーなのでエラー": {
			instanceID: instanceID,
			instance:   instance,
			executeGet: true,
			GetErr:     assert.AnError,
			expectErr:  assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := repomock.NewMockRepository(ctrl)
			manager := instancemock.NewMockManager(ctrl)
			u := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

			repo.EXPECT().FindInstance(gomock.Any(), testCase.instanceID).
				Return(testCase.instance, testCase.FindInstanceErr)
			if testCase.executeGet {
				manager.EXPECT().Get(gomock.Any(), testCase.instance.Infra.ProviderInstanceID).
					Return(testCase.infraInstance, testCase.GetErr)
			}

			result, err := u.GetInstance(t.Context(), testCase.instanceID)
			if testCase.expectErr != nil {
				assert.ErrorIs(t, err, testCase.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.result, result)
			}
		})
	}
}

func TestGetTeamInstances(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	teamID := uuid.New()
	instance1 := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "infra1",
			Status:             domain.InstanceStatusRunning,
		},
	}
	instance2 := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  2,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "infra2",
			Status:             domain.InstanceStatusStopped,
		},
	}
	instance3 := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  3,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "infra3",
			Status:             domain.InstanceStatusDeleted,
		},
		DeletedAt: ptr.Of(time.Now()),
	}

	testCases := map[string]struct {
		teamID              uuid.UUID
		instances           []domain.Instance
		GetTeamInstancesErr error
		executeGetAll       bool
		infraInstances      []domain.InfraInstance
		GetAllErr           error
		result              []domain.Instance
		expectErr           error
	}{
		"正しく取得できる": {
			teamID:         teamID,
			instances:      []domain.Instance{instance1, instance2, instance3},
			executeGetAll:  true,
			infraInstances: []domain.InfraInstance{instance1.Infra, instance2.Infra},
			result:         []domain.Instance{instance1, instance2, instance3},
		},
		"GetTeamInstancesでエラー": {
			teamID:              teamID,
			GetTeamInstancesErr: assert.AnError,
			expectErr:           assert.AnError,
		},
		"GetAllでエラー": {
			teamID:        teamID,
			instances:     []domain.Instance{instance1, instance2},
			executeGetAll: true,
			GetAllErr:     assert.AnError,
			expectErr:     assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := repomock.NewMockRepository(ctrl)
			manager := instancemock.NewMockManager(ctrl)
			u := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

			repo.EXPECT().GetTeamInstances(gomock.Any(), testCase.teamID).
				Return(testCase.instances, testCase.GetTeamInstancesErr)

			if testCase.executeGetAll {
				ids := make([]string, 0, len(testCase.instances))
				for _, inst := range testCase.instances {
					ids = append(ids, inst.Infra.ProviderInstanceID)
				}
				manager.EXPECT().GetByIDs(gomock.Any(), ids).
					Return(testCase.infraInstances, testCase.GetAllErr)
			}

			result, err := u.GetTeamInstances(t.Context(), testCase.teamID)
			if testCase.expectErr != nil {
				assert.ErrorIs(t, err, testCase.expectErr)
			} else {
				assert.NoError(t, err)
			}

			if testCase.expectErr != nil {
				return
			}

			assert.Len(t, result, len(testCase.result))

			for i, inst := range result {
				testutil.CompareInstance(t, testCase.result[i], inst)
			}
		})
	}
}

func TestGetAllInstances(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	teamID := uuid.New()
	instance1 := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "infra1",
			Status:             domain.InstanceStatusRunning,
		},
	}
	instance2 := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  2,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "infra2",
			Status:             domain.InstanceStatusStopped,
		},
	}

	testCases := map[string]struct {
		instances          []domain.Instance
		GetAllInstancesErr error
		executeGetAll      bool
		infraInstances     []domain.InfraInstance
		GetAllErr          error
		result             []domain.Instance
		expectErr          error
	}{
		"正しく取得できる": {
			instances:      []domain.Instance{instance1, instance2},
			executeGetAll:  true,
			infraInstances: []domain.InfraInstance{instance1.Infra, instance2.Infra},
			result:         []domain.Instance{instance1, instance2},
		},
		"GetAllInstancesでエラー": {
			GetAllInstancesErr: assert.AnError,
			expectErr:          assert.AnError,
		},
		"GetAllでエラー": {
			instances:     []domain.Instance{instance1, instance2},
			executeGetAll: true,
			GetAllErr:     assert.AnError,
			expectErr:     assert.AnError,
		}}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := repomock.NewMockRepository(ctrl)
			manager := instancemock.NewMockManager(ctrl)
			u := usecase.NewInstanceUseCase(repo, domain.NewInstanceFactory(3), manager)

			repo.EXPECT().GetAllInstances(gomock.Any()).
				Return(testCase.instances, testCase.GetAllInstancesErr)

			if testCase.executeGetAll {
				manager.EXPECT().GetAll(gomock.Any()).
					Return(testCase.infraInstances, testCase.GetAllErr)
			}

			result, err := u.GetAllInstances(t.Context())
			if testCase.expectErr != nil {
				assert.ErrorIs(t, err, testCase.expectErr)
			} else {
				assert.NoError(t, err)
			}

			if testCase.expectErr != nil {
				return
			}

			assert.Len(t, result, len(testCase.result))

			for i, inst := range result {
				testutil.CompareInstance(t, testCase.result[i], inst)
			}
		})
	}
}
