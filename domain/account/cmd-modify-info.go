package account

import "github.com/msqtt/moj/domain/pkg/queue"

type ModifyInfoAccountCmd struct {
	AccountID  string
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
	acc, err := m.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountNotFound
	}
	err = acc.modifyInfo(queue, cmd)
	if err != nil {
		return err
	}
	return m.repo.Save(acc)
}
