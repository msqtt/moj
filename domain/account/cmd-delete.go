package account

import "moj/domain/pkg/queue"

type DeleteAccountCmd struct {
	AccountID int
	Time      int64
}

type DeleteAccountCmdHandler struct {
	repo AccountRepository
}

func NewDeleteAccountCmdHandler(repo AccountRepository) *DeleteAccountCmdHandler {
	return &DeleteAccountCmdHandler{
		repo: repo,
	}
}

func (d *DeleteAccountCmdHandler) Handle(queue queue.EventQueue, cmd DeleteAccountCmd) error {
	acc, err := d.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	if acc == nil {
		return ErrAccountNotFound
	}

	err = acc.delete(queue, cmd)
	if err != nil {
		return err
	}
	return d.repo.Save(acc)
}
