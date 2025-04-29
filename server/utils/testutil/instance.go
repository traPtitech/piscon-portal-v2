package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func CompareInstance(t *testing.T, want, got domain.Instance) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID, "instance.ID mismatch")
	assert.Equal(t, want.TeamID, got.TeamID, "instance.TeamID mismatch")
	assert.Equal(t, want.Index, got.Index, "instance.InstanceNumber mismatch")
	assert.Equal(t, want.Status, got.Status, "instance.Status mismatch")
	assert.Equal(t, want.Infra.PrivateIP, got.Infra.PrivateIP, "instance.PrivateIP mismatch")
	assert.Equal(t, want.Infra.PublicIP, got.Infra.PublicIP, "instance.PublicIP mismatch")
}
