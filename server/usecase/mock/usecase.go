// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mock/usecase.go -package=mock -typed=true
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	domain "github.com/traPtitech/piscon-portal-v2/server/domain"
	usecase "github.com/traPtitech/piscon-portal-v2/server/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
	isgomock struct{}
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// CreateBenchmark mocks base method.
func (m *MockUseCase) CreateBenchmark(ctx context.Context, instanceID, userID uuid.UUID) (domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBenchmark", ctx, instanceID, userID)
	ret0, _ := ret[0].(domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBenchmark indicates an expected call of CreateBenchmark.
func (mr *MockUseCaseMockRecorder) CreateBenchmark(ctx, instanceID, userID any) *MockUseCaseCreateBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBenchmark", reflect.TypeOf((*MockUseCase)(nil).CreateBenchmark), ctx, instanceID, userID)
	return &MockUseCaseCreateBenchmarkCall{Call: call}
}

// MockUseCaseCreateBenchmarkCall wrap *gomock.Call
type MockUseCaseCreateBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseCreateBenchmarkCall) Return(arg0 domain.Benchmark, arg1 error) *MockUseCaseCreateBenchmarkCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseCreateBenchmarkCall) Do(f func(context.Context, uuid.UUID, uuid.UUID) (domain.Benchmark, error)) *MockUseCaseCreateBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseCreateBenchmarkCall) DoAndReturn(f func(context.Context, uuid.UUID, uuid.UUID) (domain.Benchmark, error)) *MockUseCaseCreateBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CreateTeam mocks base method.
func (m *MockUseCase) CreateTeam(ctx context.Context, input usecase.CreateTeamInput) (domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTeam", ctx, input)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTeam indicates an expected call of CreateTeam.
func (mr *MockUseCaseMockRecorder) CreateTeam(ctx, input any) *MockUseCaseCreateTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTeam", reflect.TypeOf((*MockUseCase)(nil).CreateTeam), ctx, input)
	return &MockUseCaseCreateTeamCall{Call: call}
}

// MockUseCaseCreateTeamCall wrap *gomock.Call
type MockUseCaseCreateTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseCreateTeamCall) Return(arg0 domain.Team, arg1 error) *MockUseCaseCreateTeamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseCreateTeamCall) Do(f func(context.Context, usecase.CreateTeamInput) (domain.Team, error)) *MockUseCaseCreateTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseCreateTeamCall) DoAndReturn(f func(context.Context, usecase.CreateTeamInput) (domain.Team, error)) *MockUseCaseCreateTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FinalizeBenchmark mocks base method.
func (m *MockUseCase) FinalizeBenchmark(ctx context.Context, benchmarkID uuid.UUID, result domain.BenchmarkResult, finishedAt time.Time, errorMessage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeBenchmark", ctx, benchmarkID, result, finishedAt, errorMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeBenchmark indicates an expected call of FinalizeBenchmark.
func (mr *MockUseCaseMockRecorder) FinalizeBenchmark(ctx, benchmarkID, result, finishedAt, errorMessage any) *MockUseCaseFinalizeBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeBenchmark", reflect.TypeOf((*MockUseCase)(nil).FinalizeBenchmark), ctx, benchmarkID, result, finishedAt, errorMessage)
	return &MockUseCaseFinalizeBenchmarkCall{Call: call}
}

// MockUseCaseFinalizeBenchmarkCall wrap *gomock.Call
type MockUseCaseFinalizeBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseFinalizeBenchmarkCall) Return(arg0 error) *MockUseCaseFinalizeBenchmarkCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseFinalizeBenchmarkCall) Do(f func(context.Context, uuid.UUID, domain.BenchmarkResult, time.Time, string) error) *MockUseCaseFinalizeBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseFinalizeBenchmarkCall) DoAndReturn(f func(context.Context, uuid.UUID, domain.BenchmarkResult, time.Time, string) error) *MockUseCaseFinalizeBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBenchmark mocks base method.
func (m *MockUseCase) GetBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBenchmark", ctx, id)
	ret0, _ := ret[0].(domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBenchmark indicates an expected call of GetBenchmark.
func (mr *MockUseCaseMockRecorder) GetBenchmark(ctx, id any) *MockUseCaseGetBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBenchmark", reflect.TypeOf((*MockUseCase)(nil).GetBenchmark), ctx, id)
	return &MockUseCaseGetBenchmarkCall{Call: call}
}

