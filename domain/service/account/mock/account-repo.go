// Code generated by MockGen. DO NOT EDIT.
// Source: domain/account/repository.go
//
// Generated by this command:
//
//	mockgen -source domain/account/repository.go -destination domain/service/account/mock/account-repo.go
//

// Package mock_account is a generated GoMock package.
package mock_account

import (
	account "github.com/msqtt/moj/domain/account"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// FindAccountByID mocks base method.
func (m *MockAccountRepository) FindAccountByID(accountID int) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAccountByID", accountID)
	ret0, _ := ret[0].(*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAccountByID indicates an expected call of FindAccountByID.
func (mr *MockAccountRepositoryMockRecorder) FindAccountByID(accountID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAccountByID", reflect.TypeOf((*MockAccountRepository)(nil).FindAccountByID), accountID)
}

// Save mocks base method.
func (m *MockAccountRepository) Save(arg0 *account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockAccountRepositoryMockRecorder) Save(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAccountRepository)(nil).Save), arg0)
}
