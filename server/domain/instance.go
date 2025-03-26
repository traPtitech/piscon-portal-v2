package domain

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
)

type Instance struct {
	ID             uuid.UUID
	TeamID         uuid.UUID
	InstanceNumber int
	Status         InstanceStatus
	PrivateIP      string
	PublicIP       string
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

type Instances []Instance

func NewInstance(teamID uuid.UUID, existing Instances) Instance {
	m := slices.MaxFunc(existing, func(a, b Instance) int { return cmp.Compare(a.InstanceNumber, b.InstanceNumber) })
	return Instance{
		ID:             uuid.New(),
		TeamID:         teamID,
		InstanceNumber: m.InstanceNumber + 1,
		Status:         InstanceStatusStopping,
	}
}
