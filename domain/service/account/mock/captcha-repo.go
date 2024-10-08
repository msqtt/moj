// Code generated by MockGen. DO NOT EDIT.
// Source: domain/captcha/repotitory.go
//
// Generated by this command:
//
//	mockgen -source domain/captcha/repotitory.go -destination domain/service/account/mock/captcha-repo.go -package mock_account
//

// Package mock_account is a generated GoMock package.
package mock_account

import (
	context "context"
	captcha "moj/domain/captcha"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCaptchaRepository is a mock of CaptchaRepository interface.
type MockCaptchaRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCaptchaRepositoryMockRecorder
}

// MockCaptchaRepositoryMockRecorder is the mock recorder for MockCaptchaRepository.
type MockCaptchaRepositoryMockRecorder struct {
	mock *MockCaptchaRepository
}

// NewMockCaptchaRepository creates a new mock instance.
func NewMockCaptchaRepository(ctrl *gomock.Controller) *MockCaptchaRepository {
	mock := &MockCaptchaRepository{ctrl: ctrl}
	mock.recorder = &MockCaptchaRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCaptchaRepository) EXPECT() *MockCaptchaRepositoryMockRecorder {
	return m.recorder
}

// FindLatestCaptcha mocks base method.
func (m *MockCaptchaRepository) FindLatestCaptcha(ctx context.Context, email, code string, captchaType captcha.CaptchaType) (*captcha.Captcha, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindLatestCaptcha", ctx, email, code, captchaType)
	ret0, _ := ret[0].(*captcha.Captcha)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindLatestCaptcha indicates an expected call of FindLatestCaptcha.
func (mr *MockCaptchaRepositoryMockRecorder) FindLatestCaptcha(ctx, email, code, captchaType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindLatestCaptcha", reflect.TypeOf((*MockCaptchaRepository)(nil).FindLatestCaptcha), ctx, email, code, captchaType)
}

// Save mocks base method.
func (m *MockCaptchaRepository) Save(ctx context.Context, captcha *captcha.Captcha) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, captcha)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCaptchaRepositoryMockRecorder) Save(ctx, captcha any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCaptchaRepository)(nil).Save), ctx, captcha)
}
