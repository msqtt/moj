package account

import "moj/domain/pkg/queue"

type SetStatusAccountCmd struct {
	AccountID int
	Enabled   bool
}

type SetStatusAccountCmdHandler struct {
	repo AccountRepository
}

func NewSetStatusAccountCmdHandler(repo AccountRepository) *SetStatusAccountCmdHandler {
	return &SetStatusAccountCmdHandler{
		repo: repo,
	}
}

func (s *SetStatusAccountCmdHandler) Handle(queue queue.EventQueue,
	cmd SetStatusAccountCmd) error {
	acc, err := s.repo.FindAccountByID(cmd.AccountID)
	if err != nil {
		return err
	}

	err = acc.SetStatus(queue, cmd)
	if err != nil {
		return err
	}

	return s.repo.Save(acc)
}
