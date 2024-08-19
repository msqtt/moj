package svc

import (
	"context"
	"log/slog"
	"moj/judgement/domain"
	jud_pb "moj/judgement/rpc"
	"moj/domain/judgement"
	"moj/domain/pkg/queue"
	"moj/domain/policy"
	"moj/domain/question"
	"time"
)

type Server struct {
	jud_pb.UnimplementedJudgeServiceServer
	cmdInvoker          domain.CommandInvoker
	executionCmdHandler *judgement.ExecutionCmdHandler
	questionRepo        question.QuestionRepository
	reader              policy.CaseFileService
}

// ExecuteJudge implements jud_pb.JudgeServiceServer.
func (s *Server) ExecuteJudge(ctx context.Context, req *jud_pb.ExecuteJudgeRequest) (
	resp *jud_pb.ExecuteJudgeResponse, err error) {

	slog.Debug("execute judge", "req", req)

	ques, err := s.questionRepo.FindQuestionByID(ctx, req.QuestionID)
	if err != nil {
		slog.Error("failed to find question", "error", err)
		err = responseStatusError(err)
		return nil, err
	}

	cases, err := s.reader.ReadAllCaseFile(ctx, ques.Cases)
	if err != nil {
		slog.Error("failed to read question case file", "error", err)
		err = responseStatusError(err)
		return nil, err
	}
	cmd := judgement.ExecutionCmd{
		RecordID:           req.RecordID,
		QuestionID:         req.QuestionID,
		QuestionModifyTime: ques.ModifyTime,
		Cases:              cases,
		Language:           req.Language,
		Code:               req.Code,
		CodeHash:           req.CodeHash,
		TimeLimit:          int64(ques.TimeLimit),
		MemoryLimit:        int64(ques.MemoryLimit),
		Time:               time.Now().Unix(),
	}

	slog.Debug("start to invoke execute judgement command", "cmd", cmd)
	err = s.cmdInvoker.Invoke(ctx, func(ctx context.Context, eq queue.EventQueue) (any, error) {
		return nil, s.executionCmdHandler.Handle(ctx, eq, cmd)
	})

	if err != nil {
		slog.Error("failed to invoke execute judgement command", "error", err)
		err = responseStatusError(err)
	}

	resp = &jud_pb.ExecuteJudgeResponse{
		Time: time.Now().Unix(),
	}
	return
}

func NewServer(cmdInvoker domain.CommandInvoker, executionCmdHandler *judgement.ExecutionCmdHandler,
	questionRepo question.QuestionRepository,
	reader policy.CaseFileService,
) *Server {
	return &Server{
		cmdInvoker:          cmdInvoker,
		executionCmdHandler: executionCmdHandler,
		questionRepo:        questionRepo,
		reader:              reader,
	}
}
