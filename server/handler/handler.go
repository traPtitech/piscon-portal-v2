package handler

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/handler/internal"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/oauth2"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

type Handler struct {
	repo           repository.Repository
	sessionManager SessionManager
	oauth2Service  *oauth2.Service

	useCase usecase.UseCase
}

type Config struct {
	RootURL       string
	Debug         bool
	SessionSecret string
	Oauth2        oauth2.Config
	// if SessionManager is nil, use default session manager
	SessionManager SessionManager
}

func New(useCase usecase.UseCase, repo repository.Repository, config Config) (*Handler, error) {
	var sessionManager SessionManager
	if config.SessionManager == nil {
		sessionManager = internal.NewSessionManager(config.SessionSecret, config.Debug)
	} else {
		sessionManager = config.SessionManager
	}

	redirectURI, err := url.JoinPath(config.RootURL, "/api/oauth2/callback")
	if err != nil {
		return nil, fmt.Errorf("join callback URL: %w", err)
	}
	oauth2Service, err := oauth2.NewService(config.Oauth2, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("create oauth2 service: %w", err)
	}
	return &Handler{
		repo:           repo,
		sessionManager: sessionManager,
		oauth2Service:  oauth2Service,
		useCase:        useCase,
	}, nil
}

func (h *Handler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	h.sessionManager.Init(api)

	api.GET("/oauth2/code", h.GetOauth2Code)
	api.GET("/oauth2/callback", h.Oauth2Callback)
	api.POST("/oauth2/logout", h.Logout, h.AuthMiddleware())

	users := api.Group("/users", h.AuthMiddleware())
	users.GET("", h.GetUsers)
	users.GET("/me", h.GetUserMe)

	teams := api.Group("/teams", h.AuthMiddleware())
	teams.GET("", h.GetTeams)
	teams.POST("", h.CreateTeam)
	teams.GET("/:teamID", h.GetTeam)
	// TODO: Admins can access even if they are not members of the team.
	teams.PATCH("/:teamID", h.UpdateTeam, h.TeamAuthMiddleware())
}
