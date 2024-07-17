package account

import "moj/domain/pkg/queue"

type SetAdminAccountCmd struct {
	AccountID int
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

func (s *SetAdminAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd SetAdminAccountCmd) error {
	acc, err := s.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountNotFound
	}
	err = acc.SetAdmin(queue, cmd)
	if err != nil {
		return err
	}

	return s.repo.Save(acc)
}
