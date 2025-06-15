package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type DocumentUseCase interface {
	// GetDocument retrieves the document.
	// If the document does not exist, it returns ErrNotFound.
	GetDocument(ctx context.Context) (domain.Document, error)
	// CreateDocument creates a new document with the given body.
	CreateDocument(ctx context.Context, body string) (domain.Document, error)
}

type docUseCaseImpl struct {
	repo repository.Repository
}

func NewDocumentUseCase(repo repository.Repository) DocumentUseCase {
	return &docUseCaseImpl{repo: repo}
}

func (u *docUseCaseImpl) GetDocument(ctx context.Context) (domain.Document, error) {
	doc, err := u.repo.GetDocument(ctx)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.Document{}, ErrNotFound
	}
	if err != nil {
		return domain.Document{}, fmt.Errorf("get document: %w", err)
	}

	return doc, nil
}

func (u *docUseCaseImpl) CreateDocument(ctx context.Context, body string) (domain.Document, error) {
	docID := uuid.New()
	doc, err := u.repo.CreateDocument(ctx, docID, body)
	if err != nil {
		return domain.Document{}, fmt.Errorf("create document: %w", err)
	}

	return doc, nil
}
