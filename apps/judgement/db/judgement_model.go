package db

import (
	"moj/domain/judgement"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JudgementModel struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	JudgeStatus      string             `bson:"judge_status"`
	RecordID         string             `bson:"record_id"`
	QuestionID       string             `bson:"question_id"`
	Code             string
	CodeHash         string `bson:"code_hash"`
	NumberFinishedAt int    `bson:"number_finished_at"`
	TotalQuestion    int    `bson:"total_question"`
	MemoryUsed       int    `bson:"memory_used"`
	TimeUsed         int    `bson:"time_used"`
	CpuTimeUsed      int    `bson:"cpu_time_used"`
	ExecuteTime      int    `bson:"execute_time"`
	FailedReason     string `bson:"failed_reason"`
}

func NewJudgementModelFromAggreation(ju *judgement.Judgement) *JudgementModel {
	id, _ := primitive.ObjectIDFromHex(ju.JudgementID)
	return &JudgementModel{
		ID:               id,
		JudgeStatus:      string(ju.JudgeStatus),
		RecordID:         ju.RecordID,
		QuestionID:       ju.QuestionID,
		Code:             ju.Code,
		CodeHash:         ju.CodeHash,
		NumberFinishedAt: ju.NumberFinishedAt,
		TotalQuestion:    ju.TotalQuestion,
		MemoryUsed:       ju.MemoryUsed,
		TimeUsed:         ju.TimeUsed,
		CpuTimeUsed:      ju.CPUTimeUsed,
		ExecuteTime:      int(ju.ExecuteTime),
		FailedReason:     ju.FailedReason,
	}
}

func (ju *JudgementModel) ToAggreation() *judgement.Judgement {
	return &judgement.Judgement{
		JudgementID:      ju.ID.Hex(),
		JudgeStatus:      judgement.JudgeStatusType(ju.JudgeStatus),
		RecordID:         ju.RecordID,
		QuestionID:       ju.QuestionID,
		Code:             ju.Code,
		CodeHash:         ju.CodeHash,
		NumberFinishedAt: ju.NumberFinishedAt,
		TotalQuestion:    ju.TotalQuestion,
		MemoryUsed:       ju.MemoryUsed,
		TimeUsed:         ju.TimeUsed,
		CPUTimeUsed:      ju.CpuTimeUsed,
		ExecuteTime:      int64(ju.ExecuteTime),
		FailedReason:     ju.FailedReason,
	}
}
