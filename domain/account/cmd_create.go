package account

import (
	"context"
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
	repo  AccountRepository
	crypt crypt.Cryptor
}

func NewCreateAccountCmdHandler(repo AccountRepository, crypt crypt.Cryptor) *CreateAccountCmdHandler {
	return &CreateAccountCmdHandler{
		repo:  repo,
		crypt: crypt,
	}
}

func (c *CreateAccountCmdHandler) Handle(ctx context.Context, queue queue.EventQueue,
	cmd CreateAccountCmd) (string, error) {
	acc, err := NewAccount(c.crypt, cmd.Email, cmd.Password, cmd.NickName)
	if err != nil {
		return "", err
	}

	err = c.repo.Save(ctx, acc)
	if err != nil {
		return "", err
	}
	return acc.AccountID, acc.create(queue, cmd)
}
