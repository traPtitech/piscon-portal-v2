package handler

import (
	"database/sql"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stephenafamo/bob"
)

type Handler struct {
	db             bob.DB
	sessionManager *sessionManager
	oauth2Service  *oauth2Service
}

type Config struct {
	RootURL       string
	Debug         bool
	SessionSecret string
	Oauth2        Oauth2Config
}

type Oauth2Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
}

func New(db *sql.DB, config Config) (*Handler, error) {
	bobDB := bob.NewDB(db)
	sessionManager := newSessionManager(config.SessionSecret, config.Debug)

	oauth2Service, err := newOauth2Service(config.Oauth2, strings.TrimSuffix(config.RootURL, "/")+"/api/oauth2/callback")
	if err != nil {
		return nil, err
	}
	return &Handler{
		db:             bobDB,
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
