package policy

import (
	"errors"
	"github.com/msqtt/moj/domain/judgement"
	"github.com/msqtt/moj/domain/pkg/queue"
	"github.com/msqtt/moj/domain/question"
	"github.com/msqtt/moj/domain/record"
	"time"
)

var (
	ErrFailedToExecuteCode error = errors.New("failed to execute code")
)

type JudgeOnSubmitPolicy struct {
	caseFileReader      CaseFileService
	executionCmdHandler judgement.ExecutionCmdHandler
	questionRepository  question.QuestionRepository
	queue               queue.EventQueue
}

func NewJudgeOnSubmitPolicy(caseFileReader CaseFileService,
	executionCmdHandler judgement.ExecutionCmdHandler,
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
	evt, ok := event.(record.SubmitRecordEvent)
	if !ok {
		return nil
	}
	que, err := p.questionRepository.FindQuestionByID(evt.QuestionID)
	if err != nil {
		return err
	}
	cases, err := p.caseFileReader.ReadAllCaseFile(que.Case)
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
		Time:               time.Now().Unix(),
	}

	err = p.executionCmdHandler.Handle(p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToExecuteCode, err)
	}
	return nil
}
