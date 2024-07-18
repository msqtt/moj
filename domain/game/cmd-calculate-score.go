package game

import "github.com/msqtt/moj/domain/pkg/queue"

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

func NewCalculateScoreCmdHandler(repo GameRepository) *CalculateScoreCmdHandler {
	return &CalculateScoreCmdHandler{
		repo: repo,
	}
}

func (h *CalculateScoreCmdHandler) Handle(queue queue.EventQueue, cmd CalculateScoreCmd) error {
	game, err := h.repo.FindGameByID(cmd.AccountID)
	if err != nil {
		return err
	}
	if game == nil {
		return ErrGameNotFound
	}
	return game.calculate(queue, cmd)
}
