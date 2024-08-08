package judgement

import (
	"errors"

	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
)

type JudgeStatusType string

const (
	JudgeStatusAC  JudgeStatusType = "Accepted"
	JudgeStatusCE  JudgeStatusType = "Compile Error"
	JudgeStatusWA  JudgeStatusType = "Wrong Answer"
	JudgeStatusRE  JudgeStatusType = "Runtime Error"
	JudgeStatusTLE JudgeStatusType = "Time Limit Exceeded"
	JudgeStatusMLE JudgeStatusType = "Memory Limit Exceeded"
	JudgeStatusUE  JudgeStatusType = "Unknown Error"
)

type Judgement struct {
	JudgementID      string
	JudgeStatus      JudgeStatusType
	RecordID         string
	QuestionID       string
	Language         string
	Code             string
	CodeHash         string
	NumberFinishedAt int
	TotalQuestion    int
	MemoryUsed       int
	TimeUsed         int
	CPUTimeUsed      int
	ExecuteTime      int64
	FailedReason     string
}

var (
	ErrJudgementNotFound = errors.New("judgement not found")
	ErrEmptyCase         = errors.Join(domain_err.ErrInValided, errors.New("empty case"))
)

func NewJudgement(recordID, questionID string, total int,
	lang, code, codeHash string, time int64) *Judgement {
	return &Judgement{
		RecordID:      recordID,
		QuestionID:    questionID,
		Language:      lang,
		Code:          code,
		CodeHash:      codeHash,
		TotalQuestion: total,
		ExecuteTime:   time,
	}
}

func (j *Judgement) execute(queue queue.EventQueue,
	exeService ExecutionService, cmd ExecutionCmd) error {
	if j.TotalQuestion == 0 {
		return ErrEmptyCase
	}

	// start to execute code
	result, err := exeService.Execute(cmd)
	if err != nil {
		return err
	}

	// update info
	j.JudgeStatus = result.JudgeStatus
	j.NumberFinishedAt = result.NumberFinishedAt
	j.MemoryUsed = result.MemoryUsed
	j.TimeUsed = result.TimeUsed
	j.CPUTimeUsed = result.CPUTimeUsed
	j.FailedReason = result.FailedReason

	return j.sendExecutionEvent(queue, cmd)
}

func (j *Judgement) sendExecutionEvent(queue queue.EventQueue, cmd ExecutionCmd) error {
	event := ExecutionFinishEvent{
		RecordID:         cmd.RecordID,
		CodeHash:         j.CodeHash,
		NumberFinishedAt: j.NumberFinishedAt,
		TotalQuestion:    j.TotalQuestion,
		MemoryUsed:       j.MemoryUsed,
		TimeUsed:         j.TimeUsed,
		CPUTimeUsed:      j.CPUTimeUsed,
		FailedReason:     j.FailedReason,
	}

	return queue.EnQueue(event)
}
