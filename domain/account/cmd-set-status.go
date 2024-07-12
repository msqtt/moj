package account

import "moj/domain/pkg/queue"

type SetStatusAccountCmd struct {
	AccountID int
	Enabled   bool
}

type SetStatusAccountCmdHandler struct {
	repo AccountRepo
}

func (s *SetStatusAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd SetStatusAccountCmd) error {
	acc, err := s.repo.findAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	err = acc.setStatus(queue, cmd)
	if err != nil {
		return err
	}

	return s.repo.save(acc)
}