// MockUseCaseGetBenchmarkCall wrap *gomock.Call
type MockUseCaseGetBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetBenchmarkCall) Return(arg0 domain.Benchmark, arg1 error) *MockUseCaseGetBenchmarkCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetBenchmarkCall) Do(f func(context.Context, uuid.UUID) (domain.Benchmark, error)) *MockUseCaseGetBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetBenchmarkCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.Benchmark, error)) *MockUseCaseGetBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBenchmarkLog mocks base method.
func (m *MockUseCase) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBenchmarkLog", ctx, benchmarkID)
	ret0, _ := ret[0].(domain.BenchmarkLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBenchmarkLog indicates an expected call of GetBenchmarkLog.
func (mr *MockUseCaseMockRecorder) GetBenchmarkLog(ctx, benchmarkID any) *MockUseCaseGetBenchmarkLogCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBenchmarkLog", reflect.TypeOf((*MockUseCase)(nil).GetBenchmarkLog), ctx, benchmarkID)
	return &MockUseCaseGetBenchmarkLogCall{Call: call}
}

// MockUseCaseGetBenchmarkLogCall wrap *gomock.Call
type MockUseCaseGetBenchmarkLogCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetBenchmarkLogCall) Return(arg0 domain.BenchmarkLog, arg1 error) *MockUseCaseGetBenchmarkLogCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetBenchmarkLogCall) Do(f func(context.Context, uuid.UUID) (domain.BenchmarkLog, error)) *MockUseCaseGetBenchmarkLogCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetBenchmarkLogCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.BenchmarkLog, error)) *MockUseCaseGetBenchmarkLogCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBenchmarks mocks base method.
func (m *MockUseCase) GetBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBenchmarks", ctx)
	ret0, _ := ret[0].([]domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBenchmarks indicates an expected call of GetBenchmarks.
func (mr *MockUseCaseMockRecorder) GetBenchmarks(ctx any) *MockUseCaseGetBenchmarksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBenchmarks", reflect.TypeOf((*MockUseCase)(nil).GetBenchmarks), ctx)
	return &MockUseCaseGetBenchmarksCall{Call: call}
}

// MockUseCaseGetBenchmarksCall wrap *gomock.Call
type MockUseCaseGetBenchmarksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetBenchmarksCall) Return(arg0 []domain.Benchmark, arg1 error) *MockUseCaseGetBenchmarksCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetBenchmarksCall) Do(f func(context.Context) ([]domain.Benchmark, error)) *MockUseCaseGetBenchmarksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetBenchmarksCall) DoAndReturn(f func(context.Context) ([]domain.Benchmark, error)) *MockUseCaseGetBenchmarksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetQueuedBenchmarks mocks base method.
func (m *MockUseCase) GetQueuedBenchmarks(ctx context.Context) ([]domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQueuedBenchmarks", ctx)
	ret0, _ := ret[0].([]domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueuedBenchmarks indicates an expected call of GetQueuedBenchmarks.
func (mr *MockUseCaseMockRecorder) GetQueuedBenchmarks(ctx any) *MockUseCaseGetQueuedBenchmarksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQueuedBenchmarks", reflect.TypeOf((*MockUseCase)(nil).GetQueuedBenchmarks), ctx)
	return &MockUseCaseGetQueuedBenchmarksCall{Call: call}
}

// MockUseCaseGetQueuedBenchmarksCall wrap *gomock.Call
type MockUseCaseGetQueuedBenchmarksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetQueuedBenchmarksCall) Return(arg0 []domain.Benchmark, arg1 error) *MockUseCaseGetQueuedBenchmarksCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetQueuedBenchmarksCall) Do(f func(context.Context) ([]domain.Benchmark, error)) *MockUseCaseGetQueuedBenchmarksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetQueuedBenchmarksCall) DoAndReturn(f func(context.Context) ([]domain.Benchmark, error)) *MockUseCaseGetQueuedBenchmarksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetTeam mocks base method.
func (m *MockUseCase) GetTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeam", ctx, id)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeam indicates an expected call of GetTeam.
func (mr *MockUseCaseMockRecorder) GetTeam(ctx, id any) *MockUseCaseGetTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeam", reflect.TypeOf((*MockUseCase)(nil).GetTeam), ctx, id)
	return &MockUseCaseGetTeamCall{Call: call}
}

