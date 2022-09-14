// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/transport/http/server.go

// Package mock_http is a generated GoMock package.
package mock_http

import (
	reflect "reflect"

	router "github.com/fasthttp/router"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Init mocks base method.
func (m *MockHandler) Init(r *router.Router) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init", r)
}

// Init indicates an expected call of Init.
func (mr *MockHandlerMockRecorder) Init(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockHandler)(nil).Init), r)
}
