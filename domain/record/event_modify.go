package record

type ModifyRecordEvent struct {
	RecordID         string
	AccountID        string
	QuestionID       string
	GameID           string
	JudgeStatus      string
	NumberFinishedAt int
	TotalQuestion    int
	FinishTime       int64
}
