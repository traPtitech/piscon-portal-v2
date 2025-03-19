package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type UserUseCase interface {
	GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
}

type userUseCaseImpl struct {
	repo repository.Repository
}

func NewUserUseCase(repo repository.Repository) UserUseCase {
	return &userUseCaseImpl{
		repo: repo,
	}
}

func (u *userUseCaseImpl) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := u.repo.FindUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}

func (u *userUseCaseImpl) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}
