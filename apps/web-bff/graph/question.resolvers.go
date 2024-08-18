package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"log/slog"
	ques_pb "moj/apps/question/rpc"
	red_pb "moj/apps/record/rpc"
	"moj/apps/web-bff/graph/model"
	"moj/apps/web-bff/pkg"

	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// CreateQuestion is the resolver for the createQuestion field.
func (r *mutationResolver) CreateQuestion(ctx context.Context, input model.QuestionInput) (*model.Question, error) {
	uid, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	cases := make([]*ques_pb.Case, len(input.Cases))

	for id, ca := range input.Cases {
		cases[id] = toPbCase(ca)
	}

	req := &ques_pb.UploadQuestionRequest{
		Title:            input.Title,
		AccountID:        uid,
		Content:          input.Content,
		Level:            *toPbLevel(input.Level),
		AllowedLanguages: input.AllowedLanguages,
		TimeLimit:        int64(input.TimeLimit),
		MemoryLimit:      int64(input.MemoryLimit),
		Tags:             input.Tags,
		Cases:            cases,
	}

	resp, err := r.RpcClients.QuestionClient.UploadQuestion(ctx, req)
	if err != nil {
		slog.Error("failed to upload question", "error", err)

		return nil, ErrInternal
	}

	ques, err := findQuestion(r.RpcClients.QuestionClient, ctx, resp.QuestionID)
	if err != nil {
		slog.Error("upload question: failed to get question", "error", err)
		return nil, ErrInternal
	}

	return ques, err
}

// ModifyQuestion is the resolver for the modifyQuestion field.
func (r *mutationResolver) ModifyQuestion(ctx context.Context, input model.QuestionInput) (*model.Question, error) {
	if input.ID == nil {
		return nil, ErrQuestionIDNotProvide
	}

	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	cases := make([]*ques_pb.Case, len(input.Cases))
	for id, c := range input.Cases {
		cases[id] = toPbCase(c)
	}
	req := &ques_pb.UpdateQuestionRequest{
		QuestionID:       *input.ID,
		Title:            input.Title,
		Content:          input.Content,
		Level:            *toPbLevel(input.Level),
		Enabled:          input.Enabled,
		AllowedLanguages: input.AllowedLanguages,
		TimeLimit:        int64(input.TimeLimit),
		MemoryLimit:      int64(input.MemoryLimit),
		Tags:             input.Tags,
		Cases:            cases,
	}

	_, err = r.RpcClients.QuestionClient.UpdateQuestion(ctx, req)
	if err != nil {
		slog.Error("failed to update question", "error", err)
		return nil, err
	}

	ques, err := findQuestion(r.RpcClients.QuestionClient, ctx, *input.ID)
	if err != nil {
		slog.Error("update question: failed to get question", "error", err)
		return nil, err
	}
	return ques, err
}

// DeleteQuestion is the resolver for the deleteQuestion field.
func (r *mutationResolver) DeleteQuestion(ctx context.Context, id string) (*model.Time, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	req := &ques_pb.DeleteQuestionRequest{
		QuestionID: id,
	}

	resp, err := r.RpcClients.QuestionClient.DeleteQuestion(ctx, req)
	if err != nil {
		slog.Error("failed to delete question", "error", err)
		return nil, err
	}

	return fromInt64Second(resp.GetTime()), nil
}

// CheckUserPassedQuestion is the resolver for the checkUserPassedQuestion field.
func (r *mutationResolver) CheckUserPassedQuestion(ctx context.Context, id string) (bool, error) {
	uid, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", false)
	if err != nil {
		return false, err
	}

	req := &red_pb.CheckAccountPassRequest{
		QuestionID: id,
		AccountID:  uid,
	}

	resp, err := r.RpcClients.RecordClient.CheckAccountPass(ctx, req)
	if err != nil {
		slog.Error("failed to check user passed question", "error", err)
		return false, ErrInternal
	}

	return resp.GetIsPass(), err
}

// Question is the resolver for the question field.
func (r *queryResolver) Question(ctx context.Context, id string) (*model.Question, error) {
	ques, err := findQuestion(r.RpcClients.QuestionClient, ctx, id)
	if err != nil {
		slog.Error("failed to get question", "error", err)
		return nil, err
	}

	if !ques.Enabled {
		_, err = checkUserLogin(r.RpcClients.UserClient, r.sessionManager,
			ctx, "", true)
		if err != nil {
			return nil, ErrUserUnAuthorized
		}
	}

	return ques, nil
}

