package domain

import (
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
	n := maxInstanceNumber(existing)
	return Instance{
		ID:             uuid.New(),
		TeamID:         teamID,
		InstanceNumber: n + 1,
		Status:         InstanceStatusStopping,
	}
}

func maxInstanceNumber(instances Instances) int {
	if len(instances) < 1 {
		return 0
	}
	m := instances[0].InstanceNumber
	for _, i := range instances {
		if i.InstanceNumber > m {
			m = i.InstanceNumber
		}
	}
	return m
}
