package account

import (
	"github.com/msqtt/moj/domain/pkg/crypt"
	"github.com/msqtt/moj/domain/pkg/queue"
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

func (h *ChangePasswdAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd ChangePasswdAccountCmd) error {
	acc, err := h.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	if acc == nil {
		return ErrAccountNotFound
	}

	err = acc.changePasswd(h.crypt, queue, cmd)
	if err != nil {
		return err
	}
	return h.repo.Save(acc)
}