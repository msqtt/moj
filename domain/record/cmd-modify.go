package record

import "moj/domain/pkg/queue"

type ModifyRecordCmd struct {
	RecordID       int
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

func (h *ModifyRecordCmdHandler) Handle(queue queue.EventQueue, cmd ModifyRecordCmd) error {
	rec, err := h.repo.findRecordByID(cmd.RecordID)
	if err != nil {
		return err
	}

	rec.modify(queue, cmd)

	return h.repo.save(rec)
}
