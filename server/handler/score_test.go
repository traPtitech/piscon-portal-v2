package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"go.uber.org/mock/gomock"
)

func TestGetScores(t *testing.T) {
	t.Parallel()

	teamID1 := uuid.New()
	teamScores1 := usecase.TeamScores{
		TeamID: teamID1,
		Scores: []domain.Score{
			{BenchmarkID: uuid.New(), TeamID: teamID1, Score: 100, CreatedAt: time.Now()},
		},
	}
	openapiTeamScores1 := openapi.GetScoresOKApplicationJSON{openapi.TeamScores{
		TeamId: openapi.TeamId(teamID1),
		Scores: []openapi.BenchScore{
			{
				BenchmarkId: openapi.BenchmarkId(teamScores1.Scores[0].BenchmarkID),
				TeamId:      openapi.TeamId(teamScores1.Scores[0].TeamID),
				Score:       openapi.Score(teamScores1.Scores[0].Score),
				CreatedAt:   openapi.CreatedAt(teamScores1.Scores[0].CreatedAt),
			},
		},
	}}

	teamID2 := uuid.New()
	teamScores2 := usecase.TeamScores{
		TeamID: teamID2,
		Scores: []domain.Score{
			{BenchmarkID: uuid.New(), TeamID: teamID2, Score: 200, CreatedAt: time.Now()},
			{BenchmarkID: uuid.New(), TeamID: teamID2, Score: 300, CreatedAt: time.Now()},
		},
	}
	openapiTeamScores2 := openapi.GetScoresOKApplicationJSON{
		{
			TeamId: openapi.TeamId(teamID1),
			Scores: []openapi.BenchScore{
				{
					BenchmarkId: openapi.BenchmarkId(teamScores1.Scores[0].BenchmarkID),
					TeamId:      openapi.TeamId(teamScores1.Scores[0].TeamID),
					Score:       openapi.Score(teamScores1.Scores[0].Score),
					CreatedAt:   openapi.CreatedAt(teamScores1.Scores[0].CreatedAt),
				},
			},
		},
		{
			TeamId: openapi.TeamId(teamID2),
			Scores: []openapi.BenchScore{
				{
					BenchmarkId: openapi.BenchmarkId(teamScores2.Scores[0].BenchmarkID),
					TeamId:      openapi.TeamId(teamScores2.Scores[0].TeamID),
					Score:       openapi.Score(teamScores2.Scores[0].Score),
					CreatedAt:   openapi.CreatedAt(teamScores2.Scores[0].CreatedAt),
				},
				{
					BenchmarkId: openapi.BenchmarkId(teamScores2.Scores[1].BenchmarkID),
					TeamId:      openapi.TeamId(teamScores2.Scores[1].TeamID),
					Score:       openapi.Score(teamScores2.Scores[1].Score),
					CreatedAt:   openapi.CreatedAt(teamScores2.Scores[1].CreatedAt),
				},
			},
		},
	}

	testCases := map[string]struct {
		teamScores   []usecase.TeamScores
		GetScoresErr error
		resBosy      openapi.GetScoresOKApplicationJSON
		code         int
	}{
		"GetScoresがエラーなので500": {
			teamScores:   nil,
			GetScoresErr: assert.AnError,
			code:         http.StatusInternalServerError,
		},
		"何もない": {
			teamScores: []usecase.TeamScores{},
			code:       http.StatusOK,
		},
		"GetScoresが1件": {
			teamScores: []usecase.TeamScores{teamScores1},
			code:       http.StatusOK,
			resBosy:    openapiTeamScores1,
		},
		"GetScoresが複数件": {
			teamScores: []usecase.TeamScores{teamScores1, teamScores2},
			code:       http.StatusOK,
			resBosy:    openapiTeamScores2,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)

			e := echo.New()
			h := NewHandler(useCaseMock, repoMock, nil)

			req := httptest.NewRequest(http.MethodGet, "/scores", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			useCaseMock.EXPECT().GetScores(gomock.Any()).Return(testCase.teamScores, testCase.GetScoresErr)

			_ = h.GetScores(c)

			assert.Equal(t, testCase.code, rec.Code)

			if rec.Code != http.StatusOK {
				return
			}

			var resBody openapi.GetScoresOKApplicationJSON
			err := json.Unmarshal(rec.Body.Bytes(), &resBody)
			assert.NoError(t, err)
			assert.Equal(t, len(testCase.teamScores), len(resBody))
			for i, teamScore := range testCase.teamScores {
				assert.Equal(t, teamScore.TeamID, uuid.UUID(resBody[i].TeamId))
				assert.Equal(t, len(teamScore.Scores), len(resBody[i].Scores))
				for j, score := range teamScore.Scores {
					assert.Equal(t, score.BenchmarkID, uuid.UUID(resBody[i].Scores[j].BenchmarkId))
					assert.Equal(t, score.TeamID, uuid.UUID(resBody[i].Scores[j].TeamId))
					assert.Equal(t, score.Score, int64(resBody[i].Scores[j].Score))
					assert.WithinDuration(t, score.CreatedAt, time.Time(resBody[i].Scores[j].CreatedAt), time.Second)
				}
			}

		})
	}
}
