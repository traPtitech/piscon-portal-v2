package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
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
			ID:        uuid.New(),
			Name:      "Team B",
			Members:   members,
			GitHubIDs: []string{"user1", "user2"},
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

			if !assert.Equal(t, http.StatusOK, rec.Code, "status code") {
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
		Name:      "Team A",
		Members:   []openapi.UserId{openapi.UserId(userID)},
		GithubIds: []openapi.GitHubId{"user1", "user2"},
	}
	httpReq := newJSONRequest(http.MethodPost, "/teams", req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.Set(handler.UserIDKey, userID)
	h := NewHandler(useCaseMock, repoMock, nil)

	teamID := uuid.New()
	useCaseMock.EXPECT().CreateTeam(gomock.Any(), usecase.CreateTeamInput{
		Name:      "Team A",
		MemberIDs: []uuid.UUID{userID},
		CreatorID: userID,
		GitHubIDs: []string{"user1", "user2"},
	}).Return(domain.Team{
		ID:      teamID,
		Name:    "Team A",
		Members: []domain.User{{ID: userID}},
	}, nil)

	_ = h.CreateTeam(c)

	if !assert.Equal(t, http.StatusCreated, rec.Code, "status code") {
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
	c.Set(handler.UserIDKey, userID)
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().CreateTeam(gomock.Any(), usecase.CreateTeamInput{
		Name:      "Team A",
		MemberIDs: []uuid.UUID{userID},
		CreatorID: userID,
		GitHubIDs: []string{},
	}).Return(domain.Team{}, usecase.NewUseCaseErrorFromMsg("user is already in another team"))

	_ = h.CreateTeam(c)

	if !assert.Equal(t, http.StatusBadRequest, rec.Code, "status code") {
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
		GitHubIDs: []string{"user1"},
	}

	useCaseMock.EXPECT().GetTeam(gomock.Any(), teamID).Return(team, nil)

	_ = h.GetTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code, "status code") {
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

	if !assert.Equal(t, http.StatusNotFound, rec.Code, "status code") {
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
		Name:      lo.ToPtr("Updated Team"),
		MemberIDs: lo.ToPtr([]uuid.UUID{newMemberID}),
		GitHubIDs: nil,
	}).Return(domain.Team{
		ID:      teamID,
		Name:    "Updated Team",
		Members: []domain.User{{ID: newMemberID}},
	}, nil)

	_ = h.UpdateTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code, "status code") {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, teamID, uuid.UUID(res.ID))
}

func TestUpdateTeamWithGitHubIDs(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	newMemberID := uuid.New()
	req := &openapi.PatchTeamReq{
		Name:      openapi.NewOptTeamName("Updated Team"),
		Members:   []openapi.UserId{openapi.UserId(newMemberID)},
		GithubIds: []openapi.GitHubId{"updated1", "updated2"},
	}
	httpReq := newJSONRequest(http.MethodPatch, "/teams/"+teamID.String(), req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().UpdateTeam(gomock.Any(), usecase.UpdateTeamInput{
		ID:        teamID,
		Name:      lo.ToPtr("Updated Team"),
		MemberIDs: lo.ToPtr([]uuid.UUID{newMemberID}),
		GitHubIDs: lo.ToPtr([]string{"updated1", "updated2"}),
	}).Return(domain.Team{
		ID:        teamID,
		Name:      "Updated Team",
		Members:   []domain.User{{ID: newMemberID}},
		GitHubIDs: []string{"updated1", "updated2"},
	}, nil)

	_ = h.UpdateTeam(c)

	if !assert.Equal(t, http.StatusOK, rec.Code, "status code") {
		t.Log(rec.Body.String())
	}
	var res openapi.Team
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, teamID, uuid.UUID(res.ID))
	assert.Equal(t, []openapi.GitHubId{"updated1", "updated2"}, res.GithubIds)
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
		Name:      lo.ToPtr("Updated Team"),
		MemberIDs: lo.ToPtr([]uuid.UUID{newMemberID}),
		GitHubIDs: nil,
	}).Return(domain.Team{}, usecase.NewUseCaseErrorFromMsg("team is full"))

	_ = h.UpdateTeam(c)

	if !assert.Equal(t, http.StatusBadRequest, rec.Code, "status code") {
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
	assert.ElementsMatch(t, expected.GitHubIDs,
		lo.Map(actual.GithubIds, func(id openapi.GitHubId, _ int) string { return string(id) }))
}

