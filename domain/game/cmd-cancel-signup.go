package game

import "moj/domain/pkg/queue"

type CancelSignUpGameCmd struct {
	GameID    int
	AccountID int
	Time      int64
}

type CancelSignUpGameCmdHandler struct {
	repo GameRepository
}

func (h *CancelSignUpGameCmdHandler) Handle(queue queue.EventQueue, cmd CancelSignUpGameCmd) error {
	game, err := h.repo.findGameByID(cmd.GameID)
	if err != nil {
		return err
	}

	return game.cancelSignUp(queue, h.repo.deletSignUpAccount, cmd)
}
