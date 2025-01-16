package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func newJSONRequest(method string, path string, req json.Marshaler) *http.Request {
	body, _ := req.MarshalJSON()
	httpReq := httptest.NewRequest(method, path, bytes.NewReader(body))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return httpReq
}
