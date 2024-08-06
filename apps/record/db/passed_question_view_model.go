package db

import "time"

type PassedQuestionViewModel struct {
	AccountID  string    `bson:"account_id"`
	QuestionID string    `bson:"question_id"`
	RecordID   string    `bson:"record_id"`
	FinishTime time.Time `bson:"finish_time"`
}
