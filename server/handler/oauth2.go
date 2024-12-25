package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/utils/random"
)

func (h *Handler) GetOauth2Code(c echo.Context) error {
	sessID, err := h.sessionManager.setSessionID(c, 7*24*time.Hour) // max age 1 week
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	state := random.String(16)
	url := h.oauth2Service.AuthCodeURL(sessID, state)

	return c.Redirect(http.StatusSeeOther, url)
}

func (h *Handler) Oauth2Callback(c echo.Context) error {
	ctx := c.Request().Context()

	sessionID, err := h.sessionManager.getSessionID(c)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if !h.oauth2Service.VerifyState(sessionID, state) {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	token, err := h.oauth2Service.Exchange(ctx, sessionID, code)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}
	userInfo, err := h.oauth2Service.GetUserInfo(ctx, token)
	if err != nil {
		c.Logger().Warn(err)
		return c.String(http.StatusBadRequest, "invalid token")
	}

	err = h.repo.Transaction(ctx, func(ctx context.Context, r repository.Repository) error {
		user, err := r.FindUser(ctx, userInfo.ID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		if errors.Is(err, repository.ErrNotFound) {
			// create new newUser
			user = domain.NewUser(userInfo.ID, userInfo.Name)
			if err := r.CreateUser(ctx, user); err != nil {
				return err
			}
		}

		// create new session to prevent session fixation
		sessionID, err := h.sessionManager.setSessionID(c, 7*24*time.Hour) // max age 1 week
		if err != nil {
			return err
		}
		session := domain.NewSession(sessionID, user.ID, time.Now().Add(7*24*time.Hour))
		err = r.CreateSession(ctx, session)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) Logout(c echo.Context) error {
	ctx := c.Request().Context()

	sessID, err := h.sessionManager.getSessionID(c)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	// delete sessionID
	err = h.repo.DeleteSession(ctx, sessID)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	if err := h.sessionManager.clearSessionID(c); err != nil {
		return internalServerErrorResponse(c, err)
	}

	return c.NoContent(http.StatusOK)
}
