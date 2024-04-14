// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/benderr/keypass/internal/server/domain/record/usecase (interfaces: RecordRepo,DataCrypter)
//
// Generated by this command:
//
//	mockgen -destination=internal/server/domain/record/usecase/mocks/mocks.go -package=mocks github.com/benderr/keypass/internal/server/domain/record/usecase RecordRepo,DataCrypter
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	record "github.com/benderr/keypass/internal/server/domain/record"
	gomock "go.uber.org/mock/gomock"
)

// MockRecordRepo is a mock of RecordRepo interface.
type MockRecordRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRecordRepoMockRecorder
}

// MockRecordRepoMockRecorder is the mock recorder for MockRecordRepo.
type MockRecordRepoMockRecorder struct {
	mock *MockRecordRepo
}

// NewMockRecordRepo creates a new mock instance.
func NewMockRecordRepo(ctrl *gomock.Controller) *MockRecordRepo {
	mock := &MockRecordRepo{ctrl: ctrl}
	mock.recorder = &MockRecordRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecordRepo) EXPECT() *MockRecordRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRecordRepo) Create(arg0 context.Context, arg1 string, arg2 []byte, arg3, arg4 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRecordRepoMockRecorder) Create(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRecordRepo)(nil).Create), arg0, arg1, arg2, arg3, arg4)
}

// Delete mocks base method.
func (m *MockRecordRepo) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRecordRepoMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRecordRepo)(nil).Delete), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockRecordRepo) GetByID(arg0 context.Context, arg1 string) (*record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRecordRepoMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRecordRepo)(nil).GetByID), arg0, arg1)
}

// GetByUser mocks base method.
func (m *MockRecordRepo) GetByUser(arg0 context.Context, arg1 string) ([]record.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUser", arg0, arg1)
	ret0, _ := ret[0].([]record.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUser indicates an expected call of GetByUser.
func (mr *MockRecordRepoMockRecorder) GetByUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUser", reflect.TypeOf((*MockRecordRepo)(nil).GetByUser), arg0, arg1)
}

// Update mocks base method.
func (m *MockRecordRepo) Update(arg0 context.Context, arg1 string, arg2 []byte, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRecordRepoMockRecorder) Update(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRecordRepo)(nil).Update), arg0, arg1, arg2, arg3)
}

// MockDataCrypter is a mock of DataCrypter interface.
type MockDataCrypter struct {
	ctrl     *gomock.Controller
	recorder *MockDataCrypterMockRecorder
}

// MockDataCrypterMockRecorder is the mock recorder for MockDataCrypter.
type MockDataCrypterMockRecorder struct {
	mock *MockDataCrypter
}

// NewMockDataCrypter creates a new mock instance.
func NewMockDataCrypter(ctrl *gomock.Controller) *MockDataCrypter {
	mock := &MockDataCrypter{ctrl: ctrl}
	mock.recorder = &MockDataCrypterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataCrypter) EXPECT() *MockDataCrypterMockRecorder {
	return m.recorder
}

// Decrypt mocks base method.
func (m *MockDataCrypter) Decrypt(arg0 []byte) (map[string]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", arg0)
	ret0, _ := ret[0].(map[string]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt.
func (mr *MockDataCrypterMockRecorder) Decrypt(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockDataCrypter)(nil).Decrypt), arg0)
}

// Encrypt mocks base method.
func (m *MockDataCrypter) Encrypt(arg0 any) ([]byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Encrypt indicates an expected call of Encrypt.
func (mr *MockDataCrypterMockRecorder) Encrypt(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockDataCrypter)(nil).Encrypt), arg0)
}
