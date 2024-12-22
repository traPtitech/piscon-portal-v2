package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/oauth2"
)

type Handler struct {
	repo           repository.Repository
	sessionManager *sessionManager
	oauth2Service  *oauth2.Service
}

type Config struct {
	RootURL       string
	Debug         bool
	SessionSecret string
	Oauth2        oauth2.Config
}

func New(repo repository.Repository, config Config) (*Handler, error) {
	sessionManager := newSessionManager(config.SessionSecret, config.Debug)

	oauth2Service, err := oauth2.NewService(config.Oauth2, strings.TrimSuffix(config.RootURL, "/")+"/api/oauth2/callback")
	if err != nil {
		return nil, err
	}
	return &Handler{
		repo:           repo,
		sessionManager: sessionManager,
		oauth2Service:  oauth2Service,
	}, nil
}

func (h *Handler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	h.sessionManager.init(api)

	api.GET("/oauth2/code", h.GetOauth2Code)
	api.GET("/oauth2/callback", h.Oauth2Callback)
	api.POST("/oauth2/logout", h.Logout, h.AuthMiddleware())
}
