package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) FindUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return findUser(ctx, r.executor(ctx), id.String())
}

func (r *Repository) CreateUser(ctx context.Context, user domain.User) error {
	return createUser(ctx, r.executor(ctx), user)
}

func (r *Repository) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := models.Users.Query().All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	res := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUser, err := toDomainUser(user)
		if err != nil {
			return nil, fmt.Errorf("convert user: %w", err)
		}
		res = append(res, domainUser)
	}

	return res, nil
}

func findUser(ctx context.Context, executor bob.Executor, id string) (domain.User, error) {
	user, err := models.FindUser(ctx, executor, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}

	return toDomainUser(user)
}

func createUser(ctx context.Context, executor bob.Executor, user domain.User) error {
	_, err := models.Users.Insert(&models.UserSetter{
		ID:   omit.From(user.ID.String()),
		Name: omit.From(user.Name),
	}).Exec(ctx, executor)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *Repository) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.User, error) {
	if len(ids) == 0 {
		return []domain.User{}, nil
	}

	anyIDs := make([]any, len(ids))
	for i, id := range ids {
		anyIDs[i] = id.String()
	}
	users, err := models.Users.Query(
		sm.Where(models.UserColumns.ID.In(mysql.Arg(anyIDs...))),
	).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get users by ids: %w", err)
	}
	res := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUser, err := toDomainUser(user)
		if err != nil {
			return nil, fmt.Errorf("convert user: %w", err)
		}
		res = append(res, domainUser)
	}

	return res, nil
}

func (r *Repository) GetAdmins(ctx context.Context) ([]domain.User, error) {
	adminUsers, err := models.Users.Query(
		sm.Where(models.UserColumns.IsAdmin.EQ(mysql.Arg(true))),
	).All(ctx, r.executor(ctx))
	if err != nil {
		return nil, fmt.Errorf("get admins: %w", err)
	}

	res := make([]domain.User, 0, len(adminUsers))
	for _, user := range adminUsers {
		domainUser, err := toDomainUser(user)
		if err != nil {
			return nil, fmt.Errorf("convert user: %w", err)
		}
		res = append(res, domainUser)
	}

	return res, nil
}

func (r *Repository) AddAdmins(ctx context.Context, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	anyIDs := make([]any, len(userIDs))
	for i, id := range userIDs {
		anyIDs[i] = id.String()
	}
	_, err := models.Users.Update(
		um.SetCol(models.ColumnNames.Users.IsAdmin).To(true),
		um.Where(models.UserColumns.ID.In(mysql.Arg(anyIDs...))),
	).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("add admins: %w", err)
	}

	return nil
}

func (r *Repository) DeleteAdmins(ctx context.Context, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	anyIDs := make([]any, len(userIDs))
	for i, id := range userIDs {
		anyIDs[i] = id.String()
	}
	_, err := models.Users.Update(
		um.SetCol(models.ColumnNames.Users.IsAdmin).To(false),
		um.Where(models.UserColumns.ID.In(mysql.Arg(anyIDs...))),
	).Exec(ctx, r.executor(ctx))
	if err != nil {
		return fmt.Errorf("delete admins: %w", err)
	}

	return nil
}

func toDomainUser(user *models.User) (domain.User, error) {
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("parse user ID: %w", err)
	}

	var teamID uuid.NullUUID
	if id, ok := user.TeamID.Get(); ok {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return domain.User{}, fmt.Errorf("parse team ID: %w", err)
		}
		teamID = uuid.NullUUID{
			UUID:  parsedID,
			Valid: true,
		}
	}

	return domain.User{
		ID:      userID,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
		TeamID:  teamID,
	}, nil
}
