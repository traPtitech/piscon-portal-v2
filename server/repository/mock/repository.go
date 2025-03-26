// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=mock/repository.go -package=mock -typed=true
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	domain "github.com/traPtitech/piscon-portal-v2/server/domain"
	repository "github.com/traPtitech/piscon-portal-v2/server/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateBenchmark mocks base method.
func (m *MockRepository) CreateBenchmark(ctx context.Context, benchmark domain.Benchmark) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBenchmark", ctx, benchmark)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBenchmark indicates an expected call of CreateBenchmark.
func (mr *MockRepositoryMockRecorder) CreateBenchmark(ctx, benchmark any) *MockRepositoryCreateBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBenchmark", reflect.TypeOf((*MockRepository)(nil).CreateBenchmark), ctx, benchmark)
	return &MockRepositoryCreateBenchmarkCall{Call: call}
}

// MockRepositoryCreateBenchmarkCall wrap *gomock.Call
type MockRepositoryCreateBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryCreateBenchmarkCall) Return(arg0 error) *MockRepositoryCreateBenchmarkCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryCreateBenchmarkCall) Do(f func(context.Context, domain.Benchmark) error) *MockRepositoryCreateBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryCreateBenchmarkCall) DoAndReturn(f func(context.Context, domain.Benchmark) error) *MockRepositoryCreateBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CreateSession mocks base method.
func (m *MockRepository) CreateSession(ctx context.Context, session domain.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockRepositoryMockRecorder) CreateSession(ctx, session any) *MockRepositoryCreateSessionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockRepository)(nil).CreateSession), ctx, session)
	return &MockRepositoryCreateSessionCall{Call: call}
}

// MockRepositoryCreateSessionCall wrap *gomock.Call
type MockRepositoryCreateSessionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryCreateSessionCall) Return(arg0 error) *MockRepositoryCreateSessionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryCreateSessionCall) Do(f func(context.Context, domain.Session) error) *MockRepositoryCreateSessionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryCreateSessionCall) DoAndReturn(f func(context.Context, domain.Session) error) *MockRepositoryCreateSessionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CreateTeam mocks base method.
func (m *MockRepository) CreateTeam(ctx context.Context, team domain.Team) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTeam", ctx, team)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTeam indicates an expected call of CreateTeam.
func (mr *MockRepositoryMockRecorder) CreateTeam(ctx, team any) *MockRepositoryCreateTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTeam", reflect.TypeOf((*MockRepository)(nil).CreateTeam), ctx, team)
	return &MockRepositoryCreateTeamCall{Call: call}
}

// MockRepositoryCreateTeamCall wrap *gomock.Call
type MockRepositoryCreateTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryCreateTeamCall) Return(arg0 error) *MockRepositoryCreateTeamCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryCreateTeamCall) Do(f func(context.Context, domain.Team) error) *MockRepositoryCreateTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryCreateTeamCall) DoAndReturn(f func(context.Context, domain.Team) error) *MockRepositoryCreateTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user any) *MockRepositoryCreateUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
	return &MockRepositoryCreateUserCall{Call: call}
}

// MockRepositoryCreateUserCall wrap *gomock.Call
type MockRepositoryCreateUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryCreateUserCall) Return(arg0 error) *MockRepositoryCreateUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryCreateUserCall) Do(f func(context.Context, domain.User) error) *MockRepositoryCreateUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryCreateUserCall) DoAndReturn(f func(context.Context, domain.User) error) *MockRepositoryCreateUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteSession mocks base method.
func (m *MockRepository) DeleteSession(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockRepositoryMockRecorder) DeleteSession(ctx, id any) *MockRepositoryDeleteSessionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockRepository)(nil).DeleteSession), ctx, id)
	return &MockRepositoryDeleteSessionCall{Call: call}
}

