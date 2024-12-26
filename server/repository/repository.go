package repository

import (
	"context"
	"errors"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type Repository interface {
	Transaction(ctx context.Context, f func(ctx context.Context, r Repository) error) error

	FindUser(ctx context.Context, id string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) error

	FindSession(ctx context.Context, id string) (domain.Session, error)
	CreateSession(ctx context.Context, session domain.Session) error
	DeleteSession(ctx context.Context, id string) error
}

var (
	ErrNotFound = errors.New("not found")
)
