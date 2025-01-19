package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type TeamUseCase interface {
	GetTeams(ctx context.Context) ([]domain.Team, error)
	GetTeam(ctx context.Context, id uuid.UUID) (domain.Team, error)
	CreateTeam(ctx context.Context, input CreateTeamInput) (domain.Team, error)
	UpdateTeam(ctx context.Context, input UpdateTeamInput) (domain.Team, error)
}

type teamUseCaseImpl struct {
	repo repository.Repository
}

func NewTeamUseCase(repo repository.Repository) *teamUseCaseImpl {
	return &teamUseCaseImpl{
		repo: repo,
	}
}

func (u *teamUseCaseImpl) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return u.repo.GetTeams(ctx)
}

func (u *teamUseCaseImpl) GetTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	team, err := u.repo.FindTeam(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Team{}, ErrNotFound
		}
		return domain.Team{}, err
	}
	return team, nil
}

type CreateTeamInput struct {
	Name      string
	MemberIDs []uuid.UUID
	// CreatorID is the ID of the user who creates the team
	CreatorID uuid.UUID
}

func (u *teamUseCaseImpl) CreateTeam(ctx context.Context, input CreateTeamInput) (domain.Team, error) {
	// creator must be a member of the team
	isMember := slices.Contains(input.MemberIDs, input.CreatorID)
	if !isMember {
		return domain.Team{}, NewUseCaseErrorFromMsg("creator must be a member of the team")
	}

	team := domain.NewTeam(input.Name)
	for _, memberID := range input.MemberIDs {
		user, err := u.repo.FindUser(ctx, memberID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return domain.Team{}, NewUseCaseErrorFromMsg("member not found")
			}
			return domain.Team{}, err
		}
		if err := team.AddMember(user); err != nil {
			return domain.Team{}, NewUseCaseError(err)
		}
	}

	if err := u.repo.CreateTeam(ctx, team); err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

type UpdateTeamInput struct {
	ID        uuid.UUID
	Name      string
	MemberIDs []uuid.UUID
}

func (u *teamUseCaseImpl) UpdateTeam(ctx context.Context, input UpdateTeamInput) (domain.Team, error) {
	team, err := u.repo.FindTeam(ctx, input.ID)
	if err != nil {
		return domain.Team{}, err
	}

	if input.Name != "" {
		team.Name = input.Name
	}

	for _, memberID := range input.MemberIDs {
		user, err := u.repo.FindUser(ctx, memberID)
		if err != nil {
			return domain.Team{}, err
		}
		if err := team.AddMember(user); err != nil {
			return domain.Team{}, err
		}
	}

	err = u.repo.UpdateTeam(ctx, team)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}