// MockRepositoryDeleteSessionCall wrap *gomock.Call
type MockRepositoryDeleteSessionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryDeleteSessionCall) Return(arg0 error) *MockRepositoryDeleteSessionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryDeleteSessionCall) Do(f func(context.Context, string) error) *MockRepositoryDeleteSessionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryDeleteSessionCall) DoAndReturn(f func(context.Context, string) error) *MockRepositoryDeleteSessionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindBenchmark mocks base method.
func (m *MockRepository) FindBenchmark(ctx context.Context, id uuid.UUID) (domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBenchmark", ctx, id)
	ret0, _ := ret[0].(domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBenchmark indicates an expected call of FindBenchmark.
func (mr *MockRepositoryMockRecorder) FindBenchmark(ctx, id any) *MockRepositoryFindBenchmarkCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBenchmark", reflect.TypeOf((*MockRepository)(nil).FindBenchmark), ctx, id)
	return &MockRepositoryFindBenchmarkCall{Call: call}
}

// MockRepositoryFindBenchmarkCall wrap *gomock.Call
type MockRepositoryFindBenchmarkCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindBenchmarkCall) Return(arg0 domain.Benchmark, arg1 error) *MockRepositoryFindBenchmarkCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindBenchmarkCall) Do(f func(context.Context, uuid.UUID) (domain.Benchmark, error)) *MockRepositoryFindBenchmarkCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindBenchmarkCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.Benchmark, error)) *MockRepositoryFindBenchmarkCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindInstance mocks base method.
func (m *MockRepository) FindInstance(ctx context.Context, id uuid.UUID) (domain.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindInstance", ctx, id)
	ret0, _ := ret[0].(domain.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindInstance indicates an expected call of FindInstance.
func (mr *MockRepositoryMockRecorder) FindInstance(ctx, id any) *MockRepositoryFindInstanceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindInstance", reflect.TypeOf((*MockRepository)(nil).FindInstance), ctx, id)
	return &MockRepositoryFindInstanceCall{Call: call}
}

// MockRepositoryFindInstanceCall wrap *gomock.Call
type MockRepositoryFindInstanceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindInstanceCall) Return(arg0 domain.Instance, arg1 error) *MockRepositoryFindInstanceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindInstanceCall) Do(f func(context.Context, uuid.UUID) (domain.Instance, error)) *MockRepositoryFindInstanceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindInstanceCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.Instance, error)) *MockRepositoryFindInstanceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindSession mocks base method.
func (m *MockRepository) FindSession(ctx context.Context, id string) (domain.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSession", ctx, id)
	ret0, _ := ret[0].(domain.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSession indicates an expected call of FindSession.
func (mr *MockRepositoryMockRecorder) FindSession(ctx, id any) *MockRepositoryFindSessionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSession", reflect.TypeOf((*MockRepository)(nil).FindSession), ctx, id)
	return &MockRepositoryFindSessionCall{Call: call}
}

// MockRepositoryFindSessionCall wrap *gomock.Call
type MockRepositoryFindSessionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindSessionCall) Return(arg0 domain.Session, arg1 error) *MockRepositoryFindSessionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindSessionCall) Do(f func(context.Context, string) (domain.Session, error)) *MockRepositoryFindSessionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindSessionCall) DoAndReturn(f func(context.Context, string) (domain.Session, error)) *MockRepositoryFindSessionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindTeam mocks base method.
func (m *MockRepository) FindTeam(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTeam", ctx, id)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTeam indicates an expected call of FindTeam.
func (mr *MockRepositoryMockRecorder) FindTeam(ctx, id any) *MockRepositoryFindTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTeam", reflect.TypeOf((*MockRepository)(nil).FindTeam), ctx, id)
	return &MockRepositoryFindTeamCall{Call: call}
}

// MockRepositoryFindTeamCall wrap *gomock.Call
type MockRepositoryFindTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindTeamCall) Return(arg0 domain.Team, arg1 error) *MockRepositoryFindTeamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindTeamCall) Do(f func(context.Context, uuid.UUID) (domain.Team, error)) *MockRepositoryFindTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindTeamCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.Team, error)) *MockRepositoryFindTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindUser mocks base method.
func (m *MockRepository) FindUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockRepositoryMockRecorder) FindUser(ctx, id any) *MockRepositoryFindUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockRepository)(nil).FindUser), ctx, id)
	return &MockRepositoryFindUserCall{Call: call}
}

