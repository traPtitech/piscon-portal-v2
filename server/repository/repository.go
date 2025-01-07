package repository

import (
	"context"
	"errors"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type Repository interface {
	// Transaction starts a transaction and calls f with the transaction.
	// If f returns an error, the transaction is rolled back and the error is returned.
	Transaction(ctx context.Context, f func(ctx context.Context, r Repository) error) error

	// FindUser finds a user by id. If the user is not found, it returns [ErrNotFound].
	FindUser(ctx context.Context, id string) (domain.User, error)
	// CreateUser creates a user.
	CreateUser(ctx context.Context, user domain.User) error

	// FindSession finds a session by id. If the session is not found, it returns [ErrNotFound].
	FindSession(ctx context.Context, id string) (domain.Session, error)
	// CreateSession creates a session.
	CreateSession(ctx context.Context, session domain.Session) error
	// DeleteSession deletes a session by id.
	DeleteSession(ctx context.Context, id string) error
}

var (
	ErrNotFound = errors.New("not found")
)
