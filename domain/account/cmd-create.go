package account

import (
	"moj/domain/pkg/crypt"
	"moj/domain/pkg/queue"
)

type CreateAccountCmd struct {
	Email    string
	NickName string
	Password string
	Time     int64
}

type CreateAccountCmdHandler struct {
	repo  AccountRepo
	crypt crypt.Cryptor
}

func (c *CreateAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd CreateAccountCmd) error {
	acc, err := NewAccount(c.crypt, cmd.Email, cmd.Password, cmd.NickName)
	if err != nil {
		return err
	}

	err = c.repo.save(acc)
	if err != nil {
		return err
	}
	return acc.create(queue, cmd)
}
