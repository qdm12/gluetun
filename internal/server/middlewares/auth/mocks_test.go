// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/gluetun/internal/server/middlewares/auth (interfaces: DebugLogger)

// Package auth is a generated GoMock package.
package auth

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDebugLogger is a mock of DebugLogger interface.
type MockDebugLogger struct {
	ctrl     *gomock.Controller
	recorder *MockDebugLoggerMockRecorder
}

// MockDebugLoggerMockRecorder is the mock recorder for MockDebugLogger.
type MockDebugLoggerMockRecorder struct {
	mock *MockDebugLogger
}

// NewMockDebugLogger creates a new mock instance.
func NewMockDebugLogger(ctrl *gomock.Controller) *MockDebugLogger {
	mock := &MockDebugLogger{ctrl: ctrl}
	mock.recorder = &MockDebugLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDebugLogger) EXPECT() *MockDebugLoggerMockRecorder {
	return m.recorder
}

// Debugf mocks base method.
func (m *MockDebugLogger) Debugf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockDebugLoggerMockRecorder) Debugf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockDebugLogger)(nil).Debugf), varargs...)
}

// Warnf mocks base method.
func (m *MockDebugLogger) Warnf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockDebugLoggerMockRecorder) Warnf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*MockDebugLogger)(nil).Warnf), varargs...)
}
