package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestGetDocument(t *testing.T) {
	t.Parallel()

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
		code              int
		res               openapi.GetDocsOK
	}{
		"正しく取得できる": {
			GetDocumentResult: &GetDocumentResult{doc: doc},
			code:              http.StatusOK,
			res: openapi.GetDocsOK{
				Body: openapi.NewOptMarkdownDocument(openapi.MarkdownDocument(doc.Body)),
			},
		},
		"ドキュメントが存在しない": {
			GetDocumentResult: &GetDocumentResult{doc: domain.Document{}, err: usecase.ErrNotFound},
			code:              http.StatusOK,
			res: openapi.GetDocsOK{
				Body: openapi.OptMarkdownDocument{},
			},
		},
		"GetDocumentでエラーが発生する": {
			GetDocumentResult: &GetDocumentResult{doc: domain.Document{}, err: assert.AnError},
			code:              http.StatusInternalServerError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/docs", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := NewHandler(useCaseMock, repoMock, nil)

			if testCase.GetDocumentResult != nil {
				useCaseMock.EXPECT().GetDocument(gomock.Any()).
					Return(testCase.GetDocumentResult.doc, testCase.GetDocumentResult.err)
			}

			_ = h.GetDocument(c)

			assert.Equal(t, testCase.code, rec.Code)

			if testCase.code != http.StatusOK {
				return
			}

			var body openapi.GetDocsOK

			b := rec.Body.Bytes()
			err := json.Unmarshal(b, &body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.res.Body, body.Body)
		})
	}
}

func TestPatchDocument(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	testCases := map[string]struct {
		reqBody             openapi.PatchDocsReq
		doc                 domain.Document
		CreateDocumentError error
		code                int
		res                 openapi.PatchDocsOK
	}{
		"正しくドキュメントを作成できる": {
			reqBody: openapi.PatchDocsReq{
				Body: openapi.MarkdownDocument("test document body"),
			},
			doc: domain.Document{
				ID:        uuid.New(),
				Body:      "test document body",
				CreatedAt: time.Now(),
			},
			res: openapi.PatchDocsOK{
				Body: openapi.NewOptMarkdownDocument(openapi.MarkdownDocument("test document body")),
			},
			code: http.StatusOK,
		},
		"CreateDocumentでエラーが発生する": {
			reqBody: openapi.PatchDocsReq{
				Body: openapi.MarkdownDocument("test document body"),
			},
			CreateDocumentError: assert.AnError,
			code:                http.StatusInternalServerError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			usecaseMock := usecasemock.NewMockUseCase(ctrl)

			usecaseMock.EXPECT().CreateDocument(gomock.Any(), string(testCase.reqBody.Body)).
				Return(testCase.doc, testCase.CreateDocumentError)

			h := NewHandler(usecaseMock, nil, nil)

			var reqBody bytes.Buffer
			err := json.NewEncoder(&reqBody).Encode(testCase.reqBody)
			require.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/docs", &reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = h.PatchDocument(c)

			assert.Equal(t, testCase.code, rec.Code)

			if testCase.code != http.StatusOK {
				return
			}

			assert.NoError(t, err)
			var body openapi.PatchDocsOK
			b := rec.Body.Bytes()
			err = json.Unmarshal(b, &body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.res.Body, body.Body)
		})
	}
}
