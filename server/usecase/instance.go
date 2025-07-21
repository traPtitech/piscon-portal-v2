package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance"
)

type InstanceUseCase interface {
	// GetInstance returns the instance with the given ID. If the instance is not found, [ErrNotFound] is returned.
	GetInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error)
	// GetTeamInstances returns all instances for the given team. Deleted instances are not included.
	GetTeamInstances(ctx context.Context, teamID uuid.UUID) ([]domain.Instance, error)
	// GetAllInstances returns all instances. Deleted instances are not included.
	GetAllInstances(ctx context.Context) ([]domain.Instance, error)
	// CreateInstance creates a new instance for the given team. If the team already has the maximum number of instances, [UseCaseError] is returned.
	CreateInstance(ctx context.Context, teamID uuid.UUID) (domain.Instance, error)
	// UpdateInstance updates the instance with the given ID. If the instance is not found, [ErrNotFound] is returned.
	// If the given operation is invalid, [UseCaseError] is returned.
	UpdateInstance(ctx context.Context, id uuid.UUID, op domain.InstanceOperation) error
	// DeleteInstance deletes the instance with the given ID. If the instance is not found, [ErrNotFound] is returned.
	// If the instance is already deleted, [UseCaseError] is returned.
	DeleteInstance(ctx context.Context, id uuid.UUID) error
}

type InstanceUseCaseImpl struct {
	repo    repository.Repository
	factory *domain.InstanceFactory
	manager instance.Manager
}

func NewInstanceUseCase(repo repository.Repository, factory *domain.InstanceFactory, manager instance.Manager) InstanceUseCase {
	return &InstanceUseCaseImpl{
		repo:    repo,
		factory: factory,
		manager: manager,
	}
}

func (i *InstanceUseCaseImpl) GetInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error) {
	instance, err := i.repo.FindInstance(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.Instance{}, ErrNotFound
	} else if err != nil {
		return domain.Instance{}, fmt.Errorf("find instance: %w", err)
	}

	infraInstance, err := i.manager.Get(ctx, instance.Infra.ProviderInstanceID)
	if err != nil {
		return domain.Instance{}, fmt.Errorf("get infra instance: %w", err)
	}

	instance.Infra = infraInstance

	return instance, nil
}

func (i *InstanceUseCaseImpl) CreateInstance(ctx context.Context, teamID uuid.UUID) (domain.Instance, error) {
	var instance domain.Instance
	var infraCreated bool
	err := i.repo.Transaction(ctx, func(ctx context.Context) error {
		existing, err := i.repo.GetTeamInstances(ctx, teamID)
		if err != nil {
			return fmt.Errorf("get team instances: %w", err)
		}

		instance, err = i.factory.Create(teamID, existing)
		if err != nil {
			return NewUseCaseError(err)
		}

		providerInstanceID, err := i.manager.Create(ctx, instanceName(teamID, instance.Index), nil) // TODO: pass team members' ssh keys
		if err != nil {
			return fmt.Errorf("create infra instance: %w", err)
		}
		instance.Infra.ProviderInstanceID = providerInstanceID
		infraCreated = true

		err = i.repo.CreateInstance(ctx, instance)
		if err != nil {
			return fmt.Errorf("create instance: %w", err)
		}

		return nil
	})
	if err != nil {
		if infraCreated {
			// rollback the created instance
			_ = i.manager.Delete(ctx, instance.Infra)
		}
		return domain.Instance{}, fmt.Errorf("transaction: %w", err)
	}

	return instance, nil
}

func (i *InstanceUseCaseImpl) GetTeamInstances(ctx context.Context, teamID uuid.UUID) ([]domain.Instance, error) {
	instances, err := i.repo.GetTeamInstances(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get team instances: %w", err)
	}

	infraInstanceIDs := make([]string, 0, len(instances))
	for _, inst := range instances {
		infraInstanceIDs = append(infraInstanceIDs, inst.Infra.ProviderInstanceID)
	}

	infraInstances, err := i.manager.GetByIDs(ctx, infraInstanceIDs)
	if err != nil {
		return nil, fmt.Errorf("get infra instances: %w", err)
	}

	result, err := setInfraInstances(instances, infraInstances)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *InstanceUseCaseImpl) GetAllInstances(ctx context.Context) ([]domain.Instance, error) {
	instances, err := i.repo.GetAllInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all instances: %w", err)
	}

	infraInstances, err := i.manager.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get infra instances: %w", err)
	}

	result, err := setInfraInstances(instances, infraInstances)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *InstanceUseCaseImpl) DeleteInstance(ctx context.Context, id uuid.UUID) error {
	err := i.repo.Transaction(ctx, func(ctx context.Context) error {
		instance, err := i.repo.FindInstance(ctx, id)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("find instance: %w", err)
		}

		err = i.manager.Delete(ctx, instance.Infra)
		if err != nil {
			return fmt.Errorf("delete infra instance: %w", err)
		}

		err = i.repo.DeleteInstance(ctx, id)
		if errors.Is(err, repository.ErrNotFound) {
			return NewUseCaseErrorFromMsg("instance not found or already deleted")
		}
		if err != nil {
			return fmt.Errorf("delete instance: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (i *InstanceUseCaseImpl) UpdateInstance(ctx context.Context, id uuid.UUID, op domain.InstanceOperation) error {
	err := i.repo.Transaction(ctx, func(ctx context.Context) error {
		instance, err := i.repo.FindInstance(ctx, id)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		} else if err != nil {
			return fmt.Errorf("find instance: %w", err)
		}

		switch op {
		case domain.InstanceOperationStart:
			err = i.manager.Start(ctx, instance.Infra)
		case domain.InstanceOperationStop:
			err = i.manager.Stop(ctx, instance.Infra)
		default:
			return NewUseCaseErrorFromMsg("invalid operation")
		}
		if err != nil {
			return fmt.Errorf("update infra instance: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func instanceName(teamID uuid.UUID, index int) string {
	return fmt.Sprintf("%s-%d", teamID.String(), index)
}

func setInfraInstances(instances []domain.Instance, infraInstances []domain.InfraInstance) ([]domain.Instance, error) {
	infraInstancesMap := make(map[string]domain.InfraInstance, len(infraInstances))
	for _, infra := range infraInstances {
		infraInstancesMap[infra.ProviderInstanceID] = infra
	}

	for i, inst := range instances {
		infra, ok := infraInstancesMap[inst.Infra.ProviderInstanceID]
		if !ok {
			return nil, fmt.Errorf("infra instance not found: %s", inst.Infra.ProviderInstanceID)
		}
		instances[i].Infra = infra
	}

	return instances, nil
}
