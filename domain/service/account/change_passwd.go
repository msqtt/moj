package account

import (
	"errors"

	"github.com/msqtt/moj/domain/account"
	"github.com/msqtt/moj/domain/captcha"
	"github.com/msqtt/moj/domain/pkg/queue"
)

var ErrFailedToChangePasswd = errors.New("failed to change password")

type ChangePasswdCmd struct {
	AccountID string
	Email     string
	Password  string
	Captcha   string
	Time      int64
}

type ChangePasswdService struct {
	changePasswdAccountCmdHandler account.ChangePasswdAccountCmdHandler
	captchaRepository             captcha.CaptchaRepository
}

func NewChangePasswdService(changePasswdAccountCmdHandler account.ChangePasswdAccountCmdHandler,
	captchaRepository captcha.CaptchaRepository) *ChangePasswdService {
	return &ChangePasswdService{
		changePasswdAccountCmdHandler: changePasswdAccountCmdHandler,
		captchaRepository:             captchaRepository,
	}
}

func (s *ChangePasswdService) Handle(queue queue.EventQueue, cmd ChangePasswdCmd) error {
	cap, err := s.captchaRepository.FindLatestCaptcha(cmd.Email, cmd.Captcha,
		captcha.CaptchaTypeChangePasswd)
	if err != nil {
		return err
	}
	if cap == nil {
		return ErrCaptchaNotFound
	}
	if cap.IsExpired(cmd.Time) {
		return ErrCaptchaAlreadyExpired
	}

	cap.SetDisable()
	s.captchaRepository.Save(cap)

	changePasswdAccountCmd := account.ChangePasswdAccountCmd{
		AccountID: cmd.AccountID,
		Password:  cmd.Password,
		Time:      cmd.Time,
	}
	err = s.changePasswdAccountCmdHandler.Handle(queue, changePasswdAccountCmd)
	if err != nil {
		return errors.Join(ErrFailedToChangePasswd, err)
	}
	return nil
}