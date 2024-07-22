package judgement

type ExecutionFinishEvent struct {
	RecordID         string
	CodeHash         string
	JudgeStatus      string
	NumberFinishedAt int
	TotalQuestion    int
	MemoryUsed       int
	TimeUsed         int
	CPUTimeUsed      int
	FailedReason     string
}