// Questions is the resolver for the questions field.
func (r *queryResolver) Questions(ctx context.Context, pageSize int, afterID *string, filter *model.QuestionsFilter) (*model.QuestionPage, error) {
	var filterOptions *ques_pb.GetQuestionPageRequest_Option
	if filter == nil || filter.Enabled == nil || !*filter.Enabled {
		_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
		if err != nil {
			// block
			if filter == nil {
				filter = &model.QuestionsFilter{Enabled: new(bool)}
			}
			*filter.Enabled = true
		}
		// pass
	}

	if filter != nil {
		var level *ques_pb.QuestionLevel
		if filter.Level != nil {
			level = toPbLevel(*filter.Level)
		}
		filterOptions = &ques_pb.GetQuestionPageRequest_Option{
			Word:      filter.Word,
			Enabled:   filter.Enabled,
			Level:     level,
			Language:  filter.Language,
			AccountID: filter.CreaterID,
		}
	}

	req := &ques_pb.GetQuestionPageRequest{
		Cursor:        *afterID,
		PageSize:      int64(pageSize),
		FilterOptions: filterOptions,
	}

	resp, err := r.RpcClients.QuestionClient.GetQuestionPage(ctx, req)
	if err != nil {
		slog.Error("failed to get question page", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrQuestionNotFound
			}
		}
		return nil, ErrInternal
	}

	questions := make([]*model.Question, len(resp.Questions))
	for id, q := range resp.Questions {
		questions[id] = fromPbQuestion(q)
	}

	return &model.QuestionPage{
		NextID:    resp.NextCursor,
		Questions: questions,
	}, err
}

// QuestionSubmitCount is the resolver for the questionSubmitCount field.
func (r *queryResolver) QuestionSubmitCount(ctx context.Context, qid string, gid *string) (*model.SubmitCount, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	req := &red_pb.GetQuestionRecordCountRequest{
		QuestionID: qid,
		GameID:     gid,
	}
	resp, err := r.RpcClients.RecordClient.GetQuestionRecordCount(ctx, req)
	if err != nil {
		slog.Error("failed to get question submit count", "error", err)
		return nil, ErrInternal
	}

	return &model.SubmitCount{
		SubmitCount: int(resp.SubmitTotal),
		PassedCount: int(resp.PassedCount),
	}, err
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func findQuestion(client ques_pb.QuestionServiceClient, ctx context.Context, id string) (*model.Question, error) {
	req := &ques_pb.GetQuestionRequest{
		QuestionID: id,
	}
	resp, err := client.GetQuestion(ctx, req)
	if err != nil {
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrQuestionNotFound
			}
		}
		return nil, ErrInternal
	}
	return fromPbQuestion(resp.Question), nil
}

var (
	ErrQuestionNotFound     = errors.New("题目未找到")
	ErrQuestionIDNotProvide = errors.New("题目ID未提供")
)

func fromPbCase(cas *ques_pb.Case) *model.Case {
	return &model.Case{
		Number:         int(cas.Number),
		InputFilePath:  cas.InputFilePath,
		OutputFilePath: cas.OutputFilePath,
	}
}
func fromPbQuestion(q *ques_pb.Question) *model.Question {
	cases := make([]*model.Case, len(q.Cases))
	for id, ca := range q.Cases {
		cases[id] = fromPbCase(ca)
	}
	return &model.Question{
		ID:               q.QuestionID,
		CreaterID:        q.AccountID,
		Enabled:          q.Enabled,
		Title:            q.Title,
		Content:          q.Content,
		Level:            model.Level(q.Level.String()),
		AllowedLanguages: q.AllowedLanguages,
		TimeLimit:        int(q.TimeLimit),
		MemoryLimit:      int(q.MemoryLimit),
		Tags:             q.Tags,
		CreateTime:       pkg.Int64ToString(q.CreateTime),
		ModifyTime:       pkg.Int64ToString(q.ModifyTime),
		Cases:            cases,
	}
}
func toPbCase(cas *model.CaseInput) *ques_pb.Case {
	return &ques_pb.Case{
		Number:         int64(cas.Number),
		InputFilePath:  cas.InputFilePath,
		OutputFilePath: cas.OutputFilePath,
	}
}
func toPbLevel(level model.Level) *ques_pb.QuestionLevel {
	m := map[string]int{"Eazy": 0, "Normal": 1, "Hard": 2}
	ret := ques_pb.QuestionLevel(m[level.String()])
	return &ret
}
