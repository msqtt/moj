package svc

import (
	"context"
	"log/slog"
	"moj/question/db"
	ques_pb "moj/question/rpc"
	"moj/domain/question"
	"time"
)

// GetQuestion implements ques_pb.QuestionServiceServer.
func (s *Server) GetQuestion(ctx context.Context, req *ques_pb.GetQuestionRequest) (
	resp *ques_pb.GetQuestionResponse, err error) {
	slog.Debug("get question info", "req", req)
	ques, err := s.questionRepository.FindQuestionByID(ctx, req.QuestionID)

	if err != nil {
		slog.Error("get question info error", "error", err)
		err = responseStatusError(err)
		return
	}

	langs := db.FromQuestionLangs(ques.AllowedLanguages)
	cases := FromQuestionCases(ques.Cases)

	resp = &ques_pb.GetQuestionResponse{
		Question: &ques_pb.Question{
			QuestionID:       ques.QuestionID,
			AccountID:        ques.AccountID,
			Enabled:          ques.Enabled,
			Title:            ques.Title,
			Content:          ques.Content,
			Level:            ques_pb.QuestionLevel(ques.Level),
			AllowedLanguages: langs,
			Cases:            cases,
			TimeLimit:        int64(ques.TimeLimit),
			MemoryLimit:      int64(ques.MemoryLimit),
			Tags:             ques.Tags,
			CreateTime:       ques.CreateTime,
			ModifyTime:       ques.ModifyTime,
		},
	}
	return
}

// GetQuestionPage implements ques_pb.QuestionServiceServer.
func (s *Server) GetQuestionPage(ctx context.Context, req *ques_pb.GetQuestionPageRequest) (
	resp *ques_pb.GetQuestionPageResponse, err error) {
	slog.Debug("get question page", "req", req)

	m := make(map[string]any)
	if req.FilterOptions != nil {
		if req.FilterOptions.Word != nil {
			m["word"] = req.FilterOptions.Word
		}
		if req.FilterOptions.Enabled != nil {
			m["enabled"] = req.FilterOptions.Enabled
		}
		if req.FilterOptions.Level != nil {
			m["level"] = req.FilterOptions.Level
		}
		if req.FilterOptions.Language != nil {
			m["language"] = req.FilterOptions.Language
		}
		if req.FilterOptions.AccountID != nil {
			m["account_id"] = req.FilterOptions.AccountID
		}
	}

	ques, err := s.questionDao.FindQuestionPage(ctx, req.GetCursor(), int(req.PageSize), m)
	if err != nil {
		slog.Error("get question page error", "error", err)
		err = responseStatusError(err)
		return
	}

	var nextCursor string
	quesView := make([]*ques_pb.Question, len(ques))
	for i, q := range ques {
		agg := q.ToAggregate()
		quesView[i] = &ques_pb.Question{
			QuestionID:       q.ID.Hex(),
			AccountID:        q.AccountID,
			Enabled:          q.Enabled,
			Title:            q.Title,
			Content:          q.Content,
			Level:            ques_pb.QuestionLevel(agg.Level),
			AllowedLanguages: q.AllowedLanguages,
			Cases:            FromModelCases(q.Cases),
			TimeLimit:        int64(q.TimeLimit),
			MemoryLimit:      int64(q.MemoryLimit),
			Tags:             q.Tags,
			CreateTime:       q.CreateTime.Unix(),
			ModifyTime:       q.ModifyTime.Unix(),
		}
	}
	if len(quesView) > 0 {
		nextCursor = quesView[len(quesView)-1].QuestionID
	}

	resp = &ques_pb.GetQuestionPageResponse{
		Questions:  quesView,
		NextCursor: nextCursor,
	}
	return
}

