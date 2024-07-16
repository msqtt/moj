package record

type ModifyRecordEvent struct {
	RecordID         int
	AccountID        int
	QuestionID       int
	GameID           int
	JudgeStatus      string
	NumberFinishedAt int
	TotalQuestion    int
	FinishTime       int64
}
