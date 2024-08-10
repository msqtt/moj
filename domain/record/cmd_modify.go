package record

import (
	"context"
	"moj/domain/pkg/queue"
)

type ModifyRecordCmd struct {
	RecordID       string
	JudgeStatus    string
	FailedReason   string
	MemoryUsed     int
	TimeUsed       int
	CPUTimeUsed    int
	NumberFinishAt int
	TotalQuestion  int
	Time           int64
}

type ModifyRecordCmdHandler struct {
	repo RecordRepository
}

func NewModifyRecordCmdHandler(repo RecordRepository) *ModifyRecordCmdHandler {
	return &ModifyRecordCmdHandler{
		repo: repo,
	}
}

func (h *ModifyRecordCmdHandler) Handle(ctx context.Context, queue queue.EventQueue, cmd ModifyRecordCmd) (any, error) {
	rec, err := h.repo.FindRecordByID(ctx, cmd.RecordID)
	if err != nil {
		return nil, err
	}
	rec.modify(queue, cmd)

	return h.repo.Save(ctx, rec)
}
