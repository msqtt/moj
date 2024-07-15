package policy

import (
	"errors"
	"moj/domain/judgement"
	"moj/domain/pkg/queue"
	"moj/domain/record"
	"time"
)

var ErrFailedToModifyRecord error = errors.New("failed to modify record")

type ModifyRecordAfterExecutionPolicy struct {
	modifyRecordCmdHandler record.ModifyRecordCmdHandler
	queue                  queue.EventQueue
}

func NewModifyRecordAfterExecutionPolicy(modifyRecordCmdHandler record.ModifyRecordCmdHandler,
	queue queue.EventQueue,
) ModifyRecordAfterExecutionPolicy {
	return ModifyRecordAfterExecutionPolicy{
		modifyRecordCmdHandler: modifyRecordCmdHandler,
		queue:                  queue,
	}
}

func (p *ModifyRecordAfterExecutionPolicy) OnEvent(event any) error {
	evt, ok := event.(judgement.ExecutionFinishEvent)
	if !ok {
		return nil
	}

	cmd := record.ModifyRecordCmd{
		RecordID:       evt.RecordID,
		JudgeStatus:    evt.JudgeStatus,
		FailedReason:   evt.FailedReason,
		MemoryUsed:     evt.MemoryUsed,
		TimeUsed:       evt.TimeUsed,
		CPUTimeUsed:    evt.CPUTimeUsed,
		NumberFinishAt: evt.NumberFinishedAt,
		TotalQuestion:  evt.TotalQuestion,
		Time:           time.Now().Unix(),
	}

	err := p.modifyRecordCmdHandler.Handle(p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToModifyRecord, err)
	}
	return nil
}
