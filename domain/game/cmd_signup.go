package game

import "moj/domain/pkg/queue"

type SignUpGameCmd struct {
	GameID    string
	AccountID string
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
	return game.signUp(queue, h.repo.InsertSignUpAccount, cmd)
}
