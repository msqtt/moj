package svc

import (
	"context"
	"log/slog"
	"moj/apps/question/db"
	ques_pb "moj/apps/question/rpc"
	"moj/domain/question"
	"time"
)

// GetQuestionInfo implements ques_pb.QuestionServiceServer.
func (s *Server) GetQuestionInfo(ctx context.Context, req *ques_pb.GetQuestionInfoRequest) (
	resp *ques_pb.GetQuestionInfoResponse, err error) {
	slog.Debug("get question info", "req", req)
	ques, err := s.questionRepository.FindQuestionByID(req.QuestionID)

	if err != nil {
		slog.Error("get question info error", "error", err)
		err = responseStatusError(err)
		return
	}

	langs := db.FromQuestionLangs(ques.AllowedLanguages)
	cases := FromQuestionCases(ques.Cases)

	resp = &ques_pb.GetQuestionInfoResponse{
		QuestionInfo: &ques_pb.QuestionInfo{
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

	ques, err := s.questionDao.FindQuestionPage(req.GetCursor(), int(req.PageSize), m)
	if err != nil {
		slog.Error("get question page error", "error", err)
		err = responseStatusError(err)
		return
	}

	var nextCursor string
	quesView := make([]*ques_pb.QuestionView, len(ques))
	for i, q := range ques {
		quesView[i] = &ques_pb.QuestionView{
			QuestionID:      q.QuestionID,
			AccountID:       q.AccountID,
			Title:           q.Title,
			Enabled:         q.Enabled,
			Level:           ques_pb.QuestionLevel(q.Level),
			Tags:            q.Tags,
			TotalCaseNumber: int64(q.TotalCaseNumber),
			CreateTime:      q.CreateTime.Unix(),
		}
	}
	if len(quesView) > 0 {
		nextCursor = quesView[len(quesView)-1].QuestionID
	}

	resp = &ques_pb.GetQuestionPageResponse{
		QuestionView: quesView,
		NextCursor:   nextCursor,
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

	_, err = s.modifyQuestionCmdHandler.Handle(cmd)
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
	id, err := s.createQuestionCmdHandler.Handle(cmd)
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
			OutputFilePath: c.OutFilePath,
		}
	}
	return ret
}

func FromQuestionCases(cases []question.Case) []*ques_pb.Case {
	ret := make([]*ques_pb.Case, len(cases))
	for i, c := range cases {
		ret[i] = &ques_pb.Case{
			Number:        int64(c.Number),
			InputFilePath: c.InputFilePath,
			OutFilePath:   c.OutputFilePath,
		}
	}
	return ret
}

var _ ques_pb.QuestionServiceServer = (*Server)(nil)