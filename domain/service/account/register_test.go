package account_test

import (
	"context"
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
	s := saccount.NewAccountRegisterService(accHandler, mCRepo, mARepo)

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
		FindLatestCaptcha(gomock.Any(), gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	mCRepo.EXPECT().
		Save(gomock.Any(), gomock.Eq(cap)).
		Return(nil)

	mCryp.EXPECT().
		Encrypt(gomock.Eq(cmd.Password)).
		Return(cmd.Password, nil)

	mARepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	mARepo.EXPECT().FindAccountByEmail(gomock.Any(), gomock.Eq(cmd.Email)).Return(nil, account.ErrAccountNotFound)

	mQueue.EXPECT().EnQueue(gomock.Eq(event)).Return(nil)

	_, err = s.Handle(context.TODO(), mQueue, cmd)
	require.NoError(t, err)

	// Test case 2: Captcha not found
	mARepo.EXPECT().FindAccountByEmail(gomock.Any(), gomock.Eq(cmd.Email)).Return(nil, account.ErrAccountNotFound)
	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Any(), gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(nil, saccount.ErrCaptchaNotFound)

	_, err = s.Handle(context.TODO(), mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaNotFound)

	// // Test case 3: Captcha expired
	cap, err = captcha.NewCaptcha("", cmd.Email, captcha.CaptchaTypeRegister,
		"IP_ADDRESS", 0, time.Now().Unix()-1000)
	require.NoError(t, err)
	require.NotNil(t, cap)

	cap.Code = cmd.Captcha

	mARepo.EXPECT().FindAccountByEmail(gomock.Any(), gomock.Eq(cmd.Email)).Return(nil, account.ErrAccountNotFound)
	mCRepo.EXPECT().
		FindLatestCaptcha(gomock.Any(), gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	_, err = s.Handle(context.TODO(), mQueue, cmd)
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
		FindLatestCaptcha(gomock.Any(), gomock.Eq(cmd.Email),
			gomock.Eq(cmd.Captcha), gomock.Eq(captcha.CaptchaTypeRegister)).
		Return(cap, nil)

	mCRepo.EXPECT().Save(gomock.Any(), gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		Save(context.TODO(), gomock.Any()).
		Return(errors.New("Failed to create account"))

	mARepo.EXPECT().FindAccountByEmail(gomock.Any(), gomock.Eq(cmd.Email)).Return(nil, account.ErrAccountNotFound)
	_, err = s.Handle(context.TODO(), mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrFailedToCreateAccount)
}
