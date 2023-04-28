// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/core/staking (interfaces: EthereumClientConfirmations)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	big "math/big"
	reflect "reflect"

	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockEthereumClientConfirmations is a mock of EthereumClientConfirmations interface.
type MockEthereumClientConfirmations struct {
	ctrl     *gomock.Controller
	recorder *MockEthereumClientConfirmationsMockRecorder
}

// MockEthereumClientConfirmationsMockRecorder is the mock recorder for MockEthereumClientConfirmations.
type MockEthereumClientConfirmationsMockRecorder struct {
	mock *MockEthereumClientConfirmations
}

// NewMockEthereumClientConfirmations creates a new mock instance.
func NewMockEthereumClientConfirmations(ctrl *gomock.Controller) *MockEthereumClientConfirmations {
	mock := &MockEthereumClientConfirmations{ctrl: ctrl}
	mock.recorder = &MockEthereumClientConfirmationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEthereumClientConfirmations) EXPECT() *MockEthereumClientConfirmationsMockRecorder {
	return m.recorder
}

// HeaderByNumber mocks base method.
func (m *MockEthereumClientConfirmations) HeaderByNumber(arg0 context.Context, arg1 *big.Int) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumber", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumber indicates an expected call of HeaderByNumber.
func (mr *MockEthereumClientConfirmationsMockRecorder) HeaderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumber", reflect.TypeOf((*MockEthereumClientConfirmations)(nil).HeaderByNumber), arg0, arg1)
}