package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

type DocumentUseCase interface {
	// GetDocument retrieves the document.
	// If the document does not exist, it returns ErrNotFound.
	GetDocument(ctx context.Context) (domain.Document, error)
}

type docUseCaseImpl struct {
	repo repository.Repository
}

func NewDocUseCase(repo repository.Repository) DocumentUseCase {
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
