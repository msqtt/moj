package account_test

import (
	"github.com/msqtt/moj/domain/account"
	"github.com/msqtt/moj/domain/captcha"
	saccount "github.com/msqtt/moj/domain/service/account"
	mock_account "github.com/msqtt/moj/domain/service/account/mock"
	"testing"
	"time"

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
	s := saccount.NewChangePasswdService(*accHandler, mCRepo)

	// Test case 1: Successful handling
	cmd := saccount.ChangePasswdCmd{
		Email:     "test@example.com",
		Captcha:   "123456",
		AccountID: 1,
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
		FindLatestCaptcha(cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap, nil)

	mCRepo.EXPECT().Save(gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		FindAccountByID(cmd.AccountID).
		Return(&account.Account{
			AccountID: cmd.AccountID,
		}, nil)

	mARepo.EXPECT().Save(gomock.Any()).Return(nil)

	mCryp.EXPECT().
		Encrypt(gomock.Eq(cmd.Password)).
		Return(cmd.Password)

	mQueue.EXPECT().
		EnQueue(gomock.Eq(event)).
		Return(nil)

	err = s.Handle(mQueue, cmd)
	require.NoError(t, err)

	// Test case 2: Captcha not found
	mCRepo.EXPECT().
		FindLatestCaptcha(cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(nil, nil)

	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaNotFound)

	// Test case 3: Captcha expired
	cap2, err := captcha.NewCaptcha(cmd.AccountID, cmd.Email, captcha.CaptchaTypeChangePasswd,
		"IP_ADDRESS", 0, time.Now().Unix()-1)

	mCRepo.EXPECT().
		FindLatestCaptcha(cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap2, nil)
	require.NotEmpty(t, cap2)
	require.NoError(t, err)

	cap.Code = cmd.Captcha

	// Set up the mock captcha repository to return an expired captcha
	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrCaptchaAlreadyExpired)

	// Test case 4: Failed to change password
	mCRepo.EXPECT().
		FindLatestCaptcha(cmd.Email, cmd.Captcha, captcha.CaptchaTypeChangePasswd).
		Return(cap, nil)

	mCRepo.EXPECT().Save(gomock.Eq(cap)).Return(nil)

	mARepo.EXPECT().
		FindAccountByID(cmd.AccountID).
		Return(nil, nil)

	err = s.Handle(mQueue, cmd)
	require.ErrorIs(t, err, saccount.ErrFailedToChangePasswd)
}
