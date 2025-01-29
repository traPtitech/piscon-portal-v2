// Code generated by MockGen. DO NOT EDIT.
// Source: execute.go
//
// Generated by this command:
//
//	mockgen -source=execute.go -destination=mock/execute.go -package=mock -typed=true
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	domain "github.com/traPtitech/piscon-portal-v2/runner/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockBenchmarker is a mock of Benchmarker interface.
type MockBenchmarker struct {
	ctrl     *gomock.Controller
	recorder *MockBenchmarkerMockRecorder
	isgomock struct{}
}

// MockBenchmarkerMockRecorder is the mock recorder for MockBenchmarker.
type MockBenchmarkerMockRecorder struct {
	mock *MockBenchmarker
}

// NewMockBenchmarker creates a new mock instance.
func NewMockBenchmarker(ctrl *gomock.Controller) *MockBenchmarker {
	mock := &MockBenchmarker{ctrl: ctrl}
	mock.recorder = &MockBenchmarkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBenchmarker) EXPECT() *MockBenchmarkerMockRecorder {
	return m.recorder
}

// CalculateScore mocks base method.
func (m *MockBenchmarker) CalculateScore(ctx context.Context, allStdout, allStderr string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateScore", ctx, allStdout, allStderr)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateScore indicates an expected call of CalculateScore.
func (mr *MockBenchmarkerMockRecorder) CalculateScore(ctx, allStdout, allStderr any) *MockBenchmarkerCalculateScoreCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateScore", reflect.TypeOf((*MockBenchmarker)(nil).CalculateScore), ctx, allStdout, allStderr)
	return &MockBenchmarkerCalculateScoreCall{Call: call}
}

// MockBenchmarkerCalculateScoreCall wrap *gomock.Call
type MockBenchmarkerCalculateScoreCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBenchmarkerCalculateScoreCall) Return(arg0 int, arg1 error) *MockBenchmarkerCalculateScoreCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBenchmarkerCalculateScoreCall) Do(f func(context.Context, string, string) (int, error)) *MockBenchmarkerCalculateScoreCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBenchmarkerCalculateScoreCall) DoAndReturn(f func(context.Context, string, string) (int, error)) *MockBenchmarkerCalculateScoreCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Start mocks base method.
func (m *MockBenchmarker) Start(ctx context.Context, job *domain.Job) (io.ReadCloser, io.ReadCloser, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx, job)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(time.Time)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Start indicates an expected call of Start.
func (mr *MockBenchmarkerMockRecorder) Start(ctx, job any) *MockBenchmarkerStartCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBenchmarker)(nil).Start), ctx, job)
	return &MockBenchmarkerStartCall{Call: call}
}

// MockBenchmarkerStartCall wrap *gomock.Call
type MockBenchmarkerStartCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBenchmarkerStartCall) Return(stdout, stderr io.ReadCloser, startedAt time.Time, err error) *MockBenchmarkerStartCall {
	c.Call = c.Call.Return(stdout, stderr, startedAt, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBenchmarkerStartCall) Do(f func(context.Context, *domain.Job) (io.ReadCloser, io.ReadCloser, time.Time, error)) *MockBenchmarkerStartCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBenchmarkerStartCall) DoAndReturn(f func(context.Context, *domain.Job) (io.ReadCloser, io.ReadCloser, time.Time, error)) *MockBenchmarkerStartCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Wait mocks base method.
func (m *MockBenchmarker) Wait(ctx context.Context) (domain.Result, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait", ctx)
	ret0, _ := ret[0].(domain.Result)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Wait indicates an expected call of Wait.
func (mr *MockBenchmarkerMockRecorder) Wait(ctx any) *MockBenchmarkerWaitCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockBenchmarker)(nil).Wait), ctx)
	return &MockBenchmarkerWaitCall{Call: call}
}

// MockBenchmarkerWaitCall wrap *gomock.Call
type MockBenchmarkerWaitCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBenchmarkerWaitCall) Return(arg0 domain.Result, arg1 time.Time, arg2 error) *MockBenchmarkerWaitCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBenchmarkerWaitCall) Do(f func(context.Context) (domain.Result, time.Time, error)) *MockBenchmarkerWaitCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBenchmarkerWaitCall) DoAndReturn(f func(context.Context) (domain.Result, time.Time, error)) *MockBenchmarkerWaitCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
