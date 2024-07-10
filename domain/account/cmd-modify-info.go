package account

import "moj/domain/pkg/queue"

type ModifyInfoAccountCmd struct {
	AccountID  int
	NickName   string
	AvatarLink string
}

type ModifyInfoAccountCmdHandler struct {
	repo AccountRepo
}

func (m *ModifyInfoAccountCmdHandler) Handle(queue queue.EventQueue, cmd ModifyInfoAccountCmd) error {
	account, err := m.repo.findAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	err = account.modifyInfo(queue, cmd)
	if err != nil {
		return err
	}
	return m.repo.save(&account)
}
