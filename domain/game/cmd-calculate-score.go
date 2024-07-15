package game

import "moj/domain/pkg/queue"

type CalculateScoreCmd struct {
	GameID           int
	RecordID         int
	AccountID        int
	QuestionID       int
	NumberFinishedAt int
	LastFinishedAt   int
	TotalQuestion    int
}

type CalculateScoreCmdHandler struct {
	repo GameRepository
}

func (h *CalculateScoreCmdHandler) Handle(queue queue.EventQueue, cmd CalculateScoreCmd) error {
	game, err := h.repo.FindGameByID(cmd.AccountID)
	if err != nil {
		return err
	}

	return game.calculate(queue, cmd)
}
