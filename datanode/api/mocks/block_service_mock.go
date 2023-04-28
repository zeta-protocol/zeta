// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/datanode/api (interfaces: BlockService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/zeta-protocol/zeta/datanode/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockBlockService is a mock of BlockService interface.
type MockBlockService struct {
	ctrl     *gomock.Controller
	recorder *MockBlockServiceMockRecorder
}

// MockBlockServiceMockRecorder is the mock recorder for MockBlockService.
type MockBlockServiceMockRecorder struct {
	mock *MockBlockService
}

// NewMockBlockService creates a new mock instance.
func NewMockBlockService(ctrl *gomock.Controller) *MockBlockService {
	mock := &MockBlockService{ctrl: ctrl}
	mock.recorder = &MockBlockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockService) EXPECT() *MockBlockServiceMockRecorder {
	return m.recorder
}

// GetLastBlock mocks base method.
func (m *MockBlockService) GetLastBlock(arg0 context.Context) (entities.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastBlock", arg0)
	ret0, _ := ret[0].(entities.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastBlock indicates an expected call of GetLastBlock.
func (mr *MockBlockServiceMockRecorder) GetLastBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastBlock", reflect.TypeOf((*MockBlockService)(nil).GetLastBlock), arg0)
}
