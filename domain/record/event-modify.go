package record

type ModifyRecordEvent struct {
	RecordID       int
	AccountID      int
	QuestionID     int
	GameID         int
	JudgeStatus    string
	NumberFinishAt int
	TotalQuestion  int
	FinishTime     int64
}
