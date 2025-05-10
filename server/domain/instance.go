package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Instance struct {
	ID     uuid.UUID
	TeamID uuid.UUID
	Index  int
	Infra  InfraInstance
}

type InfraInstance struct {
	ProviderInstanceID string
	Status             InstanceStatus
	PrivateIP          string
	PublicIP           string
}

type InstanceStatus string

const (
	InstanceStatusRunning  InstanceStatus = "running"
	InstanceStatusBuilding InstanceStatus = "building"
	InstanceStatusStarting InstanceStatus = "starting"
	InstanceStatusStopping InstanceStatus = "stopping"
	InstanceStatusStopped  InstanceStatus = "stopped"
	InstanceStatusDeleting InstanceStatus = "deleting"
	InstanceStatusDeleted  InstanceStatus = "deleted"
)

type InstanceOperation int

const (
	InstanceOperationStart InstanceOperation = iota
	InstanceOperationStop
)

type InstanceFactory struct {
	limit int
}

var ErrInstanceLimitExceeded = errors.New("instance limit exceeded")

func NewInstanceFactory(limit int) *InstanceFactory {
	return &InstanceFactory{
		limit: limit,
	}
}

func (f *InstanceFactory) Create(teamID uuid.UUID, existing []Instance) (Instance, error) {
	if len(existing) >= f.limit {
		return Instance{}, ErrInstanceLimitExceeded
	}
	// use the next available index
	used := make(map[int]struct{}, len(existing))
	for _, instance := range existing {
		used[instance.Index] = struct{}{}
	}
	var index int
	for i := 1; i <= f.limit; i++ {
		if _, ok := used[i]; !ok {
			index = i
			break
		}
	}

	instance := Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  index,
		Infra:  InfraInstance{}, // initialize with empty InfraInstance
	}
	return instance, nil
}
