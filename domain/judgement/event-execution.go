package judgement

type ExecutionFinishEvent struct {
	RecordID         int
	CodeHash         string
	NumberFinishedAt int
	TotalQuestion    int
	MemoryUsed       int
	TimeUsed         int
	CPUTimeUsed      int
	FailedReason     string
}
