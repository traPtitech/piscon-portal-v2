package instance

import (
	"context"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true

// Manager is an interface for managing infrastructure instances.
type Manager interface {
	// Create creates a new infrastructure instance with the given name and SSH public keys.
	// Returns the provider instance ID or an error if the creation fails.
	Create(ctx context.Context, name string, sshPubKeys []string) (string, error)
	// Get retrieves an infrastructure instance by its provider ID.
	// Returns the instance or an error if the retrieval fails.
	Get(ctx context.Context, id string) (domain.InfraInstance, error)
	// GetAll retrieves all infrastructure instances.
	// Returns a slice of instances or an error if the retrieval fails.
	GetAll(ctx context.Context) ([]domain.InfraInstance, error)
	// GetByIDs retrieves infrastructure instances by their provider instance IDs.
	// Returns a slice of instances or an error if the retrieval fails.
	GetByIDs(ctx context.Context, ids []string) ([]domain.InfraInstance, error)
	// Update updates the infrastructure instance with the given provider ID.
	// Returns the updated instance or an error if the update fails.
	Delete(ctx context.Context, instance domain.InfraInstance) error
	// Stop stops the infrastructure instance with the given provider ID.
	// Returns an error if the stop operation fails.
	Stop(ctx context.Context, instance domain.InfraInstance) error
	// Start starts the infrastructure instance with the given provider ID.
	// Returns an error if the start operation fails.
	Start(ctx context.Context, instance domain.InfraInstance) error
}
