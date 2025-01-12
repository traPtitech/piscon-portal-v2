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
	req := httptest.NewRequest(http.MethodGet, "/teams", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(useCaseMock, repoMock, nil)

	teamID := uuid.New()
	memberIDs := []uuid.UUID{uuid.New(), uuid.New()}

	useCaseMock.EXPECT().GetTeams(gomock.Any()).Return([]domain.Team{
		{
			ID:   teamID,
			Name: "Team A",
			Members: []domain.User{
				{ID: memberIDs[0]}, {ID: memberIDs[1]},
			},
		},
	}, nil)

	_ = h.GetTeams(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.GetTeamsOKApplicationJSON
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Len(t, res, 1)
	assert.Equal(t, teamID, uuid.UUID(res[0].ID))
}

func TestCreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	userID := uuid.New()

	e := echo.New()
	req := openapi.PostTeamReq{
		Name:    "Team A",
		Members: []openapi.UserId{openapi.UserId(userID)},
	}
	body, _ := req.MarshalJSON()
	httpReq := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

	memberID := uuid.New()
	useCaseMock.EXPECT().GetTeam(gomock.Any(), teamID).Return(domain.Team{
		ID:   teamID,
		Name: "Team A",
		Members: []domain.User{
			{ID: memberID},
		},
	}, nil)

	_ = h.GetTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, teamID, uuid.UUID(res.ID))
}

func TestUpdateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	newMemberID := uuid.New()
	req := openapi.PatchTeamReq{
		Name:    openapi.NewOptTeamName("Updated Team"),
		Members: []openapi.UserId{openapi.UserId(newMemberID)},
	}
	body, _ := req.MarshalJSON()
	httpReq := httptest.NewRequest(http.MethodPatch, "/teams/"+teamID.String(), bytes.NewReader(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
