package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func (h *Handler) GetBenchmark(c echo.Context) error {
	benchmarkID, err := uuid.Parse(c.Param("benchmarkID"))
	if err != nil {
		return badRequestResponse(c, "invalid benchmark id")
	}

	benchmark, err := h.useCase.GetBenchmark(c.Request().Context(), benchmarkID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}
	log, err := h.useCase.GetBenchmarkLog(c.Request().Context(), benchmarkID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			// log has not been created yet
			// no need to return 404
		} else {
			return internalServerErrorResponse(c, err)
		}
	}

	res := toOpenAPIBenchmarkAdminResult(benchmark, log)

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) EnqueueBenchmark(c echo.Context) error {
	userID := getUserIDFromSession(c)

	var req openapi.PostBenchmarkReq
	if err := c.Bind(&req); err != nil {
		log.Println(err)
		return badRequestResponse(c, err.Error())
	}

	benchmark, err := h.useCase.CreateBenchmark(c.Request().Context(), uuid.UUID(req.InstanceId), userID)
	if err != nil {
		if usecase.IsUseCaseError(err) {
			return badRequestResponse(c, err.Error())
		}
		return internalServerErrorResponse(c, err)
	}

	res := toOpenAPIBenchmarkListItem(benchmark)

	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetBenchmarks(c echo.Context) error {
	benchmarks, err := h.useCase.GetBenchmarks(c.Request().Context())
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]*openapi.BenchmarkListItem, 0, len(benchmarks))
	for _, benchmark := range benchmarks {
		res = append(res, &openapi.BenchmarkListItem{OneOf: toOpenAPIBenchmarkListItem(benchmark)})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetQueuedBenchmarks(c echo.Context) error {
	benchmarks, err := h.useCase.GetQueuedBenchmarks(c.Request().Context())
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]*openapi.BenchmarkListItem, 0, len(benchmarks))
	for _, benchmark := range benchmarks {
		res = append(res, &openapi.BenchmarkListItem{OneOf: toOpenAPIBenchmarkListItem(benchmark)})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllTeamBenchmarks(c echo.Context) error {
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, "invalid team id")
	}

	benchmarks, err := h.useCase.GetTeamBenchmarks(c.Request().Context(), teamID)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]*openapi.BenchmarkListItem, 0, len(benchmarks))
	for _, benchmark := range benchmarks {
		res = append(res, &openapi.BenchmarkListItem{OneOf: toOpenAPIBenchmarkListItem(benchmark)})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetTeamBenchmarks(c echo.Context) error {
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, "invalid team id")
	}

	benchmarks, err := h.useCase.GetTeamBenchmarks(c.Request().Context(), teamID)
	if err != nil {
		return internalServerErrorResponse(c, err)
	}

	res := make([]*openapi.BenchmarkListItem, 0, len(benchmarks))
	for _, benchmark := range benchmarks {
		res = append(res, &openapi.BenchmarkListItem{OneOf: toOpenAPIBenchmarkListItem(benchmark)})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetTeamBenchmark(c echo.Context) error {
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		return badRequestResponse(c, "invalid team id")
	}
	benchmarkID, err := uuid.Parse(c.Param("benchmarkID"))
	if err != nil {
		return badRequestResponse(c, "invalid benchmark id")
	}

	benchmark, err := h.useCase.GetBenchmark(c.Request().Context(), benchmarkID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return notFoundResponse(c)
		}
		return internalServerErrorResponse(c, err)
	}
	if benchmark.TeamID != teamID {
		return notFoundResponse(c)
	}

	log, err := h.useCase.GetBenchmarkLog(c.Request().Context(), benchmarkID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			// log has not been created yet
			// no need to return 404
		} else {
			return internalServerErrorResponse(c, err)
		}
	}

	res := toOpenAPIBenchmark(benchmark, log)

	return c.JSON(http.StatusOK, res)
}

