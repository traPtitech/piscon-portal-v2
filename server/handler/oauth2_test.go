package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	usecaseMock := usecasemock.NewMockUseCase(ctrl)

	server := NewPortalServer(usecaseMock, mockRepo)
	client := NewClient(server)
	userID := uuid.New()

	testFirstLogin(t, mockRepo, server, client, userID)
}

func TestLoginAsExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	usecaseMock := usecasemock.NewMockUseCase(ctrl)

	server := NewPortalServer(usecaseMock, mockRepo)
	client := NewClient(server)
	userID := uuid.New()

	// user already exists, so only create session
	mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		})
	mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(domain.User{ID: userID}, nil)
	mockRepo.EXPECT().
		CreateSession(gomock.Any(), gomock.Cond(func(s domain.Session) bool { return s.UserID == userID })).
		Return(nil)

	if err := Login(t, server, client, userID); err != nil {
		t.FailNow()
	}
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repomock.NewMockRepository(ctrl)
	mockUseCase := usecasemock.NewMockUseCase(ctrl)

	server := NewPortalServer(mockUseCase, mockRepo)
	client := NewClient(server)
	userID := uuid.New()

	testFirstLogin(t, mockRepo, server, client, userID)

	// logout
	// return not expired session
	mockRepo.EXPECT().FindSession(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, sid string) (domain.Session, error) {
			return domain.Session{
				ID:        sid,
				UserID:    userID,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil
		})
	mockRepo.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil)

	res, err := client.Post(joinPath(t, server.URL, "/api/oauth2/logout"), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(res.Body)
		t.Fatalf("status code is %d: %s", res.StatusCode, msg)
	}
}

func testFirstLogin(t *testing.T, mockRepo *repomock.MockRepository, server *httptest.Server, client *http.Client, userID uuid.UUID) {
	// create user and session
	mockRepo.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		})
	mockRepo.EXPECT().FindUser(gomock.Any(), gomock.Eq(userID)).Return(domain.User{}, repository.ErrNotFound)
	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Cond(func(u domain.User) bool { return u.ID == userID })).
		Return(nil)
	mockRepo.EXPECT().
		CreateSession(gomock.Any(), gomock.Cond(func(s domain.Session) bool { return s.UserID == userID })).
		Return(nil)

	if err := Login(t, server, client, userID); err != nil {
		t.FailNow()
	}
}
