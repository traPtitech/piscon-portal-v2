package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	return findTeam(ctx, r.db, id.String())
}

func (t *txRepository) FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	return findTeam(ctx, t.tx, id.String())
}

func (r *Repository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return getTeams(ctx, r.db)
}

func (t *txRepository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return getTeams(ctx, t.tx)
}

func (r *Repository) CreateTeam(ctx context.Context, team domain.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint errcheck

	if err := createTeam(ctx, tx, team); err != nil {
		return fmt.Errorf("create team: %w", err)
	}

	return tx.Commit()
}

func (t *txRepository) CreateTeam(ctx context.Context, team domain.Team) error {
	return createTeam(ctx, t.tx, team)
}

func (r *Repository) UpdateTeam(ctx context.Context, team domain.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() //nolint errcheck

	if err := updateTeam(ctx, tx, team); err != nil {
		return fmt.Errorf("update team: %w", err)
	}

	return tx.Commit()
}

func (t *txRepository) UpdateTeam(ctx context.Context, team domain.Team) error {
	return updateTeam(ctx, t.tx, team)
}

func findTeam(ctx context.Context, executor bob.Executor, id string) (domain.Team, error) {
	team, err := models.Teams.Query(models.SelectWhere.Teams.ID.EQ(id), models.ThenLoadTeamUsers()).One(ctx, executor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, repository.ErrNotFound
		}
		return domain.Team{}, fmt.Errorf("find team: %w", err)
	}
	return toDomainTeam(team)
}

func getTeams(ctx context.Context, executor bob.Executor) ([]domain.Team, error) {
	teams, err := models.Teams.Query(models.ThenLoadTeamUsers()).All(ctx, executor)
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

func createTeam(ctx context.Context, tx bob.Tx, team domain.Team) error {
	_, err := models.Teams.Insert(&models.TeamSetter{
		ID:        omit.From(team.ID.String()),
		Name:      omit.From(team.Name),
		CreatedAt: omit.From(team.CreatedAt),
	}).Exec(ctx, tx)
	if err != nil {
		return fmt.Errorf("create team: %w", err)
	}

	memberIDs := lo.Map(team.Members, func(m domain.User, _ int) string { return m.ID.String() })
	_, err = models.Users.Update(
		models.UpdateWhere.Users.ID.In(memberIDs...),
		models.UserSetter{
			TeamID: omitnull.From(team.ID.String()),
		}.UpdateMod(),
	).Exec(ctx, tx)

	return err
}

func updateTeam(ctx context.Context, tx bob.Tx, team domain.Team) error {
	_, err := models.Teams.Update(
		models.UpdateWhere.Teams.ID.EQ(team.ID.String()),
		models.TeamSetter{
			Name: omit.From(team.Name),
		}.UpdateMod(),
	).Exec(ctx, tx)
	if err != nil {
		return fmt.Errorf("update team: %w", err)
	}

	memberIDs := lo.Map(team.Members, func(m domain.User, _ int) string { return m.ID.String() })
	_, err = models.Users.Update(
		models.UpdateWhere.Users.ID.In(memberIDs...),
		models.UserSetter{
			TeamID: omitnull.From(team.ID.String()),
		}.UpdateMod(),
	).Exec(ctx, tx)

	return err
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
