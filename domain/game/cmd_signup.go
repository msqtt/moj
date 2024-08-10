package game

import (
	"context"
	"moj/domain/pkg/queue"
)

type SignUpGameCmd struct {
	GameID    string
	AccountID string
	Time      int64
}

type SignupGameCmdHandler struct {
	repo GameRepository
}

func NewSignUpGameCmdHandler(repo GameRepository) *SignupGameCmdHandler {
	return &SignupGameCmdHandler{
		repo: repo,
	}
}

func (h *SignupGameCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd SignUpGameCmd) error {
	game, err := h.repo.FindGameByID(ctx, cmd.GameID)
	if err != nil {
		return err
	}
	return game.signUp(queue, func(gid, aid string, time int64) error {
		return h.repo.InsertSignUpAccount(ctx, gid, aid, time)
	}, cmd)
}
