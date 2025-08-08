package usecase

import (
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type UseCase interface {
	TeamUseCase
	UserUseCase
	BenchmarkUseCase
	ScoreUseCase
	AdminUseCase
	DocumentUseCase
	InstanceUseCase
}

type useCaseImpl struct {
	TeamUseCase
	UserUseCase
	BenchmarkUseCase
	ScoreUseCase
	AdminUseCase
	DocumentUseCase
	InstanceUseCase
}

type Config struct {
	InstanceLimit int
}

func New(config Config, repo repository.Repository, instanceManager instance.Manager) UseCase {
	return &useCaseImpl{
		TeamUseCase:      NewTeamUseCase(repo),
		UserUseCase:      NewUserUseCase(repo),
		BenchmarkUseCase: NewBenchmarkUseCase(repo, instanceManager),
		ScoreUseCase:     NewScoreUseCase(repo),
		AdminUseCase:     NewAdminUseCase(repo),
		DocumentUseCase:  NewDocumentUseCase(repo),
		InstanceUseCase:  NewInstanceUseCase(repo, domain.NewInstanceFactory(config.InstanceLimit), instanceManager),
	}
}