// MockUseCaseGetTeamCall wrap *gomock.Call
type MockUseCaseGetTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetTeamCall) Return(arg0 domain.Team, arg1 error) *MockUseCaseGetTeamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetTeamCall) Do(f func(context.Context, uuid.UUID) (domain.Team, error)) *MockUseCaseGetTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetTeamCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.Team, error)) *MockUseCaseGetTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetTeamBenchmarks mocks base method.
func (m *MockUseCase) GetTeamBenchmarks(ctx context.Context, teamID uuid.UUID) ([]domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamBenchmarks", ctx, teamID)
	ret0, _ := ret[0].([]domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamBenchmarks indicates an expected call of GetTeamBenchmarks.
func (mr *MockUseCaseMockRecorder) GetTeamBenchmarks(ctx, teamID any) *MockUseCaseGetTeamBenchmarksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamBenchmarks", reflect.TypeOf((*MockUseCase)(nil).GetTeamBenchmarks), ctx, teamID)
	return &MockUseCaseGetTeamBenchmarksCall{Call: call}
}

// MockUseCaseGetTeamBenchmarksCall wrap *gomock.Call
type MockUseCaseGetTeamBenchmarksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetTeamBenchmarksCall) Return(arg0 []domain.Benchmark, arg1 error) *MockUseCaseGetTeamBenchmarksCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetTeamBenchmarksCall) Do(f func(context.Context, uuid.UUID) ([]domain.Benchmark, error)) *MockUseCaseGetTeamBenchmarksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetTeamBenchmarksCall) DoAndReturn(f func(context.Context, uuid.UUID) ([]domain.Benchmark, error)) *MockUseCaseGetTeamBenchmarksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetTeams mocks base method.
func (m *MockUseCase) GetTeams(ctx context.Context) ([]domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeams", ctx)
	ret0, _ := ret[0].([]domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeams indicates an expected call of GetTeams.
func (mr *MockUseCaseMockRecorder) GetTeams(ctx any) *MockUseCaseGetTeamsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeams", reflect.TypeOf((*MockUseCase)(nil).GetTeams), ctx)
	return &MockUseCaseGetTeamsCall{Call: call}
}

// MockUseCaseGetTeamsCall wrap *gomock.Call
type MockUseCaseGetTeamsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetTeamsCall) Return(arg0 []domain.Team, arg1 error) *MockUseCaseGetTeamsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetTeamsCall) Do(f func(context.Context) ([]domain.Team, error)) *MockUseCaseGetTeamsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetTeamsCall) DoAndReturn(f func(context.Context) ([]domain.Team, error)) *MockUseCaseGetTeamsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUser mocks base method.
func (m *MockUseCase) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUseCaseMockRecorder) GetUser(ctx, userID any) *MockUseCaseGetUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUseCase)(nil).GetUser), ctx, userID)
	return &MockUseCaseGetUserCall{Call: call}
}

// MockUseCaseGetUserCall wrap *gomock.Call
type MockUseCaseGetUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetUserCall) Return(arg0 domain.User, arg1 error) *MockUseCaseGetUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetUserCall) Do(f func(context.Context, uuid.UUID) (domain.User, error)) *MockUseCaseGetUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetUserCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.User, error)) *MockUseCaseGetUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUsers mocks base method.
func (m *MockUseCase) GetUsers(ctx context.Context) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUseCaseMockRecorder) GetUsers(ctx any) *MockUseCaseGetUsersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUseCase)(nil).GetUsers), ctx)
	return &MockUseCaseGetUsersCall{Call: call}
}

