package record

type SubmitRecordEvent struct {
	RecordID   string
	AccountID  string
	QuestionID string
	GameID     string
	Language   string
	Code       string
	CodeHash   string
	CreateTime int64
}
