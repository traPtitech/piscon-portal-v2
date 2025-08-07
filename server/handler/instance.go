package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) GetTeamInstances(c echo.Context) error {
	ctx := c.Request().Context()
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, err.Error())
	}

	instances, err := h.useCase.GetTeamInstances(ctx, teamID)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}
	res := lo.Map(instances, toOpenAPIInstance)

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateTeamInstance(c echo.Context) error {
	ctx := c.Request().Context()
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, err.Error())
	}

	instance, err := h.useCase.CreateInstance(ctx, teamID)
	if err != nil {
		if usecase.IsUseCaseError(err) {
			return badRequestResponse(c, err.Error())
		}
		return internalServerErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, toOpenAPIInstance(instance, 0))
}

func (h *Handler) DeleteTeamInstance(c echo.Context) error {
	ctx := c.Request().Context()

	instanceID, err := uuid.Parse(c.Param("instanceID"))
	if err != nil {
		return badRequestResponse(c, "invalid instance id")
	}
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, "invalid team id")
	}
	instance, err := h.useCase.GetInstance(ctx, instanceID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}
	if instance.TeamID != teamID {
		return notFoundResponse(c)
	}

	err = h.useCase.DeleteInstance(ctx, instanceID)
	if err != nil {
		if usecase.IsUseCaseError(err) {
			return badRequestResponse(c, err.Error())
		} else if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) PatchTeamInstance(c echo.Context) error {
	ctx := c.Request().Context()
	instanceID, err := uuid.Parse(c.Param("instanceID"))
	if err != nil {
		return badRequestResponse(c, "invalid instance id")
	}
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, "invalid team id")
	}
	var req openapi.PatchTeamInstanceReq
	if err := c.Bind(&req); err != nil {
		return badRequestResponse(c, err.Error())
	}

	var op domain.InstanceOperation
	switch string(req.Operation) {
	case "start":
		op = domain.InstanceOperationStart
	case "stop":
		op = domain.InstanceOperationStop
	default:
		return badRequestResponse(c, "invalid operation")
	}

	instance, err := h.useCase.GetInstance(ctx, instanceID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}
	if instance.TeamID != teamID {
		return notFoundResponse(c)
	}

	err = h.useCase.UpdateInstance(ctx, instanceID, op)
	if err != nil {
		if usecase.IsUseCaseError(err) {
			return badRequestResponse(c, err.Error())
		}
		if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetInstances(c echo.Context) error {
	ctx := c.Request().Context()
	instances, err := h.useCase.GetAllInstances(ctx)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}
	res := lo.Map(instances, toOpenAPIInstance)
	return c.JSON(http.StatusOK, res)
}

func toOpenAPIInstance(instance domain.Instance, _ int) openapi.Instance {
	privateIP, publicIP := "", ""
	if instance.Infra.PrivateIP != nil {
		privateIP = *instance.Infra.PrivateIP
	}
	if instance.Infra.PublicIP != nil {
		publicIP = *instance.Infra.PublicIP
	}
	return openapi.Instance{
		ID:               openapi.InstanceId(instance.ID),
		TeamId:           openapi.TeamId(instance.TeamID),
		ServerId:         instance.Index,
		PublicIPAddress:  openapi.IPAddress(publicIP),
		PrivateIPAddress: openapi.IPAddress(privateIP),
		Status:           openapi.InstanceStatus(instance.Infra.Status),
		CreatedAt:        instance.CreatedAt,
	}
}
