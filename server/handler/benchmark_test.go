package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"go.uber.org/mock/gomock"
)

func TestGetBenchmark(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	benchmarkID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/benchmarks/"+benchmarkID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("benchmarkID")
	c.SetParamValues(benchmarkID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmark := domain.Benchmark{
		ID:         benchmarkID,
		TeamID:     uuid.New(),
		UserID:     uuid.New(),
		Status:     domain.BenchmarkStatusFinished,
		CreatedAt:  time.Now(),
		FinishedAt: ptr.Of(time.Now()),
		Score:      100,
		Result:     ptr.Of(domain.BenchmarkResultStatusPassed),
	}

	log := domain.BenchmarkLog{
		UserLog:  "user log",
		AdminLog: "admin log",
	}

	useCaseMock.EXPECT().GetBenchmark(gomock.Any(), benchmarkID).Return(benchmark, nil)
	useCaseMock.EXPECT().GetBenchmarkLog(gomock.Any(), benchmarkID).Return(log, nil)

	_ = h.GetBenchmark(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.BenchmarkAdminResult
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, benchmarkID, uuid.UUID(res.OneOf.FinishedBenchmark.ID))
}

func TestGetBenchmark_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	benchmarkID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/benchmarks/"+benchmarkID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("benchmarkID")
	c.SetParamValues(benchmarkID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	useCaseMock.EXPECT().GetBenchmark(gomock.Any(), benchmarkID).Return(domain.Benchmark{}, usecase.ErrNotFound)

	_ = h.GetBenchmark(c)

	if !assert.Equal(t, http.StatusNotFound, rec.Code) {
		t.Log(rec.Body.String())
	}
}

func TestEnqueueBenchmark(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	userID := uuid.New()

	e := echo.New()
	req := &openapi.PostBenchmarkReq{
		InstanceId: openapi.InstanceId(uuid.New()),
	}
	httpReq := newJSONRequest(http.MethodPost, "/benchmarks", req)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.Set(handler.UserIDKey, userID)
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmarkID := uuid.New()
	useCaseMock.EXPECT().CreateBenchmark(gomock.Any(), uuid.UUID(req.InstanceId), userID).Return(domain.Benchmark{
		ID:        benchmarkID,
		Instance:  domain.Instance{ID: uuid.UUID(req.InstanceId)},
		TeamID:    uuid.New(),
		UserID:    userID,
		Status:    domain.BenchmarkStatusWaiting,
		CreatedAt: time.Now(),
	}, nil)

	_ = h.EnqueueBenchmark(c)

	if !assert.Equal(t, http.StatusCreated, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.BenchmarkListItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, benchmarkID, uuid.UUID(res.OneOf.WaitingBenchmark.ID))
}

func TestGetBenchmarks(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/benchmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmarks := []domain.Benchmark{
		{
			ID:        uuid.New(),
			TeamID:    uuid.New(),
			UserID:    uuid.New(),
			Status:    domain.BenchmarkStatusWaiting,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			TeamID:    uuid.New(),
			UserID:    uuid.New(),
			Status:    domain.BenchmarkStatusRunning,
			CreatedAt: time.Now(),
			StartedAt: ptr.Of(time.Now()),
			Score:     50,
		},
	}

	useCaseMock.EXPECT().GetBenchmarks(gomock.Any()).Return(benchmarks, nil)

	_ = h.GetBenchmarks(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res []*openapi.BenchmarkListItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	compareBenchmarks(t, benchmarks, res)
}

func TestGetQueuedBenchmarks(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/benchmarks/queued", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmarks := []domain.Benchmark{
		{
			ID:        uuid.New(),
			TeamID:    uuid.New(),
			UserID:    uuid.New(),
			Status:    domain.BenchmarkStatusWaiting,
			CreatedAt: time.Now(),
		},
	}

	useCaseMock.EXPECT().GetQueuedBenchmarks(gomock.Any()).Return(benchmarks, nil)

	_ = h.GetQueuedBenchmarks(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res []*openapi.BenchmarkListItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	compareBenchmarks(t, benchmarks, res)
}

func TestGetAllTeamBenchmarks(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/teams/"+teamID.String()+"/benchmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamID")
	c.SetParamValues(teamID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmarks := []domain.Benchmark{
		{
			ID:        uuid.New(),
			TeamID:    teamID,
			UserID:    uuid.New(),
			Status:    domain.BenchmarkStatusWaiting,
			CreatedAt: time.Now(),
		},
	}

	useCaseMock.EXPECT().GetTeamBenchmarks(gomock.Any(), teamID).Return(benchmarks, nil)

	_ = h.GetAllTeamBenchmarks(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res []*openapi.BenchmarkListItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	compareBenchmarks(t, benchmarks, res)
}

func TestGetTeamBenchmark(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := repomock.NewMockRepository(ctrl)
	useCaseMock := usecasemock.NewMockUseCase(ctrl)

	e := echo.New()
	teamID := uuid.New()
	benchmarkID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/teams/"+teamID.String()+"/benchmarks/"+benchmarkID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamID", "benchmarkID")
	c.SetParamValues(teamID.String(), benchmarkID.String())
	h := NewHandler(useCaseMock, repoMock, nil)

	benchmark := domain.Benchmark{
		ID:         benchmarkID,
		TeamID:     teamID,
		UserID:     uuid.New(),
		Status:     domain.BenchmarkStatusFinished,
		CreatedAt:  time.Now(),
		FinishedAt: ptr.Of(time.Now()),
		Score:      100,
		Result:     ptr.Of(domain.BenchmarkResultStatusPassed),
	}

	log := domain.BenchmarkLog{
		UserLog:  "user log",
		AdminLog: "admin log",
	}

	useCaseMock.EXPECT().GetBenchmark(gomock.Any(), benchmarkID).Return(benchmark, nil)
	useCaseMock.EXPECT().GetBenchmarkLog(gomock.Any(), benchmarkID).Return(log, nil)

	_ = h.GetTeamBenchmark(c)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
	var res openapi.Benchmark
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, benchmarkID, uuid.UUID(res.OneOf.FinishedBenchmark.ID))
}

func compareWaitingBenchmark(t *testing.T, expected domain.Benchmark, actual openapi.WaitingBenchmark) {
	t.Helper()
	assert.Equal(t, expected.ID, uuid.UUID(actual.ID))
	assert.Equal(t, expected.Instance.ID, uuid.UUID(actual.InstanceId))
	assert.Equal(t, expected.TeamID, uuid.UUID(actual.TeamId))
	assert.Equal(t, expected.UserID, uuid.UUID(actual.UserId))
	assert.Equal(t, string(expected.Status), string(actual.Status))
	assert.WithinDuration(t, expected.CreatedAt, time.Time(actual.CreatedAt), time.Second)
}

func compareRunningBenchmark(t *testing.T, expected domain.Benchmark, actual openapi.RunningBenchmark) {
	t.Helper()
	assert.Equal(t, expected.ID, uuid.UUID(actual.ID))
	assert.Equal(t, expected.Instance.ID, uuid.UUID(actual.InstanceId))
	assert.Equal(t, expected.TeamID, uuid.UUID(actual.TeamId))
	assert.Equal(t, expected.UserID, uuid.UUID(actual.UserId))
	assert.Equal(t, string(expected.Status), string(actual.Status))
	assert.WithinDuration(t, expected.CreatedAt, time.Time(actual.CreatedAt), time.Second)
	assert.WithinDuration(t, *expected.StartedAt, time.Time(actual.StartedAt), time.Second)
	assert.Equal(t, expected.Score, int64(actual.Score))
}

func compareFinishedBenchmark(t *testing.T, expected domain.Benchmark, actual openapi.FinishedBenchmark) {
	t.Helper()
	assert.Equal(t, expected.ID, uuid.UUID(actual.ID))
	assert.Equal(t, expected.Instance.ID, uuid.UUID(actual.InstanceId))
	assert.Equal(t, expected.TeamID, uuid.UUID(actual.TeamId))
	assert.Equal(t, expected.UserID, uuid.UUID(actual.UserId))
	assert.Equal(t, string(expected.Status), string(actual.Status))
	assert.WithinDuration(t, expected.CreatedAt, time.Time(actual.CreatedAt), time.Second)
	assert.WithinDuration(t, *expected.FinishedAt, time.Time(actual.FinishedAt), time.Second)
	assert.Equal(t, expected.Score, int64(actual.Score))
	assert.Equal(t, string(*expected.Result), string(actual.Result))
}

func compareBenchmarks(t *testing.T, expected []domain.Benchmark, actual []*openapi.BenchmarkListItem) {
	t.Helper()
	if !assert.Len(t, actual, len(expected)) {
		return
	}
	for i, b := range expected {
		switch b.Status {
		case domain.BenchmarkStatusWaiting:
			compareWaitingBenchmark(t, b, actual[i].OneOf.WaitingBenchmark)
		case domain.BenchmarkStatusRunning:
			compareRunningBenchmark(t, b, actual[i].OneOf.RunningBenchmark)
		case domain.BenchmarkStatusFinished:
			compareFinishedBenchmark(t, b, actual[i].OneOf.FinishedBenchmark)
		}
	}
}
