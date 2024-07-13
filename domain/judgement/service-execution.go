package judgement

type ExecuteResult struct {
	JudgeStatus      JudgeStatusType
	NumberFinishedAt int
	MemoryUsed       int
	TimeUsed         int
	CPUTimeUsed      int
	FailedReason     string
}

type ExecutionService interface {
	execute(ExecutionCmd) (ExecuteResult, error)
}
