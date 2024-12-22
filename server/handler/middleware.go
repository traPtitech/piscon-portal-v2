package handler

import (
	"database/sql"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/models"
)

func (h *Handler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			sessID, err := h.sessionManager.getSessionID(c)
			if err != nil {
				return internalServerErrorResponse(c, err.Error())
			}

			session, err := models.Sessions.Query(models.SelectWhere.Sessions.ID.EQ(sessID)).One(ctx, h.db)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return unauthorizedResponse(c, "session not found")
				}
				return internalServerErrorResponse(c, err.Error())
			}
			if session.ExpiredAt.Before(time.Now()) {
				// delete expired session
				_, err := models.Sessions.Delete(models.DeleteWhere.Sessions.ID.EQ(sessID)).Exec(ctx, h.db)
				if err != nil {
					return internalServerErrorResponse(c, err.Error())
				}
				return unauthorizedResponse(c, "session expired")
			}

			return next(c)
		}
	}
}
