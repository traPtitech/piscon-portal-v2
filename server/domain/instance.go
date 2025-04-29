package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Instance struct {
	ID     uuid.UUID
	TeamID uuid.UUID
	Index  int
	Status InstanceStatus
	Infra  InfraInstance
}

type InfraInstance struct {
	ProviderInstanceID string
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
	used := make(map[int]struct{})
	for _, instance := range existing {
		used[instance.Index] = struct{}{}
	}
	var index int
	for i := range f.limit {
		if _, ok := used[i+1]; !ok {
			index = i + 1
			break
		}
	}

	instance := Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  index,
		Status: InstanceStatusStarting,
		Infra:  InfraInstance{}, // initialize with empty InfraInstance
	}
	return instance, nil
}
