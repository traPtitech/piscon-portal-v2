package usecase

import (
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/traq"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type UseCase interface {
	TeamUseCase
	UserUseCase
}

type useCaseImpl struct {
	*teamUseCaseImpl
	*userUseCaseImpl
}

func New(repo repository.Repository, traqService traq.Service) UseCase {
	return &useCaseImpl{
		teamUseCaseImpl: NewTeamUseCase(repo),
		userUseCaseImpl: NewUserUseCase(repo, traqService),
	}
}
