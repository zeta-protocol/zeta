// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/core/types (interfaces: PostRestore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/zeta-protocol/zeta/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockPostRestore is a mock of PostRestore interface.
type MockPostRestore struct {
	ctrl     *gomock.Controller
	recorder *MockPostRestoreMockRecorder
}

// MockPostRestoreMockRecorder is the mock recorder for MockPostRestore.
type MockPostRestoreMockRecorder struct {
	mock *MockPostRestore
}

// NewMockPostRestore creates a new mock instance.
func NewMockPostRestore(ctrl *gomock.Controller) *MockPostRestore {
	mock := &MockPostRestore{ctrl: ctrl}
	mock.recorder = &MockPostRestoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRestore) EXPECT() *MockPostRestoreMockRecorder {
	return m.recorder
}

// GetState mocks base method.
func (m *MockPostRestore) GetState(arg0 string) ([]byte, []types.StateProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]types.StateProvider)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetState indicates an expected call of GetState.
func (mr *MockPostRestoreMockRecorder) GetState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockPostRestore)(nil).GetState), arg0)
}

// Keys mocks base method.
func (m *MockPostRestore) Keys() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Keys indicates an expected call of Keys.
func (mr *MockPostRestoreMockRecorder) Keys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockPostRestore)(nil).Keys))
}

// LoadState mocks base method.
func (m *MockPostRestore) LoadState(arg0 context.Context, arg1 *types.Payload) ([]types.StateProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadState", arg0, arg1)
	ret0, _ := ret[0].([]types.StateProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadState indicates an expected call of LoadState.
func (mr *MockPostRestoreMockRecorder) LoadState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadState", reflect.TypeOf((*MockPostRestore)(nil).LoadState), arg0, arg1)
}

// Namespace mocks base method.
func (m *MockPostRestore) Namespace() types.SnapshotNamespace {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Namespace")
	ret0, _ := ret[0].(types.SnapshotNamespace)
	return ret0
}

// Namespace indicates an expected call of Namespace.
func (mr *MockPostRestoreMockRecorder) Namespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespace", reflect.TypeOf((*MockPostRestore)(nil).Namespace))
}

// OnStateLoaded mocks base method.
func (m *MockPostRestore) OnStateLoaded(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnStateLoaded", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnStateLoaded indicates an expected call of OnStateLoaded.
func (mr *MockPostRestoreMockRecorder) OnStateLoaded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnStateLoaded", reflect.TypeOf((*MockPostRestore)(nil).OnStateLoaded), arg0)
}

// Stopped mocks base method.
func (m *MockPostRestore) Stopped() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stopped")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Stopped indicates an expected call of Stopped.
func (mr *MockPostRestoreMockRecorder) Stopped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stopped", reflect.TypeOf((*MockPostRestore)(nil).Stopped))
}
