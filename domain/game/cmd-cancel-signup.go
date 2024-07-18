package game

import "github.com/msqtt/moj/domain/pkg/queue"

type CancelSignUpGameCmd struct {
	GameID    int
	AccountID int
	Time      int64
}

type CancelSignUpGameCmdHandler struct {
	repo GameRepository
}

func NewCancelSignUpGameCmdHandler(repo GameRepository) *CancelSignUpGameCmdHandler {
	return &CancelSignUpGameCmdHandler{
		repo: repo,
	}
}

func (h *CancelSignUpGameCmdHandler) Handle(queue queue.EventQueue, cmd CancelSignUpGameCmd) error {
	game, err := h.repo.FindGameByID(cmd.GameID)
	if err != nil {
		return err
	}
	if game == nil {
		return ErrGameNotFound
	}
	return game.cancelSignUp(queue, h.repo.DeletSignUpAccount, cmd)
}
