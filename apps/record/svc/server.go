package svc

import (
	"moj/apps/record/db"
	"moj/apps/record/domain"
	red_pb "moj/apps/record/rpc"
	"moj/domain/record"
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
