package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/dm"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

func TestGetDocument(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	doc := domain.Document{
		ID:        uuid.New(),
		Body:      "test document body",
		CreatedAt: time.Now(),
	}

	testCases := map[string]struct {
		docsBefore []domain.Document
		doc        domain.Document
		err        error
	}{
		"正しく取得できる": {
			docsBefore: []domain.Document{doc},
			doc:        doc,
		},
		"ドキュメントが存在しない": {
			docsBefore: []domain.Document{},
			err:        repository.ErrNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			for _, d := range testCase.docsBefore {
				mustMakeDocument(t, db, d)
				t.Cleanup(func() {
					_, err := models.Documents.Delete(
						dm.Where(models.Documents.Columns.ID.EQ(mysql.Arg(d.ID.String()))),
					).Exec(context.Background(), db)
					require.NoError(t, err)
				})
			}

			doc, err := repo.GetDocument(t.Context())

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			assert.Equal(t, doc.ID, testCase.doc.ID)
			assert.Equal(t, doc.Body, testCase.doc.Body)
			assert.WithinDuration(t, doc.CreatedAt, testCase.doc.CreatedAt, time.Second)
		})
	}

}

func TestCreateDocument(t *testing.T) {
	t.Parallel()

	repo, db := setupRepository(t)

	doc := domain.Document{
		ID:        uuid.New(),
		Body:      "test document body",
		CreatedAt: time.Now(),
	}

	testCases := map[string]struct {
		doc domain.Document
		err error
	}{
		"正しく作成できる": {
			doc: doc,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			newDoc, err := repo.CreateDocument(t.Context(), testCase.doc.ID, testCase.doc.Body)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			assert.Equal(t, newDoc.ID, testCase.doc.ID)
			assert.Equal(t, newDoc.Body, testCase.doc.Body)
			assert.WithinDuration(t, newDoc.CreatedAt, time.Now(), time.Second)

			resultDoc, err := models.Documents.Query(
				sm.Where(models.Documents.Columns.ID.EQ(mysql.Arg(newDoc.ID.String()))),
				sm.Limit(1),
			).One(t.Context(), db)
			assert.NoError(t, err)
			assert.Equal(t, newDoc.ID.String(), resultDoc.ID)
			assert.Equal(t, newDoc.Body, resultDoc.Body)
			assert.WithinDuration(t, newDoc.CreatedAt, resultDoc.CreatedAt, time.Second)

			t.Cleanup(func() {
				_, err := models.Documents.Delete(
					dm.Where(models.Documents.Columns.ID.EQ(mysql.Arg(newDoc.ID.String()))),
				).Exec(context.Background(), db)
				require.NoError(t, err)
			})
		})
	}
}
