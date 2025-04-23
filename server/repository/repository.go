package repository

import (
	"context"
	"errors"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true
type Repository interface {
	// Transaction starts a transaction and calls f with the transaction.
	// If f returns an error, the transaction is rolled back and the error is returned.
	Transaction(ctx context.Context, f func(ctx context.Context) error) error

	UserRepository
	SessionRepository
	TeamRepository
	BenchmarkRepository
	InstanceRepository
}

var (
	ErrNotFound = errors.New("not found")
)
