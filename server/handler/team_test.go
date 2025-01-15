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
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestGetTeams(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	h := NewHandler(useCaseMock, repoMock, nil)

	members := []domain.User{
		{ID: uuid.New()}, {ID: uuid.New()}, {ID: uuid.New()},
	}
	teams := []domain.Team{
		{
			ID:      uuid.New(),
			Name:    "Team A",
			Members: members,
		},
		{
			ID:      uuid.New(),
			Name:    "Team B",
			Members: members,
		},
	}

	tests := []struct {
		name  string
		teams []domain.Team
	}{
		{
			name:  "success",
			teams: teams,
		},
		{
			name:  "empty",
			teams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/teams", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			useCaseMock.EXPECT().GetTeams(gomock.Any()).Return(tt.teams, nil)

			_ = h.GetTeams(c)

			if !assert.Equal(t, http.StatusOK, rec.Code) {
				t.Log(rec.Body.String())
			}
			var res openapi.GetTeamsOKApplicationJSON
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
			assert.Len(t, res, len(tt.teams))
			for i, team := range tt.teams {
				compareTeam(t, team, res[i])
			}
		})
	}
}

func TestCreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	userID := uuid.New()

	e := echo.New()
	req := &openapi.PostTeamReq{
		Name:    "Team A",
		Members: []openapi.UserId{openapi.UserId(userID)},
	}
	httpReq := newJSONRequest(http.MethodPost, "/teams", req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.Set("userID", userID)
	h := NewHandler(useCaseMock, repoMock, nil)

	teamID := uuid.New()
	useCaseMock.EXPECT().CreateTeam(gomock.Any(), usecase.CreateTeamInput{
		Name:      "Team A",
		MemberIDs: []uuid.UUID{userID},
		CreatorID: userID,
	}).Return(domain.Team{
		ID:      teamID,
		Name:    "Team A",
		Members: []domain.User{{ID: userID}},
	}, nil)

	_ = h.CreateTeam(c)

	if !assert.Equal(t, http.StatusCreated, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, teamID, uuid.UUID(res.ID))
}

func TestCreateTeam_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	userID := uuid.New()

	e := echo.New()
	req := &openapi.PostTeamReq{
		Name:    "Team A",
		Members: []openapi.UserId{openapi.UserId(userID)},
	}
	httpReq := newJSONRequest(http.MethodPost, "/teams", req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.Set("userID", userID)
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().CreateTeam(gomock.Any(), usecase.CreateTeamInput{
		Name:      "Team A",
		MemberIDs: []uuid.UUID{userID},
		CreatorID: userID,
	}).Return(domain.Team{}, usecase.NewUseCaseErrorFromMsg("user is already in another team"))

	_ = h.CreateTeam(c)

	if !assert.Equal(t, http.StatusBadRequest, rec.Code) {
		t.Log(rec.Body.String())
	}
}

func TestGetTeam(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/teams/"+teamID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	team := domain.Team{
		ID:   teamID,
		Name: "Team A",
		Members: []domain.User{
			{ID: uuid.New()},
		},
	}

	useCaseMock.EXPECT().GetTeam(gomock.Any(), teamID).Return(team, nil)

	_ = h.GetTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	compareTeam(t, team, res)
}

func TestGetTeam_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/teams/"+teamID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().GetTeam(gomock.Any(), teamID).Return(domain.Team{}, usecase.ErrNotFound)

	_ = h.GetTeam(c)

	if !assert.Equal(t, http.StatusNotFound, rec.Code) {
		t.Log(rec.Body.String())
	}
}

func TestUpdateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	newMemberID := uuid.New()
	req := &openapi.PatchTeamReq{
		Name:    openapi.NewOptTeamName("Updated Team"),
		Members: []openapi.UserId{openapi.UserId(newMemberID)},
	}
	httpReq := newJSONRequest(http.MethodPatch, "/teams/"+teamID.String(), req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().UpdateTeam(gomock.Any(), usecase.UpdateTeamInput{
		ID:        teamID,
		Name:      "Updated Team",
		MemberIDs: []uuid.UUID{newMemberID},
	}).Return(domain.Team{
		ID:      teamID,
		Name:    "Updated Team",
		Members: []domain.User{{ID: newMemberID}},
	}, nil)

	_ = h.UpdateTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, teamID, uuid.UUID(res.ID))
}

func TestUpdateTeam_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	newMemberID := uuid.New()
	req := &openapi.PatchTeamReq{
		Name:    openapi.NewOptTeamName("Updated Team"),
		Members: []openapi.UserId{openapi.UserId(newMemberID)},
	}
	httpReq := newJSONRequest(http.MethodPatch, "/teams/"+teamID.String(), req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().UpdateTeam(gomock.Any(), usecase.UpdateTeamInput{
		ID:        teamID,
		Name:      "Updated Team",
		MemberIDs: []uuid.UUID{newMemberID},
	}).Return(domain.Team{}, usecase.NewUseCaseErrorFromMsg("team is full"))

	_ = h.UpdateTeam(c)

	if !assert.Equal(t, http.StatusBadRequest, rec.Code) {
		t.Log(rec.Body.String())
	}
}

func compareTeam(t *testing.T, expected domain.Team, actual openapi.Team) {
	t.Helper()
	assert.Equal(t, expected.ID, uuid.UUID(actual.ID))
	assert.Equal(t, expected.Name, string(actual.Name))
	assert.Len(t, actual.Members, len(expected.Members))
	for i, member := range expected.Members {
		assert.Equal(t, member.ID, uuid.UUID(actual.Members[i]))
	}
}
