package account

import "moj/domain/pkg/queue"

type DeleteAccountCmd struct {
	AccountID int
	Time      int64
}

type DeleteAccountCmdHandler struct {
	repo AccountRepository
}

func (d *DeleteAccountCmdHandler) Handle(queue queue.EventQueue, cmd DeleteAccountCmd) error {
	acc, err := d.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	err = acc.delete(queue, cmd)
	if err != nil {
		return err
	}
	return d.repo.Save(acc)
}
