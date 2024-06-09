// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sender (interfaces: Mailer)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_mailer.go -package=mocks . Mailer
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMailer is a mock of Mailer interface.
type MockMailer struct {
	ctrl     *gomock.Controller
	recorder *MockMailerMockRecorder
}

// MockMailerMockRecorder is the mock recorder for MockMailer.
type MockMailerMockRecorder struct {
	mock *MockMailer
}

// NewMockMailer creates a new mock instance.
func NewMockMailer(ctrl *gomock.Controller) *MockMailer {
	mock := &MockMailer{ctrl: ctrl}
	mock.recorder = &MockMailerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailer) EXPECT() *MockMailerMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockMailer) Send(arg0 context.Context, arg1, arg2 string, arg3 ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockMailerMockRecorder) Send(arg0, arg1, arg2 any, arg3 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockMailer)(nil).Send), varargs...)
}
