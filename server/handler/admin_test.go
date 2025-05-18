package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestPutAdmins(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type PutAdminsResult struct {
		err error
	}

	userID := uuid.New()

	testCases := map[string]struct {
		reqBody         []string
		loginUserID     uuid.UUID
		PutAdminsResult *PutAdminsResult
		expectedStatus  int
	}{
		"無効なuuidが含まれているので400": {
			reqBody:        []string{"invalid-uuid"},
			loginUserID:    uuid.New(),
			expectedStatus: http.StatusBadRequest,
		},
		"PutAdminsでエラーが発生したので500": {
			reqBody:         []string{userID.String()},
			loginUserID:     uuid.New(),
			PutAdminsResult: &PutAdminsResult{err: assert.AnError},
			expectedStatus:  http.StatusInternalServerError,
		},
		"PutAdminsでUseCaseErrorが発生したので400": {
			reqBody:         []string{userID.String()},
			loginUserID:     uuid.New(),
			PutAdminsResult: &PutAdminsResult{err: usecase.NewUseCaseError(assert.AnError)},
			expectedStatus:  http.StatusBadRequest,
		},
		"正常系": {
			reqBody:         []string{userID.String()},
			loginUserID:     uuid.New(),
			PutAdminsResult: &PutAdminsResult{},
			expectedStatus:  http.StatusOK,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)

			if testCase.PutAdminsResult != nil {
				useCaseMock.EXPECT().
					PutAdmins(gomock.Any(), testCase.loginUserID, []uuid.UUID{userID}).
					Return(testCase.PutAdminsResult.err)
			}

			e := echo.New()

			bodyJSON, err := json.Marshal(testCase.reqBody)
			require.NoError(t, err)
			body := bytes.NewBuffer(bodyJSON)

			req := httptest.NewRequest(http.MethodGet, "/admins", body)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(handler.UserIDKey, testCase.loginUserID)
			h := NewHandler(useCaseMock, repoMock, nil)

			_ = h.PutAdmins(c)

			assert.Equal(t, testCase.expectedStatus, rec.Code)
		})
	}

}
