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

func NewSignupGameCmdHandler(repo GameRepository) *SignupGameCmdHandler {
	return &SignupGameCmdHandler{
		repo: repo,
	}
}

func (h *SignupGameCmdHandler) Handle(queue queue.EventQueue, cmd SignUpGameCmd) error {
	game, err := h.repo.FindGameByID(cmd.GameID)
	if err != nil {
		return err
	}
	if game == nil {
		return ErrGameNotFound
	}
	return game.signUp(queue, h.repo.InsertSignUpAccount, cmd)
}
