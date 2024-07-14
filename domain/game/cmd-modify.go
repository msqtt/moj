package game

import "moj/domain/pkg/queue"

type ModifyGameCmd struct {
	GameID       int
	Title        string
	Descirption  string
	StartTime    int64
	EndTime      int64
	QuestionList []GameQuestion
}

type ModifyGameCmdHandler struct {
	repo GameRepository
}

func (h *ModifyGameCmdHandler) Handle(queue queue.EventQueue, cmd ModifyGameCmd) error {
	game, err := h.repo.findGameByID(cmd.GameID)
	if err != nil {
		return err
	}

	err = game.modify(cmd)
	if err != nil {
		return err
	}

	return h.repo.save(game)
}
