package game

import (
	"context"
	"moj/domain/pkg/queue"
)

type ModifyGameCmd struct {
	GameID       string
	Title        string
	Descirption  string
	StartTime    int64
	EndTime      int64
	QuestionList []GameQuestion
}

type ModifyGameCmdHandler struct {
	repo GameRepository
}

func NewModifyGameCmdHandler(repo GameRepository) *ModifyGameCmdHandler {
	return &ModifyGameCmdHandler{
		repo: repo,
	}
}

func (h *ModifyGameCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd ModifyGameCmd) (any, error) {
	game, err := h.repo.FindGameByID(ctx, cmd.GameID)
	if err != nil {
		return nil, err
	}
	err = game.modify(cmd)
	if err != nil {
		return nil, err
	}

	return h.repo.Save(ctx, game)
}
