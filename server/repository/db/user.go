package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
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
