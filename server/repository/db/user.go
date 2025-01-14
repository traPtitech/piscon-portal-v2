package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) FindUser(ctx context.Context, id string) (domain.User, error) {
	return findUser(ctx, r.executor(ctx), id)
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

	return toDomainUser(user), nil
}

func createUser(ctx context.Context, executor bob.Executor, user domain.User) error {
	_, err := models.Users.Insert(&models.UserSetter{
		ID:   omit.From(user.ID),
		Name: omit.From(user.Name),
	}).Exec(ctx, executor)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func toDomainUser(user *models.User) domain.User {
	return domain.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
