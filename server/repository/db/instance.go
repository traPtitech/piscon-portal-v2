package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

var whereNotDeleted = models.SelectWhere.Instances.Status.NE(models.InstancesStatusDeleted)

func (r *Repository) CreateInstance(ctx context.Context, instance domain.Instance) error {
	setter, err := buildInstanceSetter(instance)
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	_, err = models.Instances.Insert(setter).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	return nil
}

func (r *Repository) UpdateInstance(ctx context.Context, instance domain.Instance) error {
	setter, err := buildInstanceSetter(instance)
	if err != nil {
		return fmt.Errorf("update instance: %w", err)
	}
	modelInstance, err := models.FindInstance(ctx, r.executor(ctx), instance.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("find instance for update: %w", err)
	}
	err = modelInstance.Update(ctx, r.executor(ctx), setter)
	if err != nil {
		return fmt.Errorf("update instance: %w", err)
	}
	return nil
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
	instances, err := models.Instances.Query(
		models.SelectWhere.Instances.TeamID.EQ(teamID.String()),
		whereNotDeleted, // Exclude deleted instances
	).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get team instances: %w", err)
	}
	result := make([]domain.Instance, 0, len(instances))
	for _, inst := range instances {
		domainInst, err := toDomainInstance(inst)
		if err != nil {
			return nil, err
		}
		result = append(result, domainInst)
	}
	return result, nil
}

func (r *Repository) GetAllInstances(ctx context.Context) ([]domain.Instance, error) {
	instances, err := models.Instances.Query(whereNotDeleted).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get all instances: %w", err)
	}
	result := make([]domain.Instance, 0, len(instances))
	for _, inst := range instances {
		domainInst, err := toDomainInstance(inst)
		if err != nil {
			return nil, err
		}
		result = append(result, domainInst)
	}
	return result, nil
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
			ProviderInstanceID: instance.ProviderInstanceID,
			Status:             status,
			PrivateIP:          instance.PrivateIP.V,
			PublicIP:           instance.PublicIP.V,
		},
		CreatedAt: instance.CreatedAt,
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

func buildInstanceSetter(instance domain.Instance) (*models.InstanceSetter, error) {
	status, err := fromDomainInstanceStatus(instance.Infra.Status)
	if err != nil {
		return nil, err
	}
	setter := &models.InstanceSetter{
		ID:                 lo.ToPtr(instance.ID.String()),
		ProviderInstanceID: lo.ToPtr(instance.Infra.ProviderInstanceID),
		TeamID:             lo.ToPtr(instance.TeamID.String()),
		InstanceNumber:     lo.ToPtr(int32(instance.Index)),
		Status:             &status,
		PublicIP:           ToSQLNull(instance.Infra.PublicIP),
		PrivateIP:          ToSQLNull(instance.Infra.PrivateIP),
		CreatedAt:          lo.Ternary(!instance.CreatedAt.IsZero(), &instance.CreatedAt, nil),
	}
	return setter, nil
}
