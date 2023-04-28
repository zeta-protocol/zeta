// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/core/checkpoint (interfaces: State)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/zeta-protocol/zeta/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// Checkpoint mocks base method.
func (m *MockState) Checkpoint() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkpoint")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Checkpoint indicates an expected call of Checkpoint.
func (mr *MockStateMockRecorder) Checkpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkpoint", reflect.TypeOf((*MockState)(nil).Checkpoint))
}

// Load mocks base method.
func (m *MockState) Load(arg0 context.Context, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockStateMockRecorder) Load(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockState)(nil).Load), arg0, arg1)
}

// Name mocks base method.
func (m *MockState) Name() types.CheckpointName {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(types.CheckpointName)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockStateMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockState)(nil).Name))
}