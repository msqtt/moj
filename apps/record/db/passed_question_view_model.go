package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PassedQuestionViewModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	AccountID  string             `bson:"account_id"`
	QuestionID string             `bson:"question_id"`
	Level      string
	RecordID   string    `bson:"record_id"`
	GameID     string    `bson:"game_id"`
	FinishTime time.Time `bson:"finish_time"`
}