func toOpenAPIBenchmarkListItem(benchmark domain.Benchmark) openapi.BenchmarkListItemSum {
	switch benchmark.Status {
	case domain.BenchmarkStatusWaiting:
		return openapi.NewWaitingBenchmarkBenchmarkListItemSum(openapi.WaitingBenchmark{
			ID:         openapi.BenchmarkId(benchmark.ID),
			InstanceId: openapi.InstanceId(benchmark.Instance.ID),
			TeamId:     openapi.TeamId(benchmark.TeamID),
			UserId:     openapi.UserId(benchmark.UserID),
			Status:     openapi.WaitingBenchmarkStatusWaiting,
			CreatedAt:  openapi.CreatedAt(benchmark.CreatedAt),
		})
	case domain.BenchmarkStatusRunning:
		return openapi.NewRunningBenchmarkBenchmarkListItemSum(openapi.RunningBenchmark{
			ID:         openapi.BenchmarkId(benchmark.ID),
			InstanceId: openapi.InstanceId(benchmark.Instance.ID),
			TeamId:     openapi.TeamId(benchmark.TeamID),
			UserId:     openapi.UserId(benchmark.UserID),
			Status:     openapi.RunningBenchmarkStatusRunning,
			CreatedAt:  openapi.CreatedAt(benchmark.CreatedAt),
			StartedAt:  openapi.StartedAt(*benchmark.StartedAt),
			Score:      openapi.Score(benchmark.Score),
		})
	case domain.BenchmarkStatusFinished:
		return openapi.NewFinishedBenchmarkBenchmarkListItemSum(openapi.FinishedBenchmark{
			ID:         openapi.BenchmarkId(benchmark.ID),
			InstanceId: openapi.InstanceId(benchmark.Instance.ID),
			TeamId:     openapi.TeamId(benchmark.TeamID),
			UserId:     openapi.UserId(benchmark.UserID),
			Status:     openapi.FinishedBenchmarkStatusFinished,
			CreatedAt:  openapi.CreatedAt(benchmark.CreatedAt),
			FinishedAt: openapi.FinishedAt(*benchmark.FinishedAt),
			Score:      openapi.Score(benchmark.Score),
			Result:     toOpenAPIBenchmarkResult(*benchmark.Result),
		})
	default:
		// unreachable
		panic(fmt.Sprintf("unexpected status: %v", benchmark.Status))
	}
}

func toOpenAPIBenchmark(benchmark domain.Benchmark, log domain.BenchmarkLog) *openapi.Benchmark {
	listItem := toOpenAPIBenchmarkListItem(benchmark)
	return &openapi.Benchmark{
		Log: log.UserLog,
		OneOf: openapi.BenchmarkSum{
			Type:              openapi.BenchmarkSumType(listItem.Type),
			WaitingBenchmark:  listItem.WaitingBenchmark,
			RunningBenchmark:  listItem.RunningBenchmark,
			FinishedBenchmark: listItem.FinishedBenchmark,
		},
	}
}

func toOpenAPIBenchmarkResult(result domain.BenchmarkResult) openapi.FinishedBenchmarkResult {
	switch result {
	case domain.BenchmarkResultStatusError:
		return openapi.FinishedBenchmarkResultError
	case domain.BenchmarkResultStatusPassed:
		return openapi.FinishedBenchmarkResultPassed
	case domain.BenchmarkResultStatusFailed:
		return openapi.FinishedBenchmarkResultFailed
	default:
		// unreachable
		panic(fmt.Sprintf("unexpected result: %v", result))
	}
}

func toOpenAPIBenchmarkAdminResult(benchmark domain.Benchmark, log domain.BenchmarkLog) *openapi.BenchmarkAdminResult {
	listItem := toOpenAPIBenchmarkListItem(benchmark)
	return &openapi.BenchmarkAdminResult{
		Log:      log.UserLog,
		AdminLog: log.AdminLog,
		OneOf: openapi.BenchmarkAdminResultSum{
			Type:              openapi.BenchmarkAdminResultSumType(listItem.Type),
			WaitingBenchmark:  listItem.WaitingBenchmark,
			RunningBenchmark:  listItem.RunningBenchmark,
			FinishedBenchmark: listItem.FinishedBenchmark,
		},
	}
}
