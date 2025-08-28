package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func (r *Repository) GetDocument(ctx context.Context) (domain.Document, error) {
	doc, err := models.Documents.Query(sm.OrderBy(models.Documents.Columns.CreatedAt).Desc()).One(ctx, r.executor(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Document{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Document{}, fmt.Errorf("get document: %w", err)
	}

	domainDoc, err := toDomainDocument(doc)
	if err != nil {
		return domain.Document{}, err
	}
	return domainDoc, nil
}

func (r *Repository) CreateDocument(ctx context.Context, id uuid.UUID, body string) (domain.Document, error) {
	_, err := models.Documents.Insert(&models.DocumentSetter{
		ID:   lo.ToPtr(id.String()),
		Body: lo.ToPtr(body),
	}).Exec(ctx, r.executor(ctx))
	if err != nil {
		return domain.Document{}, fmt.Errorf("insert document: %w", err)
	}

	return domain.Document{
		ID:        id,
		Body:      body,
		CreatedAt: time.Now(),
	}, nil
}

func toDomainDocument(doc *models.Document) (domain.Document, error) {
	id, err := uuid.Parse(doc.ID)
	if err != nil {
		return domain.Document{}, fmt.Errorf("parse document ID: %w", err)
	}

	return domain.Document{
		ID:        id,
		Body:      doc.Body,
		CreatedAt: doc.CreatedAt,
	}, nil
}
