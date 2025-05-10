package instance

import (
	"context"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed=true

// Manager is an interface for managing infrastructure instances.
type Manager interface {
	Create(ctx context.Context, name string, sshPubKeys []string) (domain.InfraInstance, error)
	Get(ctx context.Context, id string) (domain.InfraInstance, error)
	GetAll(ctx context.Context) ([]domain.InfraInstance, error)
	Delete(ctx context.Context, instance domain.InfraInstance) (domain.InfraInstance, error)
	Stop(ctx context.Context, instance domain.InfraInstance) (domain.InfraInstance, error)
	Start(ctx context.Context, instance domain.InfraInstance) (domain.InfraInstance, error)
}
