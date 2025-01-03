package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) GetTeams(c echo.Context) error {
	ctx := c.Request().Context()

	teams, err := h.teamUseCase.GetTeams(ctx)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make(openapi.GetTeamsOKApplicationJSON, 0, len(teams))
	for _, team := range teams {
		members := lo.Map(team.Members, func(u domain.User, _ int) openapi.UserId {
			return openapi.UserId(uuid.MustParse(u.ID))
		})
		res = append(res, openapi.Team{
			ID:        openapi.TeamId(uuid.MustParse(team.ID)),
			Name:      openapi.TeamName(team.Name),
			Members:   members,
			CreatedAt: team.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateTeam(c echo.Context) error {
	ctx := c.Request().Context()

	var req openapi.PostTeamReq
	if err := c.Bind(&req); err != nil {
		return badRequestResponse(c, err.Error())
	}
	userID := getUserID(c)

	team, err := h.teamUseCase.CreateTeam(ctx, usecase.CreateTeamInput{
		Name: string(req.Name),
		MemberIDs: lo.Map(req.Members, func(id openapi.UserId, _ int) string {
			return uuid.UUID(id).String()
		}),
		CreatorID: userID,
	})
	if err != nil {
		if usecase.IsErrBadRequest(err) {
			return badRequestResponse(c, err.Error())
		}
		return internalServerErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, toOpenAPITeam(team))
}

func (h *Handler) GetTeam(c echo.Context) error {
	ctx := c.Request().Context()

	teamID := c.Param("teamID")

	team, err := h.repo.FindTeam(ctx, teamID)
	if err != nil {
		if err == repository.ErrNotFound {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toOpenAPITeam(team))
}

func (h *Handler) UpdateTeam(c echo.Context) error {
	ctx := c.Request().Context()

	var req openapi.PatchTeamReq
	if err := c.Bind(&req); err != nil {
		return badRequestResponse(c, err.Error())
	}
	teamID := c.Param("teamID")

	team, err := h.teamUseCase.UpdateTeam(ctx, usecase.UpdateTeamInput{
		ID:        teamID,
		Name:      string(req.Name.Value),
		MemberIDs: lo.Map(req.Members, func(id openapi.UserId, _ int) string { return uuid.UUID(id).String() }),
	})
	if err != nil {
		if usecase.IsErrBadRequest(err) {
			return badRequestResponse(c, err.Error())
		}
		return internalServerErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toOpenAPITeam(team))
}

func toOpenAPITeam(team domain.Team) openapi.Team {
	return openapi.Team{
		ID:        openapi.TeamId(uuid.MustParse(team.ID)),
		Name:      openapi.TeamName(team.Name),
		Members:   lo.Map(team.Members, func(m domain.User, _ int) openapi.UserId { return openapi.UserId(uuid.MustParse(m.ID)) }),
		CreatedAt: team.CreatedAt,
	}
}
