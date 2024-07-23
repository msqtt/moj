package account

import "moj/domain/pkg/queue"

type LoginAccountCmd struct {
	AccountID string
	Device    string
	IPAddr    string
	Time      int64
}

type LoginAccountCmdHandler struct {
	repo AccountRepository
}

func NewLoginAccountCmdHandler(repo AccountRepository) *LoginAccountCmdHandler {
	return &LoginAccountCmdHandler{
		repo: repo,
	}
}

func (l *LoginAccountCmdHandler) Handle(queue queue.EventQueue, cmd LoginAccountCmd) error {
	acc, err := l.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	return acc.login(queue, cmd)
}
