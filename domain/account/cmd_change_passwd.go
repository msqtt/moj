package account

import (
	"context"
	"moj/domain/pkg/crypt"
	"moj/domain/pkg/queue"
)

type ChangePasswdAccountCmd struct {
	AccountID string
	Password  string
	Time      int64
}

type ChangePasswdAccountCmdHandler struct {
	repo  AccountRepository
	crypt crypt.Cryptor
}

func NewChangePasswdAccountCmdHandler(repo AccountRepository, crypt crypt.Cryptor) *ChangePasswdAccountCmdHandler {
	return &ChangePasswdAccountCmdHandler{
		repo:  repo,
		crypt: crypt,
	}
}

func (h *ChangePasswdAccountCmdHandler) Handle(ctx context.Context, queue queue.EventQueue,
	cmd ChangePasswdAccountCmd) error {
	acc, err := h.repo.FindAccountByID(ctx, cmd.AccountID)
	if err != nil {
		return err
	}
	err = acc.changePasswd(h.crypt, queue, cmd)
	if err != nil {
		return err
	}
	return h.repo.Save(ctx, acc)
}
