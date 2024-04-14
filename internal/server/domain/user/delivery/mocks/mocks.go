// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/benderr/keypass/internal/domain/user/delivery (interfaces: UserUsecase,SessionManager)
//
// Generated by this command:
//
//	mockgen -destination=internal/domain/user/delivery/mocks/mocks.go -package=mocks github.com/benderr/keypass/internal/domain/user/delivery UserUsecase,SessionManager
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	user "github.com/benderr/keypass/internal/server/domain/user"
	gomock "go.uber.org/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockUserUsecase) Login(arg0 context.Context, arg1, arg2 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1, arg2)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserUsecaseMockRecorder) Login(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUsecase)(nil).Login), arg0, arg1, arg2)
}

// Register mocks base method.
func (m *MockUserUsecase) Register(arg0 context.Context, arg1, arg2 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1, arg2)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserUsecaseMockRecorder) Register(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserUsecase)(nil).Register), arg0, arg1, arg2)
}

// MockSessionManager is a mock of SessionManager interface.
type MockSessionManager struct {
	ctrl     *gomock.Controller
	recorder *MockSessionManagerMockRecorder
}

// MockSessionManagerMockRecorder is the mock recorder for MockSessionManager.
type MockSessionManagerMockRecorder struct {
	mock *MockSessionManager
}

// NewMockSessionManager creates a new mock instance.
func NewMockSessionManager(ctrl *gomock.Controller) *MockSessionManager {
	mock := &MockSessionManager{ctrl: ctrl}
	mock.recorder = &MockSessionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionManager) EXPECT() *MockSessionManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionManager) Create(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionManagerMockRecorder) Create(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionManager)(nil).Create), arg0)
}
