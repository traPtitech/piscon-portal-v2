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

var whereNotDeleted = models.SelectWhere.Instances.DeletedAt.IsNull()

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

func (r *Repository) DeleteInstance(ctx context.Context, id uuid.UUID) error {
	rows, err := models.Instances.Delete(models.DeleteWhere.Instances.ID.EQ(id.String())).
		Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("delete instance: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
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

	return domain.Instance{
		ID:     id,
		TeamID: teamID,
		Index:  int(instance.InstanceNumber),
		Infra: domain.InfraInstance{
			ProviderInstanceID: instance.ProviderInstanceID,
		},
		CreatedAt: instance.CreatedAt,
	}, nil
}

func buildInstanceSetter(instance domain.Instance) (*models.InstanceSetter, error) {
	setter := &models.InstanceSetter{
		ID:                 lo.ToPtr(instance.ID.String()),
		ProviderInstanceID: lo.ToPtr(instance.Infra.ProviderInstanceID),
		TeamID:             lo.ToPtr(instance.TeamID.String()),
		InstanceNumber:     lo.ToPtr(int32(instance.Index)),
		CreatedAt:          lo.Ternary(!instance.CreatedAt.IsZero(), &instance.CreatedAt, nil),
	}
	return setter, nil
}
