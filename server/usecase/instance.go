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

		infra, err := i.manager.Create(ctx, instanceName(teamID, instance.Index), nil) // TODO: pass team members' ssh keys
		if err != nil {
			return fmt.Errorf("create infra instance: %w", err)
		}
		instance.Infra = infra
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
			_, _ = i.manager.Delete(ctx, instance.Infra)
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
	return instances, nil
}

func (i *InstanceUseCaseImpl) GetAllInstances(ctx context.Context) ([]domain.Instance, error) {
	instances, err := i.repo.GetAllInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all instances: %w", err)
	}
	return instances, nil
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

		deletedInstance, err := i.manager.Delete(ctx, instance.Infra)
		if err != nil {
			return fmt.Errorf("delete infra instance: %w", err)
		}
		instance.Infra = deletedInstance

		return i.repo.UpdateInstance(ctx, instance)
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

		var updatedInstance domain.InfraInstance
		switch op {
		case domain.InstanceOperationStart:
			updatedInstance, err = i.manager.Start(ctx, instance.Infra)
		case domain.InstanceOperationStop:
			updatedInstance, err = i.manager.Stop(ctx, instance.Infra)
		default:
			return NewUseCaseErrorFromMsg("invalid operation")
		}
		if err != nil {
			return fmt.Errorf("update infra instance: %w", err)
		}
		instance.Infra = updatedInstance

		return i.repo.UpdateInstance(ctx, instance)
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func instanceName(teamID uuid.UUID, index int) string {
	return fmt.Sprintf("%s-%d", teamID.String(), index)
}
