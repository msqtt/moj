package account_test

import (
	"errors"
	"testing"
	"time"

	"moj/domain/account"
	"moj/domain/captcha"
	saccount "moj/domain/service/account"
	mock_account "moj/domain/service/account/mock"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	mQueue := mock_account.NewMockEventQueue(ctrl)
	mARepo := mock_account.NewMockAccountRepository(ctrl)
	mCRepo := mock_account.NewMockCaptchaRepository(ctrl)
	mCryp := mock_account.NewMockCryptor(ctrl)

	// Create a new AccountRegisterService
	accHandler := account.NewCreateAccountCmdHandler(mARepo, mCryp)
	s := saccount.NewAccountRegisterService(*accHandler, mCRepo)

	// Test case 1: Successful registration
	cmd := saccount.RegisterCmd{
		Email:    "test@example.com",
		NickName: "testNickName",
		Password: "test!@#Pas123d",
		Captcha:  "123456",
		Time:     time.Now().Unix(),
	}
	cap, err := captcha.NewCaptcha("", cmd.Email, captcha.CaptchaTypeRegister,
		"IP_ADDRESS", 10000, time.Now().Unix())
	require.NoError(t, err)
	require.NotNil(t, cap)

	cap.Code = cmd.Captcha

	event := account.CreateAccountEvent{
		AccountID:    "",
		Email:        cmd.Email,
		NickName:     cmd.NickName,
		RegisterTime: cmd.Time,
		Enabled:      true,
	}

	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	mCRepo.EXPECT().
		Save(gomock.Eq(cap)).
		Return(nil)

	mCryp.EXPECT().
		Encrypt(gomock.Eq(cmd.Password)).
		Return(cmd.Password, nil)

	mARepo.EXPECT().
		Save(gomock.Any()).
		Return(nil)

	mQueue.EXPECT().EnQueue(gomock.Eq(event)).Return(nil)

	err = s.Handle(mQueue, cmd)
	require.NoError(t, err)

	// Test case 2: Captcha not found
	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(nil, nil)

	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaNotFound)

	// // Test case 3: Captcha expired
	cap, err = captcha.NewCaptcha("", cmd.Email, captcha.CaptchaTypeRegister,
		"IP_ADDRESS", 0, time.Now().Unix()-1000)
	require.NoError(t, err)
	require.NotNil(t, cap)

	cap.Code = cmd.Captcha

	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaAlreadyExpired)

	// Test case 4: Failed to create account
	cap, err = captcha.NewCaptcha("", cmd.Email, captcha.CaptchaTypeRegister,
		"IP_ADDRESS", 10000, time.Now().Unix())

	require.NoError(t, err)
	require.NotNil(t, cap)

	cap.Code = cmd.Captcha

	mCryp.EXPECT().
		Encrypt(gomock.Eq(cmd.Password)).
		Return(cmd.Password, nil)

	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	mCRepo.EXPECT().Save(gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		Save(gomock.Any()).
		Return(errors.New("Failed to create account"))

	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrFailedToCreateAccount)
}
