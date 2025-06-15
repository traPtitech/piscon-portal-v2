package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type DocumentRepository interface {
	// GetDocument retrieves the document.
	// If the document does not exist, it returns ErrNotFound.
	GetDocument(ctx context.Context) (domain.Document, error)
	// CreateDocument creates a new document with the given ID and body.
	CreateDocument(ctx context.Context, id uuid.UUID, body string) (domain.Document, error)
}
