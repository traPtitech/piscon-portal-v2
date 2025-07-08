package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func CompareInstance(t *testing.T, want, got domain.Instance) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID, "instance.ID mismatch")
	assert.Equal(t, want.TeamID, got.TeamID, "instance.TeamID mismatch")
	assert.Equal(t, want.Index, got.Index, "instance.Index mismatch")
	assert.Equal(t, want.Infra.ProviderInstanceID, got.Infra.ProviderInstanceID, "instance.Infra.ProviderInstanceID mismatch")
	assert.Equal(t, want.Infra.Status, got.Infra.Status, "instance.Infra.Status mismatch")
	assert.Equal(t, want.Infra.PrivateIP, got.Infra.PrivateIP, "instance.Infra.PrivateIP mismatch")
	assert.Equal(t, want.Infra.PublicIP, got.Infra.PublicIP, "instance.Infra.PublicIP mismatch")
	assert.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second, "instance.CreatedAt mismatch")
}

func ContainsInstance(t *testing.T, instances []domain.Instance, want domain.Instance) {
	t.Helper()

	for _, instance := range instances {
		if instance.ID == want.ID {
			CompareInstance(t, want, instance)
			return
		}
	}
	t.Errorf("instance with ID %s not found", want.ID)
}
