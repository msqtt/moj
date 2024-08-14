package policy

import (
	"context"
	"errors"
	"moj/domain/judgement"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
	"moj/domain/question"
	"moj/domain/record"
	"time"
)

var (
	ErrFailedToExecuteCode error = errors.New("failed to execute code")
)

type JudgeOnSubmitPolicy struct {
	caseFileReader      CaseFileService
	executionCmdHandler *judgement.ExecutionCmdHandler
	questionRepository  question.QuestionRepository
	queue               queue.EventQueue
}

func NewJudgeOnSubmitPolicy(caseFileReader CaseFileService,
	executionCmdHandler *judgement.ExecutionCmdHandler,
	questionRepository question.QuestionRepository,
	queue queue.EventQueue,
) *JudgeOnSubmitPolicy {
	return &JudgeOnSubmitPolicy{
		caseFileReader:      caseFileReader,
		executionCmdHandler: executionCmdHandler,
		questionRepository:  questionRepository,
		queue:               queue,
	}
}

func (p *JudgeOnSubmitPolicy) OnEvent(event any) error {
	ctx := context.Background()
	evt, ok := event.(record.SubmitRecordEvent)
	if !ok {
		return domain_err.ErrEventTypeInvalid
	}
	que, err := p.questionRepository.FindQuestionByID(ctx, evt.QuestionID)
	if err != nil {
		return err
	}
	cases, err := p.caseFileReader.ReadAllCaseFile(ctx, que.Cases)
	if err != nil {
		return err
	}

	cmd := judgement.ExecutionCmd{
		RecordID:           evt.RecordID,
		QuestionID:         evt.QuestionID,
		QuestionModifyTime: que.ModifyTime,
		Cases:              cases,
		Language:           evt.Language,
		Code:               evt.Code,
		CodeHash:           evt.CodeHash,
		TimeLimit:          int64(que.TimeLimit),
		MemoryLimit:        int64(que.MemoryLimit),
		Time:               time.Now().Unix(),
	}

	err = p.executionCmdHandler.Handle(ctx, p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToExecuteCode, err)
	}
	return nil
}
