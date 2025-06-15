package usecase_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"go.uber.org/mock/gomock"
)

func TestGetDocument(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type (
		GetDocumentResult struct {
			doc domain.Document
			err error
		}
	)

	doc := domain.Document{
		ID:        uuid.New(),
		Body:      "test document body",
		CreatedAt: time.Now(),
	}

	testCases := map[string]struct {
		GetDocumentResult *GetDocumentResult
		doc               domain.Document
		err               error
	}{
		"正しく取得できる": {
			GetDocumentResult: &GetDocumentResult{doc: doc},
			doc:               doc,
		},
		"ドキュメントが存在しない": {
			GetDocumentResult: &GetDocumentResult{doc: domain.Document{}, err: repository.ErrNotFound},
			err:               usecase.ErrNotFound,
		},
		"GetDocumentでエラーが発生する": {
			GetDocumentResult: &GetDocumentResult{doc: domain.Document{}, err: assert.AnError},
			err:               assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mock.NewMockRepository(ctrl)
			docUseCase := usecase.NewDocumentUseCase(mockRepo)

			if testCase.GetDocumentResult != nil {
				mockRepo.EXPECT().GetDocument(gomock.Any()).Return(testCase.GetDocumentResult.doc, testCase.GetDocumentResult.err)
			}

			doc, err := docUseCase.GetDocument(t.Context())

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
			} else {
				assert.NoError(t, err)
			}

			if err != nil {
				return
			}

			assert.Equal(t, testCase.doc, doc)
		})
	}

}

func TestCreateDocument(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	testCases := map[string]struct {
		body           string
		doc            domain.Document
		CreateDocError error
		err            error
	}{
		"正しくドキュメントを作成できる": {
			body: "test document body",
			doc: domain.Document{
				ID:        uuid.New(),
				Body:      "test document body",
				CreatedAt: time.Now(),
			},
		},
		"CreateDocumentでエラーが発生する": {
			body:           "test document body",
			CreateDocError: assert.AnError,
			err:            assert.AnError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mock.NewMockRepository(ctrl)
			mockRepo.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), testCase.body).
				Return(testCase.doc, testCase.CreateDocError)

			docUseCase := usecase.NewDocumentUseCase(mockRepo)

			doc, err := docUseCase.CreateDocument(t.Context(), testCase.body)

			if testCase.err != nil {
				assert.ErrorIs(t, err, testCase.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.doc.ID, doc.ID)
			assert.Equal(t, testCase.doc.Body, doc.Body)
			assert.WithinDuration(t, testCase.doc.CreatedAt, doc.CreatedAt, time.Second)
		})
	}
}
