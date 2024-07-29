package svc

import (
	"moj/apps/question/db"
	ques_pb "moj/apps/question/rpc"
	"moj/domain/question"
)

type Server struct {
	ques_pb.UnimplementedQuestionServiceServer
	createQuestionCmdHandler *question.CreateQuestionCmdHandler
	modifyQuestionCmdHandler *question.ModifyQuestionCmdHandler
	questionDao              db.QuestionDao
	questionRepository       question.QuestionRepository
}

func NewServer(
	createQuestionCmdHandler *question.CreateQuestionCmdHandler,
	modifyQuestionCmdHandler *question.ModifyQuestionCmdHandler,
	questionDao db.QuestionDao,
	questionRepository question.QuestionRepository,
) *Server {
	return &Server{
		createQuestionCmdHandler: createQuestionCmdHandler,
		modifyQuestionCmdHandler: modifyQuestionCmdHandler,
		questionDao:              questionDao,
		questionRepository:       questionRepository,
	}
}
