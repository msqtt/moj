package record

type SubmitRecordEvent struct {
	RecordID   int
	AccountID  int
	QuestionID int
	GameID     int
	Language   string
	Code       string
	CodeHash   string
	CreateTime int64
}
