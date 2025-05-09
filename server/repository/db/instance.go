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

func (r *Repository) CreateInstance(ctx context.Context, instance domain.Instance) error {
	return errors.New("not implemented")
}

func (r *Repository) UpdateInstance(ctx context.Context, instance domain.Instance) error {
	return errors.New("not implemented")
}

func (r *Repository) FindInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error) {
	instance, err := models.FindInstance(ctx, r.executor(ctx), id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Instance{}, repository.ErrNotFound
		}
		return domain.Instance{}, fmt.Errorf("find instance: %w", err)
	}

	return toDomainInstance(instance)
}

func (r *Repository) GetTeamInstances(ctx context.Context, teamID uuid.UUID) ([]domain.Instance, error) {
	return nil, errors.New("not implemented")
}

func (r *Repository) GetAllInstances(ctx context.Context) ([]domain.Instance, error) {
	return nil, errors.New("not implemented")
}

func toDomainInstance(instance *models.Instance) (domain.Instance, error) {
	id, err := uuid.Parse(instance.ID)
	if err != nil {
		return domain.Instance{}, fmt.Errorf("parse instance id: %w", err)
	}
	teamID, err := uuid.Parse(instance.TeamID)
	if err != nil {
		return domain.Instance{}, fmt.Errorf("parse team id: %w", err)
	}
	status, err := toDomainInstanceStatus(instance.Status)
	if err != nil {
		return domain.Instance{}, fmt.Errorf("parse instance status: %w", err)
	}

	return domain.Instance{
		ID:     id,
		TeamID: teamID,
		Index:  int(instance.InstanceNumber),
		Infra: domain.InfraInstance{
			Status:    status,
			PrivateIP: instance.PrivateIP.GetOrZero(),
			PublicIP:  instance.PublicIP.GetOrZero(),
		},
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
