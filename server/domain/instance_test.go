package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func TestNewInstance(t *testing.T) {
	teamID := uuid.New()
	factory := domain.NewInstanceFactory(3)

	cases := []struct {
		name     string
		existing []domain.Instance
		expected int
	}{
		{
			name:     "no existing instances",
			existing: []domain.Instance{},
			expected: 1,
		},
		{
			name: "one existing instance",
			existing: []domain.Instance{
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  1,
					Status: domain.InstanceStatusRunning,
				},
			},
			expected: 2,
		},
		{
			name: "two existing instances",
			existing: []domain.Instance{
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  1,
					Status: domain.InstanceStatusRunning,
				},
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  3,
					Status: domain.InstanceStatusRunning,
				},
			},
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance, err := factory.Create(teamID, c.existing)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if instance.Index != c.expected {
				t.Fatalf("expected instance number %d, got %d", c.expected, instance.Index)
			}
		})
	}
}
