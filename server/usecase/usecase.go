package usecase

import (
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type UseCase interface {
	TeamUseCase
	UserUseCase
}

type useCaseImpl struct {
	TeamUseCase
	UserUseCase
}

func New(repo repository.Repository) UseCase {
	return &useCaseImpl{
		TeamUseCase: NewTeamUseCase(repo),
		UserUseCase: NewUserUseCase(repo),
	}
}
