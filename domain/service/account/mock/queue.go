// Code generated by MockGen. DO NOT EDIT.
// Source: domain/pkg/queue/queue.go
//
// Generated by this command:
//
//	mockgen -source domain/pkg/queue/queue.go -destination domain/service/account/mock/queue.go -package mock_account
//

// Package mock_account is a generated GoMock package.
package mock_account

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEventQueue is a mock of EventQueue interface.
type MockEventQueue struct {
	ctrl     *gomock.Controller
	recorder *MockEventQueueMockRecorder
}

// MockEventQueueMockRecorder is the mock recorder for MockEventQueue.
type MockEventQueueMockRecorder struct {
	mock *MockEventQueue
}

// NewMockEventQueue creates a new mock instance.
func NewMockEventQueue(ctrl *gomock.Controller) *MockEventQueue {
	mock := &MockEventQueue{ctrl: ctrl}
	mock.recorder = &MockEventQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventQueue) EXPECT() *MockEventQueueMockRecorder {
	return m.recorder
}

// EnQueue mocks base method.
func (m *MockEventQueue) EnQueue(event any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnQueue", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnQueue indicates an expected call of EnQueue.
func (mr *MockEventQueueMockRecorder) EnQueue(event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnQueue", reflect.TypeOf((*MockEventQueue)(nil).EnQueue), event)
}

// Queue mocks base method.
func (m *MockEventQueue) Queue() []any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Queue")
	ret0, _ := ret[0].([]any)
	return ret0
}

// Queue indicates an expected call of Queue.
func (mr *MockEventQueueMockRecorder) Queue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Queue", reflect.TypeOf((*MockEventQueue)(nil).Queue))
}
