package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionViewModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	QuestionID      string             `bson:"question_id"`
	AccountID       string             `bson:"account_id"`
	Title           string
	Enabled         bool
	Level           int
	Tags            []string
	TotalCaseNumber int       `bson:"total_case_number"`
	CreateTime      time.Time `bson:"create_time"`
}
