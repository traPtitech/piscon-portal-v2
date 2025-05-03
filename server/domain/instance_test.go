package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

func TestCreateInstance(t *testing.T) {
	teamID := uuid.New()
	factory := domain.NewInstanceFactory(3)

	cases := []struct {
		name          string
		existing      []domain.Instance
		expectedIndex int
	}{
		{
			name:          "no existing instances",
			existing:      []domain.Instance{},
			expectedIndex: 1,
		},
		{
			name: "one existing instance",
			existing: []domain.Instance{
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  1,
				},
			},
			expectedIndex: 2,
		},
		{
			name: "two existing instances",
			existing: []domain.Instance{
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  1,
				},
				{
					ID:     uuid.New(),
					TeamID: teamID,
					Index:  3,
				},
			},
			expectedIndex: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			instance, err := factory.Create(teamID, c.existing)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if instance.Index != c.expectedIndex {
				t.Fatalf("expected instance number %d, got %d", c.expectedIndex, instance.Index)
			}
		})
	}
}

func TestCreateInstance_ExceedsInstanceLimit(t *testing.T) {
	teamID := uuid.New()
	factory := domain.NewInstanceFactory(3)

	existing := []domain.Instance{
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  1,
		},
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  2,
		},
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  3,
		},
	}
	_, err := factory.Create(teamID, existing)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != domain.ErrInstanceLimitExceeded {
		t.Fatalf("expected error %v, got %v", domain.ErrInstanceLimitExceeded, err)
	}
}
