package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
)

func TestFindInstance(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instanceID := uuid.New()
	teamID := uuid.New()
	instance := domain.Instance{
		ID:     instanceID,
		TeamID: teamID,
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "prov-instance-id",
		},
		CreatedAt: time.Now(),
	}
	mustMakeInstance(t, db, instance)

	got, err := repo.FindInstance(t.Context(), instance.ID)
	assert.NoError(t, err)

	testutil.CompareInstance(t, instance, got)
}

func TestFindInstance_NotFound(t *testing.T) {
	t.Parallel()

	repo, _ := setupRepository(t)

	id := uuid.New()
	_, err := repo.FindInstance(t.Context(), id)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestCreateInstance(t *testing.T) {
	t.Parallel()

	repo, _ := setupRepository(t)

	instanceID := uuid.New()
	teamID := uuid.New()
	instance := domain.Instance{
		ID:     instanceID,
		TeamID: teamID,
		Index:  2,
		Infra: domain.InfraInstance{
			ProviderInstanceID: "prov-instance-id",
		},
		CreatedAt: time.Now(),
	}

	err := repo.CreateInstance(t.Context(), instance)
	assert.NoError(t, err)

	got, err := repo.FindInstance(t.Context(), instance.ID)
	assert.NoError(t, err)
	testutil.CompareInstance(t, instance, got)
}

func TestGetTeamInstances(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	teamID := uuid.New()
	otherTeamID := uuid.New()
	instances := []domain.Instance{
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  1,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-1",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  2,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-2",
			},
			CreatedAt: time.Now(),
			DeletedAt: ptr.Of(time.Now()),
		},
		{
			ID:     uuid.New(),
			TeamID: otherTeamID,
			Index:  1,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-3",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: otherTeamID,
			Index:  2,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-4",
			},
			CreatedAt: time.Now(),
			DeletedAt: ptr.Of(time.Now()),
		},
	}
	for _, inst := range instances {
		mustMakeInstance(t, db, inst)
	}

	got, err := repo.GetTeamInstances(t.Context(), teamID)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	testutil.ContainsInstance(t, got, instances[0])
	testutil.ContainsInstance(t, got, instances[1])
}

func TestGetAllInstances(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instances := []domain.Instance{
		{
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  1,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-1",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  2,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-2",
			},
			CreatedAt: time.Now(),
		},
		{
			// Deleted instance is included in the results
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  3,
			Infra: domain.InfraInstance{
				ProviderInstanceID: "prov-instance-id-3",
			},
			CreatedAt: time.Now(),
			DeletedAt: ptr.Of(time.Now()),
		},
	}
	for _, inst := range instances {
		mustMakeInstance(t, db, inst)
	}

	got, err := repo.GetAllInstances(t.Context())
	assert.NoError(t, err)
	assert.Len(t, got, 3)
	for _, want := range instances {
		testutil.ContainsInstance(t, got, want)
	}
}

func TestDeleteInstance(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instanceID1 := uuid.New()
	instanceID2 := uuid.New()
	providerID1 := uuid.New()
	providerID2 := uuid.New()

	instance1 := domain.Instance{
		ID:     instanceID1,
		TeamID: uuid.New(),
		Index:  1,
		Infra: domain.InfraInstance{
			ProviderInstanceID: providerID1.String(),
		},
		CreatedAt: time.Now(),
	}
	instance2 := domain.Instance{
		ID:     instanceID2,
		TeamID: uuid.New(),
		Index:  2,
		Infra: domain.InfraInstance{
			ProviderInstanceID: providerID2.String(),
		},
		CreatedAt: time.Now(),
	}

	testCases := map[string]struct {
		instanceID      uuid.UUID
		beforeInstances []domain.Instance
		afterInstances  []domain.Instance
		err             error
	}{
		"1台だけあるときに削除": {
			instanceID:      instanceID1,
			beforeInstances: []domain.Instance{instance1},
			afterInstances:  []domain.Instance{},
		},
		"複数台あるときに削除": {
			instanceID:      instanceID1,
			beforeInstances: []domain.Instance{instance1, instance2},
			afterInstances:  []domain.Instance{instance2},
		},
		"存在しないIDを指定したのでErrNotFound": {
			instanceID:      uuid.New(),
			beforeInstances: []domain.Instance{instance1, instance2},
			afterInstances:  []domain.Instance{instance1, instance2},
			err:             repository.ErrNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, inst := range testCase.beforeInstances {
				mustMakeInstance(t, db, inst)
				t.Cleanup(func() {
					_, err := models.Instances.Delete(models.DeleteWhere.Instances.ID.EQ(inst.ID.String())).
						Exec(context.Background(), db)
					assert.NoError(t, err)
				})
			}

			err := repo.DeleteInstance(t.Context(), testCase.instanceID)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			afterInstances, err := models.Instances.Query(
				models.SelectWhere.Instances.DeletedAt.IsNull(),
			).All(t.Context(), db)
			assert.NoError(t, err)

			assert.Len(t, afterInstances, len(testCase.afterInstances))

			afterInstancesMap := make(map[uuid.UUID]domain.Instance, len(testCase.afterInstances))
			for _, inst := range afterInstances {
				assert.NoError(t, err)
				id, err := uuid.Parse(inst.ID)
				assert.NoError(t, err)
				teamID, err := uuid.Parse(inst.TeamID)
				assert.NoError(t, err)
				afterInstancesMap[id] = domain.Instance{
					ID:     id,
					TeamID: teamID,
					Index:  int(inst.InstanceNumber),
					Infra: domain.InfraInstance{
						ProviderInstanceID: inst.ProviderInstanceID,
					},
					CreatedAt: inst.CreatedAt,
				}
			}

			for _, want := range testCase.afterInstances {
				afterInst, ok := afterInstancesMap[want.ID]
				assert.True(t, ok, "instance %s not found in after instances", want.ID)
				testutil.CompareInstance(t, want, afterInst)
			}
		})
	}
}
