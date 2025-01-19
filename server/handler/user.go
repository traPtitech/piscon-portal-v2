package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
)

func (h *Handler) GetUserMe(c echo.Context) error {
	ctx := c.Request().Context()

	userID := getUserIDFromSession(c)

	user, err := h.useCase.GetUser(ctx, userID)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := toOpenAPIUser(user)

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.useCase.GetUsers(ctx)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]*openapi.User, 0, len(users))
	for _, user := range users {
		res = append(res, toOpenAPIUser(user))
	}

	return c.JSON(http.StatusOK, res)
}

func toOpenAPIUser(user domain.User) *openapi.User {
	res := &openapi.User{
		ID:      openapi.UserId(user.ID),
		Name:    openapi.UserName(user.Name),
		IsAdmin: user.IsAdmin,
	}
	if user.TeamID.Valid {
		res.TeamId = openapi.NewOptTeamId(openapi.TeamId(user.TeamID.UUID))
	}
	return res
}