// Exercise JSON binding and the real use case together so omitted fields stay omitted.
func TestUpdateTeamPartialUpdate(t *testing.T) {
	teamID := uuid.New()
	member := domain.User{ID: uuid.New(), TeamID: uuid.NullUUID{UUID: teamID, Valid: true}}
	otherMember := domain.User{ID: uuid.New(), TeamID: uuid.NullUUID{UUID: teamID, Valid: true}}
	original := domain.Team{
		ID: teamID, Name: "Team A", Members: []domain.User{member, otherMember},
		GitHubIDs: []string{"existing"}, CreatedAt: time.Now().UTC(),
	}
	tests := []struct {
		name  string
		body  map[string]any
		input usecase.UpdateTeamInput
		want  domain.Team
	}{
		{
			name:  "github only preserves all members",
			body:  map[string]any{"githubIds": []string{"existing", "added"}},
			input: usecase.UpdateTeamInput{GitHubIDs: lo.ToPtr([]string{"existing", "added"})},
			want:  domain.Team{Name: original.Name, Members: original.Members, GitHubIDs: []string{"existing", "added"}},
		},
		{
			name:  "name only preserves members and github ids",
			body:  map[string]any{"name": "Renamed"},
			input: usecase.UpdateTeamInput{Name: lo.ToPtr("Renamed")},
			want:  domain.Team{Name: "Renamed", Members: original.Members, GitHubIDs: original.GitHubIDs},
		},
		{
			name:  "explicit empty name updates name",
			body:  map[string]any{"name": ""},
			input: usecase.UpdateTeamInput{Name: lo.ToPtr("")},
			want:  domain.Team{Name: "", Members: original.Members, GitHubIDs: original.GitHubIDs},
		},
		{
			name: "empty patch preserves team",
			body: map[string]any{},
			want: original,
		},
		{
			name:  "empty github ids removes last id and preserves members",
			body:  map[string]any{"githubIds": []string{}},
			input: usecase.UpdateTeamInput{GitHubIDs: lo.ToPtr([]string{})},
			want:  domain.Team{Name: original.Name, Members: original.Members, GitHubIDs: []string{}},
		},
		{
			name:  "empty members removes all members and preserves github ids",
			body:  map[string]any{"members": []string{}},
			input: usecase.UpdateTeamInput{MemberIDs: lo.ToPtr([]uuid.UUID{})},
			want:  domain.Team{Name: original.Name, Members: []domain.User{}, GitHubIDs: original.GitHubIDs},
		},
		{
			name:  "member removal preserves remaining member and github ids",
			body:  map[string]any{"members": []uuid.UUID{member.ID}},
			input: usecase.UpdateTeamInput{MemberIDs: lo.ToPtr([]uuid.UUID{member.ID})},
			want:  domain.Team{Name: original.Name, Members: []domain.User{member}, GitHubIDs: original.GitHubIDs},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := repomock.NewMockRepository(ctrl)
			uc := usecasemock.NewMockUseCase(ctrl)
			teamUC := usecase.NewTeamUseCase(repo)
			tt.input.ID = teamID
			tt.want.ID = teamID
			tt.want.CreatedAt = original.CreatedAt
			repo.EXPECT().FindTeam(gomock.Any(), teamID).Return(original, nil)
			if tt.input.MemberIDs != nil {
				for _, id := range *tt.input.MemberIDs {
					repo.EXPECT().FindUser(gomock.Any(), id).Return(member, nil)
				}
			}
			repo.EXPECT().UpdateTeam(gomock.Any(), tt.want).Return(nil)
			uc.EXPECT().UpdateTeam(gomock.Any(), tt.input).DoAndReturn(
				func(ctx context.Context, input usecase.UpdateTeamInput) (domain.Team, error) {
					return teamUC.UpdateTeam(ctx, input)
				},
			)
			e := echo.New()
			body, err := json.Marshal(tt.body)
			require.NoError(t, err)
			req := newJSONRequest(http.MethodPatch, "/teams/"+teamID.String(), json.RawMessage(body))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("teamID")
			c.SetParamValues(teamID.String())
			h := NewHandler(uc, repo, nil)
			require.NoError(t, h.UpdateTeam(c))
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			var res openapi.Team
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
			compareTeam(t, tt.want, res)
		})
	}
}
