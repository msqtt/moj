package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PassStatus = string

const (
	PassStatusPass    PassStatus = "pass"
	PassStatusWorking PassStatus = "working"
	PassStatusUndo    PassStatus = "undo"
)

type PassedQuestionViewModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	AccountID  string             `bson:"account_id"`
	QuestionID string             `bson:"question_id"`
	Status     PassStatus
	Level      string
	RecordID   string    `bson:"record_id"`
	GameID     string    `bson:"game_id"`
	FinishTime time.Time `bson:"finish_time"`
}
