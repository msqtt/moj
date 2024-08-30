package svc

import (
	"moj/domain/record"
	"moj/record/db"
	"moj/record/domain"
	red_pb "moj/record/rpc"
)

type Server struct {
	red_pb.UnimplementedRecordServiceServer
	commandInvoker         domain.CommandInvoker
	modifyRecordCmdHandler *record.ModifyRecordCmdHandler
	submitRecordCmdHandler *record.SubmitRecordCmdHandler
	recordRepository       record.RecordRepository
	recordViewDao          db.RecordViewDao
	dayTaskViewDao         db.DailyTaskViewDao
	passQuestionViewDao    db.PassQuestionViewDao
}

func NewServer(
	commandInvoker domain.CommandInvoker,
	modifyRecordCmdHandler *record.ModifyRecordCmdHandler,
	submitRecordCmdHandler *record.SubmitRecordCmdHandler,
	recordRepository record.RecordRepository,
	recordViewDao db.RecordViewDao,
	dayTaskViewDao db.DailyTaskViewDao,
	passedQuestionViewDao db.PassQuestionViewDao,
) *Server {
	return &Server{
		commandInvoker:         commandInvoker,
		modifyRecordCmdHandler: modifyRecordCmdHandler,
		submitRecordCmdHandler: submitRecordCmdHandler,
		recordRepository:       recordRepository,
		recordViewDao:          recordViewDao,
		dayTaskViewDao:         dayTaskViewDao,
		passQuestionViewDao:    passedQuestionViewDao,
	}
}
