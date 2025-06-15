package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) GetDocument(c echo.Context) error {
	doc, err := h.useCase.GetDocument(c.Request().Context())
	if errors.Is(err, usecase.ErrNotFound) {
		body := openapi.GetDocsOK{
			Body: openapi.OptMarkdownDocument{},
		}
		// 標準の JSON エンコーダでは、空の Opt* 型のエンコードがうまくいかない。
		// この場合には、空のバイト列が返されるため、ogen により生成されたコードを使ってエンコードする。
		// https://github.com/ogen-go/ogen/issues/660
		bodyJSON, err := body.MarshalJSON()
		if err != nil {
			return internalServerErrorResponse(c, err)
		}
		// 初期状態などドキュメントがないのが正しい時もあるので、bodyフィールドが空の状態で200を返す。
		return c.JSONBlob(http.StatusOK, bodyJSON)
	}
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	body := openapi.GetDocsOK{
		Body: openapi.NewOptMarkdownDocument(openapi.MarkdownDocument(doc.Body)),
	}
	return c.JSON(http.StatusOK, body)
}

func (h *Handler) PatchDocument(c echo.Context) error {
	var body openapi.PatchDocsReq
	if err := c.Bind(&body); err != nil {
		return badRequestResponse(c, err.Error())
	}

	doc, err := h.useCase.CreateDocument(c.Request().Context(), string(body.Body))
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := openapi.PatchDocsOK{
		Body: openapi.NewOptMarkdownDocument(openapi.MarkdownDocument(doc.Body)),
	}

	return c.JSON(http.StatusOK, res)
}
