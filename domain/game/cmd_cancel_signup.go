package game

import (
	"context"
	"moj/domain/pkg/queue"
)

type CancelSignUpGameCmd struct {
	GameID    string
	AccountID string
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

func (h *CancelSignUpGameCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd CancelSignUpGameCmd) error {
	game, err := h.repo.FindGameByID(ctx, cmd.GameID)
	if err != nil {
		return err
	}
	return game.cancelSignUp(queue, func(gid, aid string) error {
		return h.repo.DeletSignUpAccount(ctx, gid, aid)
	}, cmd)
}