// UpdateQuestion implements ques_pb.QuestionServiceServer.
func (s *Server) UpdateQuestion(ctx context.Context, req *ques_pb.UpdateQuestionRequest) (
	resp *ques_pb.UpdateQuestionResponse, err error) {
	slog.Debug("update question", "req", req)

	cmd := question.ModifyQuestionCmd{
		QuestionID:       req.QuestionID,
		Enabled:          req.Enabled,
		Title:            req.Title,
		Content:          req.Content,
		Level:            question.QuestionLevel(req.Level),
		AllowedLanguages: db.ToQuestionLangs(req.GetAllowedLanguages()),
		TimeLimit:        int(req.TimeLimit),
		MemoryLimit:      int(req.MemoryLimit),
		Tags:             req.Tags,
		Time:             time.Now().Unix(),
		Cases:            ToQuestionCases(req.GetCases()),
	}

	slog.Info("update question command execution", "cmd", cmd)

	_, err = s.modifyQuestionCmdHandler.Handle(ctx, cmd)
	if err != nil {
		slog.Error("update question command execution error", "error", err)
		err = responseStatusError(err)
		return
	}
	resp = &ques_pb.UpdateQuestionResponse{
		Time: time.Now().Unix(),
	}
	return
}

// UploadQuestion implements ques_pb.QuestionServiceServer.
func (s *Server) UploadQuestion(ctx context.Context, req *ques_pb.UploadQuestionRequest) (
	resp *ques_pb.UploadQuestionResponse, err error) {
	slog.Debug("upload question", "req", req)

	langs := db.ToQuestionLangs(req.GetAllowedLanguages())
	cases := ToQuestionCases(req.GetCases())
	cmd := question.CreateQuestionCmd{
		AccountID:        req.AccountID,
		Title:            req.Title,
		Content:          req.Content,
		Level:            question.QuestionLevel(req.Level),
		AllowedLanguages: langs,
		TimeLimit:        int(req.TimeLimit),
		MemoryLimit:      int(req.MemoryLimit),
		Tags:             req.Tags,
		Time:             time.Now().Unix(),
		Cases:            cases,
	}

	slog.Info("upload question command execution", "cmd", cmd)
	id, err := s.createQuestionCmdHandler.Handle(ctx, cmd)
	if err != nil {
		slog.Error("upload question command execution error", "error", err)
		err = responseStatusError(err)
		return
	}

	resp = &ques_pb.UploadQuestionResponse{
		QuestionID: id.(string),
		Time:       time.Now().Unix(),
	}
	return
}

func ToQuestionCases(cases []*ques_pb.Case) []question.Case {
	ret := make([]question.Case, len(cases))
	for i, c := range cases {
		ret[i] = question.Case{
			Number:         int(c.GetNumber()),
			InputFilePath:  c.InputFilePath,
			OutputFilePath: c.OutputFilePath,
		}
	}
	return ret
}

func FromQuestionCases(cases []question.Case) []*ques_pb.Case {
	ret := make([]*ques_pb.Case, len(cases))
	for i, c := range cases {
		ret[i] = &ques_pb.Case{
			Number:         int64(c.Number),
			InputFilePath:  c.InputFilePath,
			OutputFilePath: c.OutputFilePath,
		}
	}
	return ret
}

func FromModelCases(cases []db.Case) []*ques_pb.Case {
	ret := make([]*ques_pb.Case, len(cases))
	for i, c := range cases {
		ret[i] = &ques_pb.Case{
			Number:         int64(c.Number),
			InputFilePath:  c.InputFilePath,
			OutputFilePath: c.OutputFilePath,
		}
	}
	return ret
}

// DeleteQuestion implements ques_pb.QuestionServiceServer.
func (s *Server) DeleteQuestion(ctx context.Context, req *ques_pb.DeleteQuestionRequest) (
	resp *ques_pb.DeleteQuestionResponse, err error) {
	slog.Debug("delete question", "req", req)

	err = s.questionDao.DeleteQuestion(ctx, req.QuestionID)
	if err != nil {
		slog.Error("delete question error", "error", err)
		return nil, err
	}

	return &ques_pb.DeleteQuestionResponse{
		Time: time.Now().Unix(),
	}, nil
}

var _ ques_pb.QuestionServiceServer = (*Server)(nil)
