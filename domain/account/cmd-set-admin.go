package account

import "moj/domain/pkg/queue"

type SetAdminAccountCmd struct {
	AccountID int
	IsAdmin   bool
}

type SetAdminAccountCmdHandler struct {
	repo AccountRepo
}

func (s *SetAdminAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd SetAdminAccountCmd) error {
	acc, err := s.repo.findAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	err = acc.setAdmin(queue, cmd)
	if err != nil {
		return err
	}

	return s.repo.save(acc)
}
