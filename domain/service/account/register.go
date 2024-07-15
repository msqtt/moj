package account

import (
	"errors"
	"moj/domain/account"
	"moj/domain/captcha"
	"moj/domain/pkg/queue"
	"time"
)

var (
	ErrCaptchaNotFound       = errors.New("captcha not found")
	ErrCaptchaAlreadyExpired = errors.New("captcha already expired")
	ErrFailedToCreateAccount = errors.New("failed to create account")
)

type RegisterCmd struct {
	Email    string
	NickName string
	Password string
	Captcha  string
}

type AccountRegisterService struct {
	createAccountCmdHandler account.CreateAccountCmdHandler
	captchaRepository       captcha.CaptchaRepository
}

func NewAccountRegisterService(createAccountCmdHandler account.CreateAccountCmdHandler,
	captchaRepository captcha.CaptchaRepository) *AccountRegisterService {
	return &AccountRegisterService{createAccountCmdHandler: createAccountCmdHandler, captchaRepository: captchaRepository}
}

func (s *AccountRegisterService) Register(queue queue.EventQueue, cmd RegisterCmd) error {
	cap, err := s.captchaRepository.FindLatestCaptcha(cmd.Email, cmd.Captcha,
		captcha.CaptchaTypeRegister)
	if err != nil {
		return err
	}
	if cap == nil {
		return ErrCaptchaNotFound
	}
	if cap.IsExpired(time.Now().Unix()) {
		return ErrCaptchaAlreadyExpired
	}
	cmdCreateAccount := account.CreateAccountCmd{
		Email:    cmd.Email,
		NickName: cmd.NickName,
		Password: cmd.Password,
	}
	err = s.createAccountCmdHandler.Handle(queue, cmdCreateAccount)
	if err != nil {
		return errors.Join(ErrFailedToCreateAccount, err)
	}
	return nil
}
