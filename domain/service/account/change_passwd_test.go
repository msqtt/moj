package account_test

import (
	"context"
	"testing"
	"time"

	"moj/domain/account"
	"moj/domain/captcha"
	saccount "moj/domain/service/account"
	mock_account "moj/domain/service/account/mock"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestChangePasswd(t *testing.T) {
	ctrl := gomock.NewController(t)
	mQueue := mock_account.NewMockEventQueue(ctrl)
	mARepo := mock_account.NewMockAccountRepository(ctrl)
	mCRepo := mock_account.NewMockCaptchaRepository(ctrl)
	mCryp := mock_account.NewMockCryptor(ctrl)

	// Create a new AccountRegisterService
	accHandler := account.NewChangePasswdAccountCmdHandler(mARepo, mCryp)
	s := saccount.NewChangePasswdService(accHandler, mCRepo)

	// Test case 1: Successful handling
	cmd := saccount.ChangePasswdCmd{
		Email:     "test@example.com",
		Captcha:   "123456",
		AccountID: "1",
		Password:  "newPassword!@#123",
		Time:      time.Now().Unix(),
	}

	event := account.ChangePasswdAccountEvent{
		AccountID:  cmd.AccountID,
		ChangeTime: cmd.Time,
	}

	cap, err := captcha.NewCaptcha(cmd.AccountID, cmd.Email, captcha.CaptchaTypeChangePasswd,
		"IP_ADDRESS", 10000, time.Now().Unix())

	require.NotEmpty(t, cap)
	require.NoError(t, err)

	cap.Code = cmd.Captcha

	mCRepo.EXPECT().
		FindLatestCaptcha(context.TODO(), cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap, nil)

	mCRepo.EXPECT().Save(gomock.Any(), gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		FindAccountByID(context.TODO(), cmd.AccountID).
		Return(&account.Account{
			AccountID: cmd.AccountID,
		}, nil)

	mARepo.EXPECT().Save(context.TODO(), gomock.Any()).Return(nil)

	mCryp.EXPECT().
		Encrypt(gomock.Eq(cmd.Password)).
		Return(cmd.Password, nil)

	mQueue.EXPECT().
		EnQueue(gomock.Eq(event)).
		Return(nil)

	err = s.Handle(context.TODO(), mQueue, cmd)
	require.NoError(t, err)

	// Test case 2: Captcha not found
	mCRepo.EXPECT().
		FindLatestCaptcha(context.TODO(), cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(nil, saccount.ErrCaptchaNotFound)

	err = s.Handle(context.TODO(), mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaNotFound)

	// Test case 3: Captcha expired
	cap2, err := captcha.NewCaptcha(cmd.AccountID, cmd.Email, captcha.CaptchaTypeChangePasswd,
		"IP_ADDRESS", 0, time.Now().Unix()-1)

	mCRepo.EXPECT().
		FindLatestCaptcha(context.TODO(), cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap2, nil)
	require.NotEmpty(t, cap2)
	require.NoError(t, err)

	cap.Code = cmd.Captcha

	// Set up the mock captcha repository to return an expired captcha
	err = s.Handle(context.TODO(), mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaAlreadyExpired)

	// Test case 4: Failed to change password
	cap.Enabled = true
	mCRepo.EXPECT().
		FindLatestCaptcha(context.TODO(), cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap, nil)

	mCRepo.EXPECT().Save(context.TODO(), gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		FindAccountByID(context.TODO(), cmd.AccountID).
		Return(nil, account.ErrAccountNotFound)

	err = s.Handle(context.TODO(), mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrFailedToChangePasswd)
}
