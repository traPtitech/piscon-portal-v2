package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) GetScores(c echo.Context) error {
	teamScores, err := h.useCase.GetScores(c.Request().Context())
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]openapi.TeamScores, 0, len(teamScores))
	for _, teamScore := range teamScores {
		scores := make([]openapi.BenchScore, 0, len(teamScore.Scores))
		for _, score := range teamScore.Scores {
			scores = append(scores, toBenchScore(score))
		}

		res = append(res, openapi.TeamScores{
			TeamId: openapi.TeamId(teamScore.TeamID),
			Scores: scores,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetRanking(c echo.Context) error {
	queryOrderBy := c.QueryParam("orderBy")

	var orderBy domain.RankingOrderBy
	switch queryOrderBy {
	case string(openapi.RankingOrderByHighest):
		orderBy = domain.RankingOrderByHighestScore
	case string(openapi.RankingOrderByLatest):
		orderBy = domain.RankingOrderByLatestScore
	default:
		return badRequestResponse(c, fmt.Sprintf("invalid order by: %s", queryOrderBy))
	}

	ranking, err := h.useCase.GetRanking(c.Request().Context(), usecase.RankingQuery{OrderBy: orderBy})
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make(openapi.GetRankingOKApplicationJSON, 0, len(ranking))
	for _, r := range ranking {
		res = append(res, openapi.RankingItem{
			TeamId:    openapi.TeamId(r.TeamID),
			Score:     openapi.Score(r.Score.Score),
			CreatedAt: openapi.CreatedAt(r.CreatedAt),
			Rank:      r.Rank,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func toBenchScore(score domain.Score) openapi.BenchScore {
	return openapi.BenchScore{
		BenchmarkId: openapi.BenchmarkId(score.BenchmarkID),
		TeamId:      openapi.TeamId(score.TeamID),
		Score:       openapi.Score(score.Score),
		CreatedAt:   openapi.CreatedAt(score.CreatedAt),
	}
}
