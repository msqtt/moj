package game

import "moj/domain/pkg/queue"

type CreateGameCmd struct {
	Title        string
	Description  string
	AccountID    int
	StartTime    int64
	EndTime      int64
	QuestionList []GameQuestion
	Time         int64
}

type CreateGameCmdHandler struct {
	repo GameRepository
}

func (h *CreateGameCmdHandler) Handle(queue queue.EventQueue, cmd CreateGameCmd) error {
	game, err := NewGame(cmd.AccountID, cmd.Title, cmd.Description, cmd.Time,
		cmd.StartTime, cmd.EndTime, cmd.QuestionList)
	if err != nil {
		return err
	}
	return h.repo.Save(game)
}
