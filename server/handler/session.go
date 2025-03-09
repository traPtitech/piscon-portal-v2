package handler

import (
	"time"

	"github.com/labstack/echo/v4"
)

//go:generate go tool mockgen -source=$GOFILE -destination=internal/mock/$GOFILE -package=mock -typed=true
type SessionManager interface {
	Init(e *echo.Group)

	GetSessionID(c echo.Context) (string, error)
	SetSessionID(c echo.Context, maxAge time.Duration) (string, error)
	ClearSessionID(c echo.Context) error
}
