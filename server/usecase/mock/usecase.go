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
