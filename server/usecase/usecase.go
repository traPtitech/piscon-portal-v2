package usecase

import (
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type UseCase interface {
	TeamUseCase
	UserUseCase
	BenchmarkUseCase
}

type useCaseImpl struct {
	TeamUseCase
	UserUseCase
	BenchmarkUseCase
}

func New(repo repository.Repository) UseCase {
	return &useCaseImpl{
		TeamUseCase:      NewTeamUseCase(repo),
		UserUseCase:      NewUserUseCase(repo),
		BenchmarkUseCase: NewBenchmarkUseCase(repo),
	}
}
