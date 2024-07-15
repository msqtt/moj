package account

import "moj/domain/pkg/queue"

type LoginAccountCmd struct {
	AccountID int
	Device    string
	IPAddr    string
	Time      int64
}

type LoginAccountCmdHandler struct {
	repo AccountRepository
}

func (l *LoginAccountCmdHandler) Handle(queue queue.EventQueue, cmd LoginAccountCmd) error {
	account, err := l.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	return account.login(queue, cmd)
}
