package db

import (
	"context"
	"database/sql"

	"github.com/stephenafamo/bob"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type Repository struct {
	db bob.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: bob.NewDB(db),
	}
}

type txRepository struct {
	tx bob.Tx
}

func newTxRepository(tx bob.Tx) *txRepository {
	return &txRepository{
		tx: tx,
	}
}

func (r *Repository) Transaction(ctx context.Context, f func(ctx context.Context, r repository.Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint errcheck

	txRepo := newTxRepository(tx)

	err = f(ctx, txRepo)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *txRepository) Transaction(ctx context.Context, f func(ctx context.Context, r repository.Repository) error) error {
	return f(ctx, t)
}
