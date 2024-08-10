package policy

import (
	"context"
	"errors"
	"moj/domain/judgement"
	"moj/domain/pkg/queue"
	"moj/domain/record"
	"time"
)

var ErrFailedToModifyRecord error = errors.New("failed to modify record")

type ModifyRecordPolicy struct {
	modifyRecordCmdHandler *record.ModifyRecordCmdHandler
	queue                  queue.EventQueue
}

func NewModifyRecordAfterExecutionPolicy(modifyRecordCmdHandler *record.ModifyRecordCmdHandler,
	queue queue.EventQueue,
) *ModifyRecordPolicy {
	return &ModifyRecordPolicy{
		modifyRecordCmdHandler: modifyRecordCmdHandler,
		queue:                  queue,
	}
}

func (p *ModifyRecordPolicy) OnEvent(event any) error {
	ctx := context.Background()
	evt, ok := event.(judgement.ExecutionFinishEvent)
	if !ok {
		return errors.New("invalid event type")
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

	_, err := p.modifyRecordCmdHandler.Handle(ctx, p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToModifyRecord, err)
	}
	return nil
}
