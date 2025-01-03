package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type TeamUseCase struct {
	repo repository.Repository
}

func NewTeamUseCase(repo repository.Repository) *TeamUseCase {
	return &TeamUseCase{
		repo: repo,
	}
}

func (u *TeamUseCase) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return u.repo.GetTeams(ctx)
}

func (u *TeamUseCase) GetTeam(ctx context.Context, id string) (domain.Team, error) {
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
	MemberIDs []string
	// CreatorID is the ID of the user who creates the team
	CreatorID string
}

func (u *TeamUseCase) CreateTeam(ctx context.Context, input CreateTeamInput) (domain.Team, error) {
	// creator must be a member of the team
	isMember := slices.Contains(input.MemberIDs, input.CreatorID)
	if !isMember {
		return domain.Team{}, NewErrBadRequest("creator must be a member of the team")
	}

	// the number of members must be less than or equal to 3
	if len(input.MemberIDs) > 3 {
		return domain.Team{}, NewErrBadRequest("the number of members must be less than or equal to 3")
	}

	team := domain.NewTeam(input.Name)
	for _, memberID := range input.MemberIDs {
		user, err := u.repo.FindUser(ctx, memberID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return domain.Team{}, NewErrBadRequest("member not found")
			}
			return domain.Team{}, err
		}
		if err := team.AddMember(user); err != nil {
			return domain.Team{}, NewErrBadRequestFromErr(err)
		}
	}

	if err := u.repo.CreateTeam(ctx, team); err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

type UpdateTeamInput struct {
	ID        string
	Name      string
	MemberIDs []string
}

func (u *TeamUseCase) UpdateTeam(ctx context.Context, input UpdateTeamInput) (domain.Team, error) {
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
