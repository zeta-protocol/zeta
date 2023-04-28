// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/wallet/service (interfaces: NetworkStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	network "github.com/zeta-protocol/zeta/wallet/network"
	gomock "github.com/golang/mock/gomock"
)

// MockNetworkStore is a mock of NetworkStore interface.
type MockNetworkStore struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkStoreMockRecorder
}

// MockNetworkStoreMockRecorder is the mock recorder for MockNetworkStore.
type MockNetworkStoreMockRecorder struct {
	mock *MockNetworkStore
}

// NewMockNetworkStore creates a new mock instance.
func NewMockNetworkStore(ctrl *gomock.Controller) *MockNetworkStore {
	mock := &MockNetworkStore{ctrl: ctrl}
	mock.recorder = &MockNetworkStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkStore) EXPECT() *MockNetworkStoreMockRecorder {
	return m.recorder
}

// GetNetwork mocks base method.
func (m *MockNetworkStore) GetNetwork(arg0 string) (*network.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetwork", arg0)
	ret0, _ := ret[0].(*network.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetwork indicates an expected call of GetNetwork.
func (mr *MockNetworkStoreMockRecorder) GetNetwork(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetwork", reflect.TypeOf((*MockNetworkStore)(nil).GetNetwork), arg0)
}

// NetworkExists mocks base method.
func (m *MockNetworkStore) NetworkExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkExists indicates an expected call of NetworkExists.
func (mr *MockNetworkStoreMockRecorder) NetworkExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkExists", reflect.TypeOf((*MockNetworkStore)(nil).NetworkExists), arg0)
}
