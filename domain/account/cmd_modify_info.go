package account

import (
	"context"
	"moj/domain/pkg/queue"
)

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

func (m *ModifyInfoAccountCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd ModifyInfoAccountCmd) error {
	acc, err := m.repo.FindAccountByID(ctx, cmd.AccountID)
	if err != nil {
		return err
	}
	err = acc.modifyInfo(queue, cmd)
	if err != nil {
		return err
	}
	return m.repo.Save(ctx, acc)
}
