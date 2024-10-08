package account

import (
	"context"
	"moj/domain/pkg/queue"
)

type SetAdminAccountCmd struct {
	AccountID string
	IsAdmin   bool
}

type SetAdminAccountCmdHandler struct {
	repo AccountRepository
}

func NewSetAdminAccountCmdHandler(repo AccountRepository) *SetAdminAccountCmdHandler {
	return &SetAdminAccountCmdHandler{
		repo: repo,
	}
}

func (s *SetAdminAccountCmdHandler) Handle(
	ctx context.Context,
	queue queue.EventQueue,
	cmd SetAdminAccountCmd) error {
	acc, err := s.repo.FindAccountByID(ctx, cmd.AccountID)
	if err != nil {
		return err
	}
	err = acc.SetAdmin(queue, cmd)
	if err != nil {
		return err
	}

	return s.repo.Save(ctx, acc)
}
