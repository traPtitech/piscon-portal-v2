package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestGetUserMe(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	h := NewHandler(useCaseMock, repoMock, nil)

	user := domain.User{
		ID:     uuid.New(),
		Name:   "user1",
		TeamID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	}

	useCaseMock.EXPECT().GetUser(gomock.Any(), user.ID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", user.ID)

	_ = h.GetUserMe(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var res openapi.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	compareUser(t, user, res)
}

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	h := NewHandler(useCaseMock, repoMock, nil)

	users := []domain.User{
		{
			ID:     uuid.New(),
			Name:   "user1",
			TeamID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
		},
		{
			ID:   uuid.New(),
			Name: "user2",
		},
		{
			ID:      uuid.New(),
			Name:    "user3",
			IsAdmin: true,
		},
	}
	useCaseMock.EXPECT().GetUsers(gomock.Any()).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.GetUsers(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var res []openapi.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Len(t, res, 3)
	for i, user := range users {
		compareUser(t, user, res[i])
	}
}

func compareUser(t *testing.T, want domain.User, got openapi.User) {
	assert.Equal(t, want.ID, uuid.UUID(got.ID))
	assert.Equal(t, want.Name, string(got.Name))
	assert.Equal(t, want.IsAdmin, got.IsAdmin)
	if want.TeamID.Valid {
		assert.Equal(t, want.TeamID.UUID, uuid.UUID(got.TeamId.Value))
	} else {
		assert.False(t, got.TeamId.IsSet())
	}
}
