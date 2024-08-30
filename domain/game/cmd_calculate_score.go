package game

import (
	"context"
	"moj/domain/pkg/queue"
	"moj/domain/record"
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
	gameRepo   GameRepository
	recordRepo record.RecordRepository
}

func NewCalculateScoreCmdHandler(
	gameRepo GameRepository,
	recordRepo record.RecordRepository,
) *CalculateScoreCmdHandler {
	return &CalculateScoreCmdHandler{
		gameRepo:   gameRepo,
		recordRepo: recordRepo,
	}
}

func (h *CalculateScoreCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd CalculateScoreCmd) error {
	game, err := h.gameRepo.FindGameByID(ctx, cmd.GameID)
	if err != nil {
		return err
	}
	// find best record to calulate score
	red, err := h.recordRepo.FindBestRecord(ctx, cmd.AccountID, cmd.QuestionID, cmd.GameID)
	if err != nil {
		return err
	}
	cmd.LastMostFinishedAt = red.NumberFinishedAt
	return game.calculate(queue, cmd)
}
