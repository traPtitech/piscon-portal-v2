package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type repoDB struct {
	// prevent direct access to bob.DB by beginning with an underscore
	_db bob.DB
}

type executorCtxKeyT struct{}

var executorCtxKey = executorCtxKeyT{}

func (db *repoDB) executor(ctx context.Context) bob.Executor {
	if v := ctx.Value(executorCtxKey); v != nil {
		exe, ok := v.(bob.Executor)
		if ok {
			return exe
		}
	}

	return db._db
}

type Repository struct {
	*repoDB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		repoDB: &repoDB{
			_db: bob.NewDB(db),
		},
	}
}

func (r *Repository) Transaction(ctx context.Context, f func(ctx context.Context, r repository.Repository) error) error {
	tx, err := r._db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint errcheck

	ctx = context.WithValue(ctx, executorCtxKey, tx)

	err = f(ctx, r)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
