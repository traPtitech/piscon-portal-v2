package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/traq"
)

type UserUseCase interface {
	GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
}

type userUseCaseImpl struct {
	repo        repository.Repository
	traqService traq.Service
}

func NewUserUseCase(repo repository.Repository, traqService traq.Service) *userUseCaseImpl {
	return &userUseCaseImpl{
		repo:        repo,
		traqService: traqService,
	}
}

func (u *userUseCaseImpl) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := u.repo.FindUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (u *userUseCaseImpl) GetUsers(ctx context.Context) ([]domain.User, error) {
	traqUsers, err := u.traqService.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	portalUsers, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	portalUserIDs := lo.SliceToMap(portalUsers, func(u domain.User) (uuid.UUID, struct{}) { return u.ID, struct{}{} })

	// merge users
	// if a user exists in both traq and portal, use the user in portal
	users := make([]domain.User, 0, len(traqUsers))
	for _, traqUser := range traqUsers {
		// check if the user exists in portal
		_, exists := portalUserIDs[traqUser.ID]
		if !exists {
			users = append(users, domain.User{
				ID:   traqUser.ID,
				Name: traqUser.Name,
			})
		}
	}
	users = append(users, portalUsers...)

	return users, nil
}
