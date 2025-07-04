package db_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
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
			Status:    domain.InstanceStatusRunning,
			PrivateIP: "192.0.2.0",
			PublicIP:  "192.0.2.0",
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
			Status:    domain.InstanceStatusBuilding,
			PrivateIP: "10.0.0.1",
			PublicIP:  "203.0.113.1",
		},
		CreatedAt: time.Now(),
	}

	err := repo.CreateInstance(t.Context(), instance)
	assert.NoError(t, err)

	got, err := repo.FindInstance(t.Context(), instance.ID)
	assert.NoError(t, err)
	testutil.CompareInstance(t, instance, got)
}

func TestUpdateInstance(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instanceID := uuid.New()
	teamID := uuid.New()
	instance := domain.Instance{
		ID:     instanceID,
		TeamID: teamID,
		Index:  3,
		Infra: domain.InfraInstance{
			Status:    domain.InstanceStatusRunning,
			PrivateIP: "10.0.0.2",
			PublicIP:  "203.0.113.2",
		},
		CreatedAt: time.Now(),
	}
	mustMakeInstance(t, db, instance)

	instance.Infra.Status = domain.InstanceStatusStopped
	instance.Infra.PrivateIP = "10.0.0.3"
	instance.Infra.PublicIP = "203.0.113.3"
	err := repo.UpdateInstance(t.Context(), instance)
	assert.NoError(t, err)

	got, err := repo.FindInstance(t.Context(), instance.ID)
	assert.NoError(t, err)
	testutil.CompareInstance(t, instance, got)
}

func TestUpdateInstance_NotFound(t *testing.T) {
	t.Parallel()

	repo, _ := setupRepository(t)

	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: uuid.New(),
		Index:  99,
		Infra: domain.InfraInstance{
			Status:    domain.InstanceStatusRunning,
			PrivateIP: "10.0.0.99",
			PublicIP:  "203.0.113.99",
		},
		CreatedAt: time.Now(),
	}
	err := repo.UpdateInstance(t.Context(), instance)
	assert.ErrorIs(t, err, repository.ErrNotFound)
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
				Status:    domain.InstanceStatusRunning,
				PrivateIP: "10.0.0.4",
				PublicIP:  "203.0.113.4",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  2,
			Infra: domain.InfraInstance{
				Status:    domain.InstanceStatusStopped,
				PrivateIP: "10.0.0.5",
				PublicIP:  "203.0.113.5",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: otherTeamID,
			Index:  1,
			Infra: domain.InfraInstance{
				Status:    domain.InstanceStatusRunning,
				PrivateIP: "10.0.0.6",
				PublicIP:  "203.0.113.6",
			},
			CreatedAt: time.Now(),
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
				Status:    domain.InstanceStatusRunning,
				PrivateIP: "10.0.0.7",
				PublicIP:  "203.0.113.7",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  2,
			Infra: domain.InfraInstance{
				Status:    domain.InstanceStatusStopped,
				PrivateIP: "10.0.0.8",
				PublicIP:  "203.0.113.8",
			},
			CreatedAt: time.Now(),
		},
	}
	for _, inst := range instances {
		mustMakeInstance(t, db, inst)
	}

	got, err := repo.GetAllInstances(t.Context())
	assert.NoError(t, err)
	assert.Len(t, got, len(instances))
	for _, want := range instances {
		testutil.ContainsInstance(t, got, want)
	}
}
