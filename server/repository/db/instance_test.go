package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/utils/testutil"
)

func TestFindInstance(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	instanceID := uuid.New()
	teamID := uuid.New()
	instance := domain.Instance{
		ID:             instanceID,
		TeamID:         teamID,
		InstanceNumber: 1,
		Status:         domain.InstanceStatusRunning,
		PrivateIP:      "192.0.2.0",
		PublicIP:       "192.0.2.0",
	}
	mustMakeInstance(t, db, instance)

	got, err := repo.FindInstance(t.Context(), instance.ID)
	assert.NoError(t, err)

	testutil.CompareInstance(t, instance, got)
}
