package db

import (
	"moj/domain/record"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	AccountID       string             `bson:"account_id"`
	GameID          string             `bson:"game_id"`
	QuestionID      string             `bson:"question_id"`
	Language        string
	Code            string
	CodeHash        string    `bson:"code_hash"`
	JudgeStatus     string    `bson:"judge_status"`
	FailedReason    string    `bson:"failed_reason"`
	NumberFinisheAt int       `bson:"number_finish_at"`
	TotalQuestion   int       `bson:"total_question"`
	CreateTime      time.Time `bson:"create_time"`
	FinishTime      time.Time `bson:"finish_time"`
	MemoryUsed      int       `bson:"memory_used"`
	TimeUsed        int       `bson:"time_used"`
	CPUTimeUsed     int       `bson:"cpu_time_used"`
}

func NewRecordFromAggregation(r *record.Record) *RecordModel {
	id, _ := primitive.ObjectIDFromHex(r.RecordID)
	return &RecordModel{
		ID:              id,
		AccountID:       r.AccountID,
		GameID:          r.GameID,
		QuestionID:      r.QuestionID,
		Language:        r.Language,
		Code:            r.Code,
		CodeHash:        r.CodeHash,
		JudgeStatus:     r.JudgeStatus,
		FailedReason:    r.FailedReason,
		NumberFinisheAt: r.NumberFinishedAt,
		TotalQuestion:   r.TotalQuestion,
		CreateTime:      time.Unix(r.CreateTime, 0),
		FinishTime:      time.Unix(r.FinishTime, 0),
		MemoryUsed:      r.MemoryUsed,
		TimeUsed:        r.TimeUsed,
		CPUTimeUsed:     r.CPUTimeUsed,
	}
}

func (r *RecordModel) ToAggregation() *record.Record {
	return &record.Record{
		RecordID:         r.ID.Hex(),
		AccountID:        r.AccountID,
		GameID:           r.GameID,
		QuestionID:       r.QuestionID,
		Language:         r.Language,
		Code:             r.Code,
		CodeHash:         r.CodeHash,
		JudgeStatus:      r.JudgeStatus,
		FailedReason:     r.FailedReason,
		NumberFinishedAt: r.NumberFinisheAt,
		TotalQuestion:    r.TotalQuestion,
		CreateTime:       r.CreateTime.Unix(),
		FinishTime:       r.FinishTime.Unix(),
		MemoryUsed:       r.MemoryUsed,
		TimeUsed:         r.TimeUsed,
		CPUTimeUsed:      r.CPUTimeUsed,
	}
}
