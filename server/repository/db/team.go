package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	team, err := models.Teams.Query(models.SelectWhere.Teams.ID.EQ(id.String()), models.SelectThenLoad.Team.Users()).One(ctx, r.executor(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, repository.ErrNotFound
		}
		return domain.Team{}, fmt.Errorf("find team: %w", err)
	}
	return toDomainTeam(team)
}

func (r *Repository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	teams, err := models.Teams.Query(models.SelectThenLoad.Team.Users()).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get teams: %w", err)
	}
	res := make([]domain.Team, 0, len(teams))
	for _, team := range teams {
		domainTeam, err := toDomainTeam(team)
		if err != nil {
			return nil, fmt.Errorf("to domain team: %w", err)
		}
		res = append(res, domainTeam)
	}
	return res, nil
}

func (r *Repository) CreateTeam(ctx context.Context, team domain.Team) error {
	if ctx.Value(executorCtxKey) == nil {
		return r.Transaction(ctx, func(ctx context.Context) error {
			return r.CreateTeam(ctx, team)
		})
	}

	_, err := models.Teams.Insert(&models.TeamSetter{
		ID:        lo.ToPtr(team.ID.String()),
		Name:      lo.ToPtr(team.Name),
		CreatedAt: lo.ToPtr(team.CreatedAt),
	}).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("create team: %w", err)
	}

	memberIDs := lo.Map(team.Members, func(m domain.User, _ int) string { return m.ID.String() })
	_, err = models.Users.Update(
		models.UpdateWhere.Users.ID.In(memberIDs...),
		models.UserSetter{
			TeamID: ToSQLNull(team.ID.String()),
		}.UpdateMod(),
	).Exec(ctx, r.executor(ctx))

	if err != nil {
		return fmt.Errorf("update users: %w", err)
	}

	return nil
}

func (r *Repository) UpdateTeam(ctx context.Context, team domain.Team) error {
	if ctx.Value(executorCtxKey) == nil {
		return r.Transaction(ctx, func(ctx context.Context) error {
			return r.UpdateTeam(ctx, team)
		})
	}

	_, err := models.Teams.Update(
		models.UpdateWhere.Teams.ID.EQ(team.ID.String()),
		models.TeamSetter{
			Name: lo.ToPtr(team.Name),
		}.UpdateMod(),
	).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("update team: %w", err)
	}

	memberIDs := lo.Map(team.Members, func(m domain.User, _ int) string { return m.ID.String() })
	_, err = models.Users.Update(
		models.UpdateWhere.Users.ID.In(memberIDs...),
		models.UserSetter{
			TeamID: ToSQLNull(team.ID.String()),
		}.UpdateMod(),
	).Exec(ctx, r.executor(ctx))

	if err != nil {
		return fmt.Errorf("update users: %w", err)
	}

	return nil
}

func toDomainTeam(team *models.Team) (domain.Team, error) {
	members := make([]domain.User, 0, len(team.R.Users))
	for _, user := range team.R.Users {
		domainUser, err := toDomainUser(user)
		if err != nil {
			return domain.Team{}, fmt.Errorf("to domain user: %w", err)
		}
		members = append(members, domainUser)
	}
	teamID, err := uuid.Parse(team.ID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("parse team ID: %w", err)
	}
	return domain.Team{
		ID:        teamID,
		Name:      team.Name,
		Members:   members,
		CreatedAt: team.CreatedAt,
	}, nil
}
