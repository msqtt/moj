package record

import (
	"errors"
	"moj/domain/pkg/common"
	"moj/domain/pkg/queue"
)

var ErrRecordNotFound = errors.New("record not found")

type Record struct {
	RecordID         int
	AccountID        int
	GameID           int
	QuestionID       int
	Language         string
	Code             string
	CodeHash         string
	JudgeStatus      string
	FailedReason     string
	NumberFinishedAt int
	TotalQuestion    int
	CreateTime       int64
	FinishTime       int64
	MemoryUsed       int
	TimeUsed         int
	CPUTimeUsed      int
}

func NewRecord(accountID, gameID, questionID int, lang, code string, time int64) *Record {
	codeHash := common.Sha1(code)
	return &Record{
		AccountID:  accountID,
		GameID:     gameID,
		QuestionID: questionID,
		Language:   lang,
		Code:       code,
		CodeHash:   codeHash,
		CreateTime: time,
	}
}

func (r *Record) submit(queue queue.EventQueue) error {
	event := SubmitRecordEvent{
		RecordID:   r.RecordID,
		AccountID:  r.AccountID,
		QuestionID: r.QuestionID,
		GameID:     r.GameID,
		Language:   r.Language,
		Code:       r.Code,
		CodeHash:   r.CodeHash,
		CreateTime: r.CreateTime,
	}
	return queue.EnQueue(event)
}

func (r *Record) modify(queue queue.EventQueue, cmd ModifyRecordCmd) error {
	r.JudgeStatus = cmd.JudgeStatus
	r.FailedReason = cmd.FailedReason
	r.NumberFinishedAt = cmd.NumberFinishAt
	r.TotalQuestion = cmd.TotalQuestion
	r.FinishTime = cmd.Time
	r.MemoryUsed = cmd.MemoryUsed
	r.TimeUsed = cmd.TimeUsed
	r.CPUTimeUsed = cmd.CPUTimeUsed

	event := ModifyRecordEvent{
		RecordID:         r.RecordID,
		AccountID:        r.AccountID,
		QuestionID:       r.QuestionID,
		GameID:           r.GameID,
		JudgeStatus:      r.JudgeStatus,
		NumberFinishedAt: r.NumberFinishedAt,
		TotalQuestion:    r.TotalQuestion,
		FinishTime:       r.FinishTime,
	}
	return queue.EnQueue(event)
}
