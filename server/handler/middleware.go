package handler

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

func (h *Handler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			sessID, err := h.sessionManager.getSessionID(c)
			if err != nil {
				return internalServerErrorResponse(c, err.Error())
			}

			session, err := h.repo.FindSession(ctx, sessID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return unauthorizedResponse(c, "session not found")
				}
				return internalServerErrorResponse(c, err.Error())
			}
			if session.ExpiresAt.Before(time.Now()) {
				// delete expired session
				if err := h.repo.DeleteSession(ctx, sessID); err != nil {
					return internalServerErrorResponse(c, err.Error())
				}
				return unauthorizedResponse(c, "session expired")
			}

			return next(c)
		}
	}
}
