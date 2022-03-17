// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/transport/http/handler/handler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRequestParser is a mock of RequestParser interface.
type MockRequestParser struct {
	ctrl     *gomock.Controller
	recorder *MockRequestParserMockRecorder
}

// MockRequestParserMockRecorder is the mock recorder for MockRequestParser.
type MockRequestParserMockRecorder struct {
	mock *MockRequestParser
}

// NewMockRequestParser creates a new mock instance.
func NewMockRequestParser(ctrl *gomock.Controller) *MockRequestParser {
	mock := &MockRequestParser{ctrl: ctrl}
	mock.recorder = &MockRequestParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequestParser) EXPECT() *MockRequestParserMockRecorder {
	return m.recorder
}

// ParseInt64FromBytes mocks base method.
func (m *MockRequestParser) ParseInt64FromBytes(arg0 []byte) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseInt64FromBytes", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseInt64FromBytes indicates an expected call of ParseInt64FromBytes.
func (mr *MockRequestParserMockRecorder) ParseInt64FromBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseInt64FromBytes", reflect.TypeOf((*MockRequestParser)(nil).ParseInt64FromBytes), arg0)
}

// ParseInt64FromInterface mocks base method.
func (m *MockRequestParser) ParseInt64FromInterface(arg0 interface{}) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseInt64FromInterface", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseInt64FromInterface indicates an expected call of ParseInt64FromInterface.
func (mr *MockRequestParserMockRecorder) ParseInt64FromInterface(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseInt64FromInterface", reflect.TypeOf((*MockRequestParser)(nil).ParseInt64FromInterface), arg0)
}