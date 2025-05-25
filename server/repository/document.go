package repository

import (
	"context"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type DocumentRepository interface {
	// GetDocument retrieves the document.
	// If the document does not exist, it returns ErrNotFound.
	GetDocument(ctx context.Context) (domain.Document, error)
}
