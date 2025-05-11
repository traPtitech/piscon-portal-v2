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

func TestGetRanking(t *testing.T) {
	t.Parallel()

	ranking := []usecase.RankingItem{
		{
			Score: domain.Score{
				BenchmarkID: uuid.New(),
				TeamID:      uuid.New(),
				Score:       1000,
				CreatedAt:   time.Now(),
			},
			Rank: 1,
		},
		{
			Score: domain.Score{
				BenchmarkID: uuid.New(),
				TeamID:      uuid.New(),
				Score:       200,
				CreatedAt:   time.Now(),
			},
			Rank: 2,
		},
	}

	testCases := map[string]struct {
		orderBy           openapi.RankingOrderBy
		executeGetRanking bool
		ranking           []usecase.RankingItem
		GetRankingErr     error
		code              int
		res               openapi.GetRankingOKApplicationJSON
	}{
		"無効なorderByなので400": {
			orderBy: openapi.RankingOrderBy("invalid"),
			code:    http.StatusBadRequest,
		},
		"GetRankingがエラーなので500": {
			orderBy:           openapi.RankingOrderByHighest,
			executeGetRanking: true,
			GetRankingErr:     assert.AnError,
			code:              http.StatusInternalServerError,
		},
		"GetRankingが成功したので200": {
			orderBy:           openapi.RankingOrderByHighest,
			executeGetRanking: true,
			ranking:           ranking,
			code:              http.StatusOK,
			res: openapi.GetRankingOKApplicationJSON{
				{
					TeamId:    openapi.TeamId(ranking[0].Score.TeamID),
					Score:     openapi.Score(ranking[0].Score.Score),
					CreatedAt: openapi.CreatedAt(ranking[0].Score.CreatedAt),
					Rank:      ranking[0].Rank,
				},
				{
					TeamId:    openapi.TeamId(ranking[1].Score.TeamID),
					Score:     openapi.Score(ranking[1].Score.Score),
					CreatedAt: openapi.CreatedAt(ranking[1].Score.CreatedAt),
					Rank:      ranking[1].Rank,
				},
			},
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

			req := httptest.NewRequest(http.MethodGet, "/score/ranking", nil)
			q := req.URL.Query()
			q.Add("orderBy", string(testCase.orderBy))
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("orderBy")
			c.SetParamValues(string(testCase.orderBy))

			if testCase.executeGetRanking {
				var orderBy domain.RankingOrderBy
				switch testCase.orderBy {
				case openapi.RankingOrderByHighest:
					orderBy = domain.RankingOrderByHighestScore
				case openapi.RankingOrderByLatest:
					orderBy = domain.RankingOrderByLatestScore
				default:
					t.Fatalf("invalid order by: %s", testCase.orderBy)
				}

				useCaseMock.EXPECT().
					GetRanking(gomock.Any(), usecase.RankingQuery{OrderBy: orderBy}).
					Return(testCase.ranking, testCase.GetRankingErr)
			}

			_ = h.GetRanking(c)

			assert.Equal(t, testCase.code, rec.Code)

			if rec.Code != http.StatusOK {
				return
			}
			var resBody openapi.GetRankingOKApplicationJSON
			err := json.Unmarshal(rec.Body.Bytes(), &resBody)
			assert.NoError(t, err)
			assert.Equal(t, len(testCase.res), len(resBody))
			for i, want := range testCase.res {
				assert.Equal(t, want.TeamId, resBody[i].TeamId)
				assert.Equal(t, want.Score, resBody[i].Score)
				assert.WithinDuration(t, time.Time(want.CreatedAt), time.Time(resBody[i].CreatedAt), time.Second)
				assert.Equal(t, want.Rank, resBody[i].Rank)
			}
		})
	}
}
