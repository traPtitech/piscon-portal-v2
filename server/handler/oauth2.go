package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/models"
	"golang.org/x/oauth2"
)

func (h *Handler) GetOauth2Code(c echo.Context) error {
	sessID, err := h.sessionManager.setSessionID(c, 7*24*time.Hour) // max age 1 week
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	state := generateRandomString(16)
	url := h.oauth2Service.authCodeURL(sessID, state)

	return c.Redirect(http.StatusSeeOther, url)
}

func (h *Handler) Oauth2Callback(c echo.Context) error {
	ctx := c.Request().Context()

	sessionID, err := h.sessionManager.getSessionID(c)
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if !h.oauth2Service.verifyState(sessionID, state) {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	token, err := h.oauth2Service.exchange(ctx, sessionID, code)
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}
	userName, err := h.oauth2Service.getUserName(ctx, token)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Check if user exists
	user, err := models.Users.Query(models.SelectWhere.Users.Name.EQ(userName)).One(ctx, h.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return internalServerErrorResponse(c, err.Error())
	}
	var userID string
	if errors.Is(err, sql.ErrNoRows) {
		// create new newUser
		newUser := domain.NewUser(userName)
		_, err := models.Users.Insert(&models.UserSetter{
			ID:   omit.From(newUser.ID),
			Name: omit.From(newUser.Name),
		}).Exec(ctx, h.db)
		if err != nil {
			return internalServerErrorResponse(c, err.Error())
		}
		userID = newUser.ID
	} else {
		userID = user.ID
	}
	// create new session
	_, err = models.Sessions.Insert(&models.SessionSetter{
		ID:        omit.From(sessionID),
		UserID:    omit.From(userID),
		ExpiredAt: omit.From(time.Now().Add(7 * 24 * time.Hour)), // 1 week
	}).Exec(ctx, h.db)
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) Logout(c echo.Context) error {
	sessID, err := h.sessionManager.getSessionID(c)
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	// delete sessionID
	_, err = models.Sessions.Delete(models.DeleteWhere.Sessions.ID.EQ(sessID)).Exec(c.Request().Context(), h.db)
	if err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	if err := h.sessionManager.clearSessionID(c); err != nil {
		return internalServerErrorResponse(c, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

type oauth2Service struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier

	stateStorage        *ttlcache.Cache[string, string]
	codeVerifierStorage *ttlcache.Cache[string, string]
}

func newOauth2Service(config Oauth2Config, redirectURL string) (*oauth2Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"openid", "profile"},
	}

	provider, err := oidc.NewProvider(context.Background(), config.Issuer)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	stateStorage := ttlcache.New[string, string]()
	codeVerifierStorage := ttlcache.New[string, string]()
	go stateStorage.Start()
	go codeVerifierStorage.Start()

	return &oauth2Service{
		config:   oauth2Config,
		verifier: verifier,

		stateStorage:        stateStorage,
		codeVerifierStorage: codeVerifierStorage,
	}, nil
}

func (s *oauth2Service) authCodeURL(sessionID string, state string) string {
	s.stateStorage.Set(sessionID, state, 15*time.Minute)

	codeVerifier := oauth2.GenerateVerifier()
	s.codeVerifierStorage.Set(sessionID, codeVerifier, 15*time.Minute)

	return s.config.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
}

func (s *oauth2Service) verifyState(sessionID, state string) bool {
	storedState, found := s.stateStorage.GetAndDelete(sessionID)
	if !found {
		return false
	}
	return storedState.Value() == state
}

func (s *oauth2Service) exchange(ctx context.Context, sessionID, code string) (*oauth2.Token, error) {
	codeVerifier, found := s.codeVerifierStorage.GetAndDelete(sessionID)
	if !found {
		return nil, errors.New("code verifier not found")
	}
	return s.config.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier.Value()))
}

func (s *oauth2Service) getUserName(ctx context.Context, token *oauth2.Token) (string, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("missing id_token")
	}
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	// traQ returns username in "name" claim
	// ref: https://github.com/traPtitech/traQ/blob/v3.21.0/service/oidc/userinfo.go#L57
	var info struct {
		Name string `json:"name"`
	}
	if err := idToken.Claims(&info); err != nil {
		return "", err
	}

	return info.Name, nil
}
