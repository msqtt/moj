package record

import "moj/domain/pkg/queue"

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

func (h *ModifyRecordCmdHandler) Handle(queue queue.EventQueue, cmd ModifyRecordCmd) (any, error) {
	rec, err := h.repo.FindRecordByID(cmd.RecordID)
	if err != nil {
		return nil, err
	}
	rec.modify(queue, cmd)

	return h.repo.Save(rec)
}
