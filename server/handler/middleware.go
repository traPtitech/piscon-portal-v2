package handler

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
)

const userIDKey = "userID"

func getUserID(c echo.Context) string {
	return c.Get(userIDKey).(string)
}

func (h *Handler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			sessID, err := h.sessionManager.GetSessionID(c)
			if err != nil {
				return internalServerErrorResponse(c, err)
			}
			if sessID == "" {
				return unauthorizedResponse(c, "session not found")
			}

			session, err := h.repo.FindSession(ctx, sessID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return unauthorizedResponse(c, "session not found")
				}
				return internalServerErrorResponse(c, err)
			}
			if session.ExpiresAt.Before(time.Now()) {
				// delete expired session
				if err := h.repo.DeleteSession(ctx, sessID); err != nil {
					return internalServerErrorResponse(c, err)
				}
				return unauthorizedResponse(c, "session expired")
			}

			c.Set(userIDKey, session.UserID)

			return next(c)
		}
	}
}

// TeamAuthMiddleware is a middleware that checks if the user is a member of the team.
// The team ID is taken from the URL parameter.
func (h *Handler) TeamAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			userID := getUserID(c)
			teamID := c.Param("teamID")
			if teamID == "" {
				return c.NoContent(http.StatusBadRequest)
			}

			team, err := h.repo.FindTeam(ctx, teamID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return c.NoContent(http.StatusNotFound)
				}
				return internalServerErrorResponse(c, err)
			}

			isMember := slices.ContainsFunc(team.Members, func(m domain.User) bool { return m.ID == userID })
			if !isMember {
				return c.NoContent(http.StatusForbidden)
			}

			return next(c)
		}
	}
}