// MockUseCaseGetUsersCall wrap *gomock.Call
type MockUseCaseGetUsersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseGetUsersCall) Return(arg0 []domain.User, arg1 error) *MockUseCaseGetUsersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseGetUsersCall) Do(f func(context.Context) ([]domain.User, error)) *MockUseCaseGetUsersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseGetUsersCall) DoAndReturn(f func(context.Context) ([]domain.User, error)) *MockUseCaseGetUsersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SaveBenchmarkProgress mocks base method.
func (m *MockUseCase) SaveBenchmarkProgress(ctx context.Context, benchmarkID uuid.UUID, benchLog domain.BenchmarkLog, score int64, startedAt time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBenchmarkProgress", ctx, benchmarkID, benchLog, score, startedAt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBenchmarkProgress indicates an expected call of SaveBenchmarkProgress.
func (mr *MockUseCaseMockRecorder) SaveBenchmarkProgress(ctx, benchmarkID, benchLog, score, startedAt any) *MockUseCaseSaveBenchmarkProgressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBenchmarkProgress", reflect.TypeOf((*MockUseCase)(nil).SaveBenchmarkProgress), ctx, benchmarkID, benchLog, score, startedAt)
	return &MockUseCaseSaveBenchmarkProgressCall{Call: call}
}

// MockUseCaseSaveBenchmarkProgressCall wrap *gomock.Call
type MockUseCaseSaveBenchmarkProgressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseSaveBenchmarkProgressCall) Return(arg0 error) *MockUseCaseSaveBenchmarkProgressCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseSaveBenchmarkProgressCall) Do(f func(context.Context, uuid.UUID, domain.BenchmarkLog, int64, time.Time) error) *MockUseCaseSaveBenchmarkProgressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseSaveBenchmarkProgressCall) DoAndReturn(f func(context.Context, uuid.UUID, domain.BenchmarkLog, int64, time.Time) error) *MockUseCaseSaveBenchmarkProgressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StartBenchmark mocks base method.
func (m *MockUseCase) StartBenchmark(ctx context.Context) (domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartBenchmark", ctx)
	ret0, _ := ret[0].(domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartBenchmark indicates an expected call of StartBenchmark.
func (mr *MockUseCaseMockRecorder) StartBenchmark(ctx any) *MockUseCaseStartBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartBenchmark", reflect.TypeOf((*MockUseCase)(nil).StartBenchmark), ctx)
	return &MockUseCaseStartBenchmarkCall{Call: call}
}

// MockUseCaseStartBenchmarkCall wrap *gomock.Call
type MockUseCaseStartBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseStartBenchmarkCall) Return(arg0 domain.Benchmark, arg1 error) *MockUseCaseStartBenchmarkCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseStartBenchmarkCall) Do(f func(context.Context) (domain.Benchmark, error)) *MockUseCaseStartBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseStartBenchmarkCall) DoAndReturn(f func(context.Context) (domain.Benchmark, error)) *MockUseCaseStartBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateTeam mocks base method.
func (m *MockUseCase) UpdateTeam(ctx context.Context, input usecase.UpdateTeamInput) (domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTeam", ctx, input)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTeam indicates an expected call of UpdateTeam.
func (mr *MockUseCaseMockRecorder) UpdateTeam(ctx, input any) *MockUseCaseUpdateTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTeam", reflect.TypeOf((*MockUseCase)(nil).UpdateTeam), ctx, input)
	return &MockUseCaseUpdateTeamCall{Call: call}
}

// MockUseCaseUpdateTeamCall wrap *gomock.Call
type MockUseCaseUpdateTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUseCaseUpdateTeamCall) Return(arg0 domain.Team, arg1 error) *MockUseCaseUpdateTeamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUseCaseUpdateTeamCall) Do(f func(context.Context, usecase.UpdateTeamInput) (domain.Team, error)) *MockUseCaseUpdateTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUseCaseUpdateTeamCall) DoAndReturn(f func(context.Context, usecase.UpdateTeamInput) (domain.Team, error)) *MockUseCaseUpdateTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
