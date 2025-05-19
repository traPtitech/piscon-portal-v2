package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) PutAdmins(c echo.Context) error {
	var req []string
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return badRequestResponse(c, "invalid request body")
	}
	defer c.Request().Body.Close()

	ids := make([]uuid.UUID, 0, len(req))
	for _, id := range req {
		userUUID, err := uuid.Parse(id)
		if err != nil {
			return badRequestResponse(c, fmt.Sprintf("invalid uuid: '%s'", id))
		}
		ids = append(ids, userUUID)
	}

	loginUserID := getUserIDFromSession(c)

	err = h.useCase.PutAdmins(c.Request().Context(), loginUserID, ids)
	if usecase.IsUseCaseError(err) {
		return badRequestResponse(c, err.Error())
	}
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	return c.NoContent(http.StatusOK)
}
