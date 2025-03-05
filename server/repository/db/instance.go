package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) FindInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error) {
	instance, err := models.FindInstance(ctx, r.executor(ctx), id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Instance{}, repository.ErrNotFound
		}
		return domain.Instance{}, err
	}

	return toDomainInstance(instance)
}

func toDomainInstance(instance *models.Instance) (domain.Instance, error) {
	id, err := uuid.Parse(instance.ID)
	if err != nil {
		return domain.Instance{}, err
	}
	teamID, err := uuid.Parse(instance.TeamID)
	if err != nil {
		return domain.Instance{}, err
	}
	status, err := toDomainInstanceStatus(instance.Status)
	if err != nil {
		return domain.Instance{}, err
	}

	return domain.Instance{
		ID:             id,
		TeamID:         teamID,
		InstanceNumber: int(instance.InstanceNumber),
		Status:         status,
		PrivateIP:      instance.PrivateIP.GetOrZero(),
		PublicIP:       instance.PublicIP.GetOrZero(),
	}, nil
}

func toDomainInstanceStatus(status models.InstancesStatus) (domain.InstanceStatus, error) {
	switch status {
	case models.InstancesStatusRunning:
		return domain.InstanceStatusRunning, nil
	case models.InstancesStatusBuilding:
		return domain.InstanceStatusBuilding, nil
	case models.InstancesStatusStarting:
		return domain.InstanceStatusStarting, nil
	case models.InstancesStatusStopping:
		return domain.InstanceStatusStopping, nil
	case models.InstancesStatusStopped:
		return domain.InstanceStatusStopped, nil
	case models.InstancesStatusDeleting:
		return domain.InstanceStatusDeleting, nil
	case models.InstancesStatusDeleted:
		return domain.InstanceStatusDeleted, nil
	default:
		return "", errors.New("unknown instance status")
	}
}

func fromDomainInstanceStatus(status domain.InstanceStatus) (models.InstancesStatus, error) {
	switch status {
	case domain.InstanceStatusRunning:
		return models.InstancesStatusRunning, nil
	case domain.InstanceStatusBuilding:
		return models.InstancesStatusBuilding, nil
	case domain.InstanceStatusStarting:
		return models.InstancesStatusStarting, nil
	case domain.InstanceStatusStopping:
		return models.InstancesStatusStopping, nil
	case domain.InstanceStatusStopped:
		return models.InstancesStatusStopped, nil
	case domain.InstanceStatusDeleting:
		return models.InstancesStatusDeleting, nil
	case domain.InstanceStatusDeleted:
		return models.InstancesStatusDeleted, nil
	default:
		return "", fmt.Errorf("unknown instance status: %v", status)
	}
}
