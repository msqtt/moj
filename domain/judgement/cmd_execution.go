package judgement

import (
	"errors"
	"moj/domain/pkg/queue"
)

type Case struct {
	Number       int
	InputContent string
	OuputContent string
}

type ExecutionCmd struct {
	RecordID           string
	QuestionID         string
	QuestionModifyTime int64
	Cases              []Case
	Language           string
	Code               string
	CodeHash           string
	TimeLimit          int64
	MemoryLimit        int64
	Time               int64
}

type ExecutionCmdHandler struct {
	repo       JudgementRepository
	exeService ExecutionService
}

func NewExecutionCmdHandler(repo JudgementRepository, exeService ExecutionService) *ExecutionCmdHandler {
	return &ExecutionCmdHandler{
		repo:       repo,
		exeService: exeService,
	}
}

func (e *ExecutionCmdHandler) Handle(queue queue.EventQueue, cmd ExecutionCmd) error {
	// Check if there are already cached before execution
	jud, err := e.repo.FindJudgementByHash(cmd.QuestionID, cmd.CodeHash, cmd.QuestionModifyTime)
	if err != nil {
		return err
	} else if errors.Is(err, ErrJudgementNotFound) {
		jud = NewJudgement(cmd.RecordID, cmd.QuestionID, len(cmd.Cases),
			cmd.Language, cmd.Code, cmd.CodeHash, cmd.Time)

		err = jud.execute(queue, e.exeService, cmd)
		if err != nil {
			return err
		}
		return e.repo.Save(jud)
	} else {
		return jud.sendExecutionEvent(queue, cmd)
	}
}
