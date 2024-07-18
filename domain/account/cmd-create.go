package account

import (
	"github.com/msqtt/moj/domain/pkg/crypt"
	"github.com/msqtt/moj/domain/pkg/queue"
)

type CreateAccountCmd struct {
	Email    string
	NickName string
	Password string
	Time     int64
}

type CreateAccountCmdHandler struct {
	repo  AccountRepository
	crypt crypt.Cryptor
}

func NewCreateAccountCmdHandler(repo AccountRepository, crypt crypt.Cryptor) *CreateAccountCmdHandler {
	return &CreateAccountCmdHandler{
		repo:  repo,
		crypt: crypt,
	}
}

func (c *CreateAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd CreateAccountCmd) error {
	acc, err := NewAccount(c.crypt, cmd.Email, cmd.Password, cmd.NickName)
	if err != nil {
		return err
	}

	err = c.repo.Save(acc)
	if err != nil {
		return err
	}
	return acc.create(queue, cmd)
}
