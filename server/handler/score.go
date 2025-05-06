package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
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

func toBenchScore(score domain.Score) openapi.BenchScore {
	return openapi.BenchScore{
		BenchmarkId: openapi.BenchmarkId(score.BenchmarkID),
		TeamId:      openapi.TeamId(score.TeamID),
		Score:       openapi.Score(score.Score),
		CreatedAt:   openapi.CreatedAt(score.CreatedAt),
	}
}
