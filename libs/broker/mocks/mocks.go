// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zeta-protocol/zeta/libs/broker (interfaces: Subscriber)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	events "github.com/zeta-protocol/zeta/core/events"
	gomock "github.com/golang/mock/gomock"
)

// MockSubscriber is a mock of Subscriber interface.
type MockSubscriber struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriberMockRecorder
}

// MockSubscriberMockRecorder is the mock recorder for MockSubscriber.
type MockSubscriberMockRecorder struct {
	mock *MockSubscriber
}

// NewMockSubscriber creates a new mock instance.
func NewMockSubscriber(ctrl *gomock.Controller) *MockSubscriber {
	mock := &MockSubscriber{ctrl: ctrl}
	mock.recorder = &MockSubscriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriber) EXPECT() *MockSubscriberMockRecorder {
	return m.recorder
}

// Ack mocks base method.
func (m *MockSubscriber) Ack() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ack")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Ack indicates an expected call of Ack.
func (mr *MockSubscriberMockRecorder) Ack() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockSubscriber)(nil).Ack))
}

// C mocks base method.
func (m *MockSubscriber) C() chan<- []events.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "C")
	ret0, _ := ret[0].(chan<- []events.Event)
	return ret0
}

// C indicates an expected call of C.
func (mr *MockSubscriberMockRecorder) C() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "C", reflect.TypeOf((*MockSubscriber)(nil).C))
}

// Closed mocks base method.
func (m *MockSubscriber) Closed() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Closed")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Closed indicates an expected call of Closed.
func (mr *MockSubscriberMockRecorder) Closed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Closed", reflect.TypeOf((*MockSubscriber)(nil).Closed))
}

// ID mocks base method.
func (m *MockSubscriber) ID() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(int)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockSubscriberMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockSubscriber)(nil).ID))
}

// Push mocks base method.
func (m *MockSubscriber) Push(arg0 ...events.Event) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Push", varargs...)
}

// Push indicates an expected call of Push.
func (mr *MockSubscriberMockRecorder) Push(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockSubscriber)(nil).Push), arg0...)
}

// SetID mocks base method.
func (m *MockSubscriber) SetID(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetID", arg0)
}

// SetID indicates an expected call of SetID.
func (mr *MockSubscriberMockRecorder) SetID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetID", reflect.TypeOf((*MockSubscriber)(nil).SetID), arg0)
}

// Skip mocks base method.
func (m *MockSubscriber) Skip() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Skip")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Skip indicates an expected call of Skip.
func (mr *MockSubscriberMockRecorder) Skip() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skip", reflect.TypeOf((*MockSubscriber)(nil).Skip))
}

// Types mocks base method.
func (m *MockSubscriber) Types() []events.Type {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Types")
	ret0, _ := ret[0].([]events.Type)
	return ret0
}

// Types indicates an expected call of Types.
func (mr *MockSubscriberMockRecorder) Types() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Types", reflect.TypeOf((*MockSubscriber)(nil).Types))
}
