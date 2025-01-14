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

func (r *Repository) FindSession(ctx context.Context, id string) (domain.Session, error) {
	return findSession(ctx, r.executor(ctx), id)
}

func (r *Repository) CreateSession(ctx context.Context, session domain.Session) error {
	return createSession(ctx, r.executor(ctx), session)
}

func (r *Repository) DeleteSession(ctx context.Context, id string) error {
	return deleteSession(ctx, r.executor(ctx), id)
}

func findSession(ctx context.Context, executor bob.Executor, id string) (domain.Session, error) {
	session, err := models.FindSession(ctx, executor, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Session{}, repository.ErrNotFound
		}
		return domain.Session{}, fmt.Errorf("find session: %w", err)
	}

	return toDomainSession(session), nil
}

func createSession(ctx context.Context, executor bob.Executor, session domain.Session) error {
	_, err := models.Sessions.Insert(&models.SessionSetter{
		ID:        omit.From(session.ID),
		UserID:    omit.From(session.UserID),
		ExpiredAt: omit.From(session.ExpiresAt),
	}).Exec(ctx, executor)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func deleteSession(ctx context.Context, executor bob.Executor, id string) error {
	_, err := models.Sessions.Delete(models.DeleteWhere.Sessions.ID.EQ(id)).Exec(ctx, executor)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func toDomainSession(session *models.Session) domain.Session {
	return domain.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiredAt,
	}
}
