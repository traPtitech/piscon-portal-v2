package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
)

func internalServerErrorResponse(c echo.Context, err error) error {
	c.Logger().Error(err)
	return c.JSON(http.StatusInternalServerError, openapi.InternalServerError{
		Message: openapi.NewOptString("Internal Server Error"),
	})
}

func badRequestResponse(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, openapi.ErrorBadRequest{
		Message: openapi.NewOptString(msg),
	})
}

func unauthorizedResponse(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, openapi.Unauthorized{
		Message: openapi.NewOptString(msg),
	})
}

func notFoundResponse(c echo.Context) error {
	return c.JSON(http.StatusNotFound, openapi.NotFound{
		Message: openapi.NewOptString("Not Found"),
	})
}
