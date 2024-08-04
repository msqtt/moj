package game

import "moj/domain/pkg/queue"

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

func (h *ModifyGameCmdHandler) Handle(queue queue.EventQueue, cmd ModifyGameCmd) (any, error) {
	game, err := h.repo.FindGameByID(cmd.GameID)
	if err != nil {
		return nil, err
	}
	err = game.modify(cmd)
	if err != nil {
		return nil, err
	}

	return h.repo.Save(game)
}
