package game

import (
	"context"
	"moj/domain/pkg/queue"
)

type CalculateScoreCmd struct {
	GameID             string
	RecordID           string
	AccountID          string
	QuestionID         string
	NumberFinishedAt   int
	LastMostFinishedAt int
	TotalQuestion      int
}

type CalculateScoreCmdHandler struct {
	repo GameRepository
}

func NewCalculateScoreCmdHandler(repo GameRepository) *CalculateScoreCmdHandler {
	return &CalculateScoreCmdHandler{
		repo: repo,
	}
}

func (h *CalculateScoreCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd CalculateScoreCmd) error {
	game, err := h.repo.FindGameByID(ctx, cmd.GameID)
	if err != nil {
		return err
	}
	return game.calculate(queue, cmd)
}
