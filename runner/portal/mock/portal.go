// Code generated by MockGen. DO NOT EDIT.
// Source: portal.go
//
// Generated by this command:
//
//	mockgen -source=portal.go -destination=mock/portal.go -package=mock -typed=true
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	domain "github.com/traPtitech/piscon-portal-v2/runner/domain"
	portal "github.com/traPtitech/piscon-portal-v2/runner/portal"
	gomock "go.uber.org/mock/gomock"
)

// MockPortal is a mock of Portal interface.
type MockPortal struct {
	ctrl     *gomock.Controller
	recorder *MockPortalMockRecorder
	isgomock struct{}
}

// MockPortalMockRecorder is the mock recorder for MockPortal.
type MockPortalMockRecorder struct {
	mock *MockPortal
}

// NewMockPortal creates a new mock instance.
func NewMockPortal(ctrl *gomock.Controller) *MockPortal {
	mock := &MockPortal{ctrl: ctrl}
	mock.recorder = &MockPortalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPortal) EXPECT() *MockPortalMockRecorder {
	return m.recorder
}

// GetJob mocks base method.
func (m *MockPortal) GetJob(ctx context.Context) (*domain.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJob", ctx)
	ret0, _ := ret[0].(*domain.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJob indicates an expected call of GetJob.
func (mr *MockPortalMockRecorder) GetJob(ctx any) *MockPortalGetJobCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJob", reflect.TypeOf((*MockPortal)(nil).GetJob), ctx)
	return &MockPortalGetJobCall{Call: call}
}

// MockPortalGetJobCall wrap *gomock.Call
type MockPortalGetJobCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPortalGetJobCall) Return(arg0 *domain.Job, arg1 error) *MockPortalGetJobCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPortalGetJobCall) Do(f func(context.Context) (*domain.Job, error)) *MockPortalGetJobCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPortalGetJobCall) DoAndReturn(f func(context.Context) (*domain.Job, error)) *MockPortalGetJobCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MakeProgressStreamClient mocks base method.
func (m *MockPortal) MakeProgressStreamClient(ctx context.Context) (portal.ProgressStreamClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeProgressStreamClient", ctx)
	ret0, _ := ret[0].(portal.ProgressStreamClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeProgressStreamClient indicates an expected call of MakeProgressStreamClient.
func (mr *MockPortalMockRecorder) MakeProgressStreamClient(ctx any) *MockPortalMakeProgressStreamClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeProgressStreamClient", reflect.TypeOf((*MockPortal)(nil).MakeProgressStreamClient), ctx)
	return &MockPortalMakeProgressStreamClientCall{Call: call}
}

// MockPortalMakeProgressStreamClientCall wrap *gomock.Call
type MockPortalMakeProgressStreamClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPortalMakeProgressStreamClientCall) Return(arg0 portal.ProgressStreamClient, arg1 error) *MockPortalMakeProgressStreamClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPortalMakeProgressStreamClientCall) Do(f func(context.Context) (portal.ProgressStreamClient, error)) *MockPortalMakeProgressStreamClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPortalMakeProgressStreamClientCall) DoAndReturn(f func(context.Context) (portal.ProgressStreamClient, error)) *MockPortalMakeProgressStreamClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PostJobFinished mocks base method.
func (m *MockPortal) PostJobFinished(ctx context.Context, jobID string, finishedAt time.Time, result domain.Result, runnerErr error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostJobFinished", ctx, jobID, finishedAt, result, runnerErr)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostJobFinished indicates an expected call of PostJobFinished.
func (mr *MockPortalMockRecorder) PostJobFinished(ctx, jobID, finishedAt, result, runnerErr any) *MockPortalPostJobFinishedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostJobFinished", reflect.TypeOf((*MockPortal)(nil).PostJobFinished), ctx, jobID, finishedAt, result, runnerErr)
	return &MockPortalPostJobFinishedCall{Call: call}
}

// MockPortalPostJobFinishedCall wrap *gomock.Call
type MockPortalPostJobFinishedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPortalPostJobFinishedCall) Return(arg0 error) *MockPortalPostJobFinishedCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPortalPostJobFinishedCall) Do(f func(context.Context, string, time.Time, domain.Result, error) error) *MockPortalPostJobFinishedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPortalPostJobFinishedCall) DoAndReturn(f func(context.Context, string, time.Time, domain.Result, error) error) *MockPortalPostJobFinishedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockProgressStreamClient is a mock of ProgressStreamClient interface.
type MockProgressStreamClient struct {
	ctrl     *gomock.Controller
	recorder *MockProgressStreamClientMockRecorder
	isgomock struct{}
}

// MockProgressStreamClientMockRecorder is the mock recorder for MockProgressStreamClient.
type MockProgressStreamClientMockRecorder struct {
	mock *MockProgressStreamClient
}

// NewMockProgressStreamClient creates a new mock instance.
func NewMockProgressStreamClient(ctrl *gomock.Controller) *MockProgressStreamClient {
	mock := &MockProgressStreamClient{ctrl: ctrl}
	mock.recorder = &MockProgressStreamClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProgressStreamClient) EXPECT() *MockProgressStreamClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockProgressStreamClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockProgressStreamClientMockRecorder) Close() *MockProgressStreamClientCloseCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockProgressStreamClient)(nil).Close))
	return &MockProgressStreamClientCloseCall{Call: call}
}

// MockProgressStreamClientCloseCall wrap *gomock.Call
type MockProgressStreamClientCloseCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProgressStreamClientCloseCall) Return(arg0 error) *MockProgressStreamClientCloseCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProgressStreamClientCloseCall) Do(f func() error) *MockProgressStreamClientCloseCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProgressStreamClientCloseCall) DoAndReturn(f func() error) *MockProgressStreamClientCloseCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendProgress mocks base method.
func (m *MockProgressStreamClient) SendProgress(ctx context.Context, progress *domain.Progress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendProgress", ctx, progress)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendProgress indicates an expected call of SendProgress.
func (mr *MockProgressStreamClientMockRecorder) SendProgress(ctx, progress any) *MockProgressStreamClientSendProgressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendProgress", reflect.TypeOf((*MockProgressStreamClient)(nil).SendProgress), ctx, progress)
	return &MockProgressStreamClientSendProgressCall{Call: call}
}

// MockProgressStreamClientSendProgressCall wrap *gomock.Call
type MockProgressStreamClientSendProgressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProgressStreamClientSendProgressCall) Return(arg0 error) *MockProgressStreamClientSendProgressCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProgressStreamClientSendProgressCall) Do(f func(context.Context, *domain.Progress) error) *MockProgressStreamClientSendProgressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProgressStreamClientSendProgressCall) DoAndReturn(f func(context.Context, *domain.Progress) error) *MockProgressStreamClientSendProgressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
