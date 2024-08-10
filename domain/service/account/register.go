package account

import (
	"context"
	"errors"
	"moj/domain/account"
	"moj/domain/captcha"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
)

var (
	ErrAccountAlreadyRegistered = errors.Join(domain_err.ErrDuplicated, errors.New("account already been registered"))
	ErrCaptchaNotFound          = errors.New("captcha not found")
	ErrCaptchaAlreadyExpired    = errors.Join(domain_err.ErrExpired, errors.New("captcha already expired"))
	ErrFailedToCreateAccount    = errors.New("failed to create account")
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
	accountRepository       account.AccountRepository
}

func NewAccountRegisterService(createAccountCmdHandler *account.CreateAccountCmdHandler,
	captchaRepository captcha.CaptchaRepository,
	accountRepository account.AccountRepository,
) *AccountRegisterService {
	return &AccountRegisterService{
		createAccountCmdHandler: createAccountCmdHandler,
		captchaRepository:       captchaRepository,
		accountRepository:       accountRepository,
	}
}

func (s *AccountRegisterService) Handle(ctx context.Context, queue queue.EventQueue, cmd RegisterCmd) error {
	// check the account by email
	_, err := s.accountRepository.FindAccountByEmail(ctx, cmd.Email)
	if err == nil {
		return ErrAccountAlreadyRegistered
	} else if !errors.Is(err, account.ErrAccountNotFound) {
		return err
	}

	cap, err := s.captchaRepository.FindLatestCaptcha(ctx, cmd.Email, cmd.Captcha,
		captcha.CaptchaTypeRegister)
	if err != nil {
		return err
	}
	if cap.IsExpired(cmd.Time) {
		return ErrCaptchaAlreadyExpired
	}

	cap.SetDisable()
	s.captchaRepository.Save(ctx, cap)

	cmdCreateAccount := account.CreateAccountCmd{
		Email:    cmd.Email,
		NickName: cmd.NickName,
		Password: cmd.Password,
		Time:     cmd.Time,
	}
	err = s.createAccountCmdHandler.Handle(ctx, queue, cmdCreateAccount)
	if err != nil {
		return errors.Join(ErrFailedToCreateAccount, err)
	}
	return nil
}
