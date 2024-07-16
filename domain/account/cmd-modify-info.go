package account

import "moj/domain/pkg/queue"

type ModifyInfoAccountCmd struct {
	AccountID  int
	NickName   string
	AvatarLink string
}

type ModifyInfoAccountCmdHandler struct {
	repo AccountRepository
}

func NewModifyInfoAccountCmdHandler(repo AccountRepository) *ModifyInfoAccountCmdHandler {
	return &ModifyInfoAccountCmdHandler{
		repo: repo,
	}
}

func (m *ModifyInfoAccountCmdHandler) Handle(queue queue.EventQueue, cmd ModifyInfoAccountCmd) error {
	account, err := m.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	err = account.modifyInfo(queue, cmd)
	if err != nil {
		return err
	}
	return m.repo.Save(account)
}
