package github

import (
	"context"
	"errors"
	"fmt"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true

// Service provides access to GitHub API functionality.
type Service interface {
	// GetSSHKeys retrieves SSH public keys for the given GitHub usernames.
	// Returns an error if any of the users don't exist or if there's an API error.
	GetSSHKeys(ctx context.Context, githubIDs []string) ([]string, error)
}

// UserNotFoundError contains information about which user was not found.
type UserNotFoundError struct {
	Username string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("github user not found: %s", e.Username)
}

func IsNotFound(err error) bool {
	var userNotFoundErr *UserNotFoundError
	return errors.As(err, &userNotFoundErr)
}

// NewUserNotFoundError creates a new UserNotFoundError.
func NewUserNotFoundError(username string) *UserNotFoundError {
	return &UserNotFoundError{Username: username}
}
