package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordViewModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	GameID          string             `bson:"game_id"`
	Language        string
	CodeHash        string    `bson:"code_hash"`
	JudgeStatus     string    `bson:"judge_status"`
	NumberFinisheAt int       `bson:"number_finish_at"`
	TotalQuestion   int       `bson:"total_question"`
	CreateTime      time.Time `bson:"create_time"`
	MemoryUsed      int       `bson:"memory_used"`
	TimeUsed        int       `bson:"time_used"`
}
