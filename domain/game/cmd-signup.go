package game

import "moj/domain/pkg/queue"

type SignUpGameCmd struct {
	GameID    int
	AccountID int
	Time      int64
}

type SignupGameCmdHandler struct {
	repo GameRepository
}

func (h *SignupGameCmdHandler) Handle(queue queue.EventQueue, cmd SignUpGameCmd) error {
	game, err := h.repo.findGameByID(cmd.GameID)
	if err != nil {
		return err
	}

	return game.signUp(queue, h.repo.insertSignUpAccount, cmd)
}
