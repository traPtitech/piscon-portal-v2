package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	sessmock "github.com/traPtitech/piscon-portal-v2/server/handler/internal/mock"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	sessMock := sessmock.NewMockSessionManager(ctrl)
	usecaseMock := usecasemock.NewMockUseCase(ctrl)
	handler := NewHandler(usecaseMock, repoMock, sessMock)

	needAuthorize := handler.AuthMiddleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tests := []struct {
		name         string
		setup        func()
		expectStatus int
	}{
		{
			name: "ok",
			setup: func() {
				sessMock.EXPECT().GetSessionID(gomock.Any()).Return("sessionID", nil)
				repoMock.EXPECT().FindSession(gomock.Any(), "sessionID").Return(domain.Session{
					ID:        "sessionID",
					UserID:    uuid.New(),
					ExpiresAt: time.Now().Add(time.Hour),
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "session not found",
			setup: func() {
				sessMock.EXPECT().GetSessionID(gomock.Any()).Return("", nil)
			},
			expectStatus: http.StatusUnauthorized,
		},
		{
			name: "expired session",
			setup: func() {
				sessMock.EXPECT().GetSessionID(gomock.Any()).Return("sessionID", nil)
				repoMock.EXPECT().FindSession(gomock.Any(), "sessionID").Return(domain.Session{
					ID:        "sessionID",
					UserID:    uuid.New(),
					ExpiresAt: time.Now().Add(-time.Hour),
				}, nil)
				// expired session should be deleted
				repoMock.EXPECT().DeleteSession(gomock.Any(), "sessionID").Return(nil)
			},
			expectStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			_ = needAuthorize(echo.New().NewContext(req, rec))
			if rec.Code != tt.expectStatus {
				t.Errorf("unexpected status code: expected=%d, got=%d", tt.expectStatus, rec.Code)
			}
		})
	}
}

func TestTeamAuthMiddleware(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	sessMock := sessmock.NewMockSessionManager(ctrl)
	usecaseMock := usecasemock.NewMockUseCase(ctrl)
	handler := NewHandler(usecaseMock, repoMock, sessMock)

	userID := uuid.New()
	teamID := uuid.New()

	needAuthorize := handler.TeamAuthMiddleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tests := []struct {
		name         string
		setup        func()
		expectStatus int
	}{
		{
			name: "ok",
			setup: func() {
				repoMock.EXPECT().FindTeam(gomock.Any(), teamID).Return(domain.Team{
					ID: teamID,
					Members: []domain.User{
						{ID: userID},
					},
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "user is not a member of the team",
			setup: func() {
				repoMock.EXPECT().FindTeam(gomock.Any(), teamID).Return(domain.Team{
					ID:      uuid.New(),
					Members: []domain.User{{ID: uuid.New()}},
				}, nil)
			},
			expectStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			c.SetParamNames("teamID")
			c.SetParamValues(teamID.String())
			c.Set("userID", userID)

			tt.setup()

			_ = needAuthorize(c)
			if rec.Code != tt.expectStatus {
				t.Errorf("unexpected status code: expected=%d, got=%d", tt.expectStatus, rec.Code)
			}
		})
	}
}