// MockRepositoryFindUserCall wrap *gomock.Call
type MockRepositoryFindUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindUserCall) Return(arg0 domain.User, arg1 error) *MockRepositoryFindUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindUserCall) Do(f func(context.Context, uuid.UUID) (domain.User, error)) *MockRepositoryFindUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindUserCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.User, error)) *MockRepositoryFindUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBenchmarkLog mocks base method.
func (m *MockRepository) GetBenchmarkLog(ctx context.Context, benchmarkID uuid.UUID) (domain.BenchmarkLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBenchmarkLog", ctx, benchmarkID)
	ret0, _ := ret[0].(domain.BenchmarkLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBenchmarkLog indicates an expected call of GetBenchmarkLog.
func (mr *MockRepositoryMockRecorder) GetBenchmarkLog(ctx, benchmarkID any) *MockRepositoryGetBenchmarkLogCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBenchmarkLog", reflect.TypeOf((*MockRepository)(nil).GetBenchmarkLog), ctx, benchmarkID)
	return &MockRepositoryGetBenchmarkLogCall{Call: call}
}

// MockRepositoryGetBenchmarkLogCall wrap *gomock.Call
type MockRepositoryGetBenchmarkLogCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryGetBenchmarkLogCall) Return(arg0 domain.BenchmarkLog, arg1 error) *MockRepositoryGetBenchmarkLogCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryGetBenchmarkLogCall) Do(f func(context.Context, uuid.UUID) (domain.BenchmarkLog, error)) *MockRepositoryGetBenchmarkLogCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryGetBenchmarkLogCall) DoAndReturn(f func(context.Context, uuid.UUID) (domain.BenchmarkLog, error)) *MockRepositoryGetBenchmarkLogCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBenchmarks mocks base method.
func (m *MockRepository) GetBenchmarks(ctx context.Context, query repository.BenchmarkQuery) ([]domain.Benchmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBenchmarks", ctx, query)
	ret0, _ := ret[0].([]domain.Benchmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBenchmarks indicates an expected call of GetBenchmarks.
func (mr *MockRepositoryMockRecorder) GetBenchmarks(ctx, query any) *MockRepositoryGetBenchmarksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBenchmarks", reflect.TypeOf((*MockRepository)(nil).GetBenchmarks), ctx, query)
	return &MockRepositoryGetBenchmarksCall{Call: call}
}

// MockRepositoryGetBenchmarksCall wrap *gomock.Call
type MockRepositoryGetBenchmarksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryGetBenchmarksCall) Return(arg0 []domain.Benchmark, arg1 error) *MockRepositoryGetBenchmarksCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryGetBenchmarksCall) Do(f func(context.Context, repository.BenchmarkQuery) ([]domain.Benchmark, error)) *MockRepositoryGetBenchmarksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryGetBenchmarksCall) DoAndReturn(f func(context.Context, repository.BenchmarkQuery) ([]domain.Benchmark, error)) *MockRepositoryGetBenchmarksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetTeams mocks base method.
func (m *MockRepository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeams", ctx)
	ret0, _ := ret[0].([]domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeams indicates an expected call of GetTeams.
func (mr *MockRepositoryMockRecorder) GetTeams(ctx any) *MockRepositoryGetTeamsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeams", reflect.TypeOf((*MockRepository)(nil).GetTeams), ctx)
	return &MockRepositoryGetTeamsCall{Call: call}
}

// MockRepositoryGetTeamsCall wrap *gomock.Call
type MockRepositoryGetTeamsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryGetTeamsCall) Return(arg0 []domain.Team, arg1 error) *MockRepositoryGetTeamsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryGetTeamsCall) Do(f func(context.Context) ([]domain.Team, error)) *MockRepositoryGetTeamsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryGetTeamsCall) DoAndReturn(f func(context.Context) ([]domain.Team, error)) *MockRepositoryGetTeamsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUsers mocks base method.
func (m *MockRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockRepositoryMockRecorder) GetUsers(ctx any) *MockRepositoryGetUsersCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockRepository)(nil).GetUsers), ctx)
	return &MockRepositoryGetUsersCall{Call: call}
}

// MockRepositoryGetUsersCall wrap *gomock.Call
type MockRepositoryGetUsersCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryGetUsersCall) Return(arg0 []domain.User, arg1 error) *MockRepositoryGetUsersCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryGetUsersCall) Do(f func(context.Context) ([]domain.User, error)) *MockRepositoryGetUsersCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryGetUsersCall) DoAndReturn(f func(context.Context) ([]domain.User, error)) *MockRepositoryGetUsersCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Transaction mocks base method.
func (m *MockRepository) Transaction(ctx context.Context, f func(context.Context, repository.Repository) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockRepositoryMockRecorder) Transaction(ctx, f any) *MockRepositoryTransactionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockRepository)(nil).Transaction), ctx, f)
	return &MockRepositoryTransactionCall{Call: call}
}

// MockRepositoryTransactionCall wrap *gomock.Call
type MockRepositoryTransactionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryTransactionCall) Return(arg0 error) *MockRepositoryTransactionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryTransactionCall) Do(f func(context.Context, func(context.Context, repository.Repository) error) error) *MockRepositoryTransactionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryTransactionCall) DoAndReturn(f func(context.Context, func(context.Context, repository.Repository) error) error) *MockRepositoryTransactionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateTeam mocks base method.
func (m *MockRepository) UpdateTeam(ctx context.Context, team domain.Team) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTeam", ctx, team)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTeam indicates an expected call of UpdateTeam.
func (mr *MockRepositoryMockRecorder) UpdateTeam(ctx, team any) *MockRepositoryUpdateTeamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTeam", reflect.TypeOf((*MockRepository)(nil).UpdateTeam), ctx, team)
	return &MockRepositoryUpdateTeamCall{Call: call}
}

// MockRepositoryUpdateTeamCall wrap *gomock.Call
type MockRepositoryUpdateTeamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryUpdateTeamCall) Return(arg0 error) *MockRepositoryUpdateTeamCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryUpdateTeamCall) Do(f func(context.Context, domain.Team) error) *MockRepositoryUpdateTeamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryUpdateTeamCall) DoAndReturn(f func(context.Context, domain.Team) error) *MockRepositoryUpdateTeamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
