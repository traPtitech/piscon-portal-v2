package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/internal/mock"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	mockSessManager := mock.NewMockSessionManager(ctrl)
	handler := NewHandler(mockRepo, mockSessManager)

	needAuthorize := handler.AuthMiddleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	mockSessManager.EXPECT().GetSessionID(gomock.Any()).Return("sessionID", nil)
	mockRepo.EXPECT().FindSession(gomock.Any(), "sessionID").Return(domain.Session{
		ID:        "sessionID",
		UserID:    "test-user-id",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	needAuthorize(echo.New().NewContext(req, rec))
	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: %d", rec.Code)
	}
}

func TestAuthMiddleware_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	mockSessManager := mock.NewMockSessionManager(ctrl)
	handler := NewHandler(mockRepo, mockSessManager)

	needAuthorize := handler.AuthMiddleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	mockSessManager.EXPECT().GetSessionID(gomock.Any()).Return("", nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	needAuthorize(echo.New().NewContext(req, rec))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("unexpected status code: %d", rec.Code)
	}
}

func TestAuthMiddleware_ExpiredSession(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	mockSessManager := mock.NewMockSessionManager(ctrl)
	handler := NewHandler(mockRepo, mockSessManager)

	needAuthorize := handler.AuthMiddleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	mockSessManager.EXPECT().GetSessionID(gomock.Any()).Return("sessionID", nil)
	mockRepo.EXPECT().FindSession(gomock.Any(), "sessionID").Return(domain.Session{
		ID:        "sessionID",
		UserID:    "test-user-id",
		ExpiresAt: time.Now().Add(-time.Hour),
	}, nil)
	mockRepo.EXPECT().DeleteSession(gomock.Any(), "sessionID").Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	needAuthorize(echo.New().NewContext(req, rec))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("unexpected status code: %d", rec.Code)
	}
}
