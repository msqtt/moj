package account

import (
	"errors"
	"moj/domain/account"
	"moj/domain/captcha"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
)

var (
	ErrCaptchaNotFound       = errors.New("captcha not found")
	ErrCaptchaAlreadyExpired = errors.Join(domain_err.ErrExpired, errors.New("captcha already expired"))
	ErrFailedToCreateAccount = errors.New("failed to create account")
)

type RegisterCmd struct {
	Email    string
	NickName string
	Password string
	Captcha  string
	Time     int64
}

type AccountRegisterService struct {
	createAccountCmdHandler *account.CreateAccountCmdHandler
	captchaRepository       captcha.CaptchaRepository
}

func NewAccountRegisterService(createAccountCmdHandler *account.CreateAccountCmdHandler,
	captchaRepository captcha.CaptchaRepository) *AccountRegisterService {
	return &AccountRegisterService{createAccountCmdHandler: createAccountCmdHandler,
		captchaRepository: captchaRepository}
}

func (s *AccountRegisterService) Handle(queue queue.EventQueue, cmd RegisterCmd) error {
	cap, err := s.captchaRepository.FindLatestCaptcha(cmd.Email, cmd.Captcha,
		captcha.CaptchaTypeRegister)
	if err != nil {
		return err
	}
	if cap.IsExpired(cmd.Time) {
		return ErrCaptchaAlreadyExpired
	}

	cap.SetDisable()
	s.captchaRepository.Save(cap)

	cmdCreateAccount := account.CreateAccountCmd{
		Email:    cmd.Email,
		NickName: cmd.NickName,
		Password: cmd.Password,
		Time:     cmd.Time,
	}
	err = s.createAccountCmdHandler.Handle(queue, cmdCreateAccount)
	if err != nil {
		return errors.Join(ErrFailedToCreateAccount, err)
	}
	return nil
}
