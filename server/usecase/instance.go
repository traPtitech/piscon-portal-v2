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
	GetTeamInstances(ctx context.Context, teamID uuid.UUID) ([]domain.Instance, error)
	GetAllInstances(ctx context.Context) ([]domain.Instance, error)
	CreateInstance(ctx context.Context, teamID uuid.UUID) (domain.Instance, error)
	UpdateInstance(ctx context.Context, id uuid.UUID, status domain.InstanceOperation) error
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
	return i.repo.GetTeamInstances(ctx, teamID)
}

func (i *InstanceUseCaseImpl) GetAllInstances(ctx context.Context) ([]domain.Instance, error) {
	return i.repo.GetAllInstances(ctx)
}

func (i *InstanceUseCaseImpl) DeleteInstance(ctx context.Context, id uuid.UUID) error {
	return errors.New("TODO: not implemented")
}

func (i *InstanceUseCaseImpl) UpdateInstance(ctx context.Context, id uuid.UUID, op domain.InstanceOperation) error {
	return errors.New("TODO: not implemented")
}

func instanceName(teamID uuid.UUID, index int) string {
	return fmt.Sprintf("%s-%d", teamID.String(), index)
}
