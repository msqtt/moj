package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"log/slog"
	game_pb "moj/apps/game/rpc"
	"moj/apps/web-bff/graph/model"
	"moj/apps/web-bff/pkg"
	"strconv"

	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// CreateGame is the resolver for the createGame field.
func (r *mutationResolver) CreateGame(ctx context.Context, input model.GameInput) (*model.Game, error) {
	uid, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	startTime, err := stringToInt64(input.StartTime)
	if err != nil {
		slog.Error("create game: failed to parse time", err)
		return nil, ErrInternal
	}
	endTime, err := stringToInt64(input.EndTime)
	if err != nil {
		slog.Error("create game: failed to parse time", err)
		return nil, ErrInternal
	}

	req := &game_pb.CreateGameRequest{
		Title:        input.Title,
		Desc:         input.Description,
		AccountID:    uid,
		StartTime:    startTime,
		EndTime:      endTime,
		QuestionList: toPbGameQuestionList(input.QuestionList),
	}

	resp, err := r.RpcClients.GameClient.CreateGame(ctx, req)
	if err != nil {
		slog.Error("failed to create game", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.InvalidArgument {
				return nil, ErrGameInvalid
			}
		}
		return nil, err
	}

	game, err := findGame(r.RpcClients.GameClient, ctx, resp.GameID)
	if err != nil {
		slog.Error("create game: failed to find game", "error", err)
		return nil, err
	}

	return game, err
}

// ModifyGame is the resolver for the updateGame field.
func (r *mutationResolver) ModifyGame(ctx context.Context, input model.GameInput) (*model.Game, error) {
	if input.ID == nil {
		return nil, ErrGameIDNotProvide
	}

	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	startTime, err := stringToInt64(input.StartTime)
	if err != nil {
		slog.Error("modify game: failed to parse time", err)
		return nil, ErrInternal
	}
	endTime, err := stringToInt64(input.EndTime)
	if err != nil {
		slog.Error("modify game: failed to parse time", err)
		return nil, ErrInternal
	}

	req := &game_pb.UpdateGameRequest{
		GameID:       *input.ID,
		Title:        input.Title,
		Desc:         input.Description,
		StartTime:    startTime,
		EndTime:      endTime,
		QuestionList: toPbGameQuestionList(input.QuestionList),
	}

	_, err = r.RpcClients.GameClient.UpdateGame(ctx, req)
	if err != nil {
		slog.Error("failed to update game", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			switch status.Code() {
			case codes.NotFound:
				return nil, ErrGameNotFound
			case codes.InvalidArgument:
				return nil, ErrGameInvalid
			}
		}
		return nil, ErrInternal
	}

	game, err := findGame(r.RpcClients.GameClient, ctx, *input.ID)
	if err != nil {
		slog.Error("update game: failed to find game", "error", err)
		return nil, err
	}

	return game, err
}

// DeleteGame is the resolver for the deleteGame field.
func (r *mutationResolver) DeleteGame(ctx context.Context, gid string) (*model.Time, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	resp, err := r.RpcClients.GameClient.DeleteGame(ctx, &game_pb.DeleteGameRequest{GameID: gid})
	if err != nil {
		slog.Error("failed to delete game", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrGameNotFound
			}
		}
		return nil, ErrInternal
	}

	return fromInt64Second(resp.GetTime()), err
}

// SignUpGame is the resolver for the signUpGame field.
func (r *mutationResolver) SignUpGame(ctx context.Context, gid string) (*model.Time, error) {
	uid, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", false)
	if err != nil {
		return nil, err
	}

	req := &game_pb.SignUpGameRequest{
		GameID:    gid,
		AccountID: uid,
	}

	resp, err := r.RpcClients.GameClient.SignUpGame(ctx, req)
	if err != nil {
		slog.Error("failed to sign up game", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			switch status.Code() {
			case codes.NotFound:
				return nil, ErrGameNotFound
			case codes.AlreadyExists:
				return nil, ErrGameAlreadySignUp
			case codes.DeadlineExceeded:
				return nil, ErrSignUpOperationExpired
			}
		}
		return nil, ErrInternal
	}

	return fromInt64Second(resp.GetTime()), err
}

// CancelSignUpGame is the resolver for the cancelSignUpGame field.
func (r *mutationResolver) CancelSignUpGame(ctx context.Context, gid string) (*model.Time, error) {
	uid, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", false)
	if err != nil {
		return nil, err
	}

	req := &game_pb.CancelSignUpGameRequest{
		GameID:    gid,
		AccountID: uid,
	}

	resp, err := r.RpcClients.GameClient.CancelSignUpGame(ctx, req)
	if err != nil {
		slog.Error("failed to cancel sign up game", "error", err)
		if status, ok := gstatus.FromError(err); ok {
			switch status.Code() {
			case codes.NotFound:
				return nil, ErrGameNotFound
			case codes.AlreadyExists:
				return nil, ErrGameNotSignUp
			case codes.DeadlineExceeded:
				return nil, ErrSignUpOperationExpired
			}
		}
		return nil, ErrInternal
	}

	return fromInt64Second(resp.GetTime()), err
}

// Game is the resolver for the game field.
func (r *queryResolver) Game(ctx context.Context, id string) (*model.Game, error) {
	game, err := findGame(r.RpcClients.GameClient, ctx, id)
	if err != nil {
		slog.Error("failed to find game", "error", err)
		return nil, err
	}
	return game, err
}

// Games is the resolver for the games field.
func (r *queryResolver) Games(ctx context.Context, pageSize int, afterID *string, filter *model.GamesFilter) (*model.GamePage, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", false)
	if err != nil {
		return nil, err
	}

	var cursor string
	if afterID == nil {
		cursor = ""
	} else {
		cursor = *afterID
	}

	var filterOptions *game_pb.GetGamePageRequest_Option
	if filter != nil {
		var t *int64
		if filter.Time != nil {
			tmp, err := stringToInt64(*filter.Time)
			if err != nil {
				slog.Error("get games: failed to parse time", err)
				return nil, ErrInternal
			}
			t = &tmp
		}
		filterOptions = &game_pb.GetGamePageRequest_Option{
			Word: filter.Word,
			Time: t,
		}
	}

	resp, err := r.RpcClients.GameClient.GetGamePage(ctx, &game_pb.GetGamePageRequest{
		Cursor:        cursor,
		PageSize:      int64(pageSize),
		FilterOptions: filterOptions,
	})
	if err != nil {
		slog.Error("get game page error", "err", err)
		return nil, ErrInternal
	}

	games := make([]*model.Game, len(resp.GetGames()))

	for id, g := range resp.GetGames() {
		games[id] = fromPbGame(g)
	}

	return &model.GamePage{
		NextID: resp.GetCursor(),
		Games:  games,
	}, err
}

// GameScore is the resolver for the gameScore field.
func (r *queryResolver) GameScore(ctx context.Context, uid string, gid string) (*model.Score, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, uid, false)
	if err != nil {
		_, err = checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
		if err != nil {
			return nil, err
		}
	}

	resp, err := r.RpcClients.GameClient.GetScore(ctx, &game_pb.GetScoreRequest{
		GameID:    gid,
		AccountID: uid,
	})
	if err != nil {
		slog.Error("get score error", "err", err)
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrScoreNotFound
			}
		}
		return nil, ErrInternal
	}

	return fromPbScore(resp.Score), err
}

// GameScores is the resolver for the gameScores field.
func (r *queryResolver) GameScores(ctx context.Context, gid string, page int, pageSize int) ([]*model.Score, error) {
	_, err := checkUserLogin(r.RpcClients.UserClient, r.sessionManager, ctx, "", true)
	if err != nil {
		return nil, err
	}

	resp, err := r.RpcClients.GameClient.GetScorePage(ctx,
		&game_pb.GetScorePageRequest{
			GameID:   gid,
			Page:     int64(page),
			PageSize: int64(pageSize),
		})
	if err != nil {
		slog.Error("get score page error", "err", err)
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrScoreNotFound
			}
		}
		return nil, ErrInternal
	}

	scores := make([]*model.Score, len(resp.GetScores()))
	for id, s := range resp.GetScores() {
		scores[id] = fromPbScore(s)
	}

	return scores, err
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func fromPbScore(s *game_pb.Score) *model.Score {
	return &model.Score{
		UserID:     s.AccountID,
		Score:      int(s.GetScore()),
		SignUpTime: pkg.Int64ToString(s.SignUpTime),
	}
}

var (
	ErrScoreNotFound          = errors.New("成绩不存在，请检查输入用户ID与比赛ID是否正确")
	ErrGameNotFound           = errors.New("比赛不存在")
	ErrGameInvalid            = errors.New("比赛信息无效")
	ErrSignUpOperationExpired = errors.New("操作时效已过")
	ErrGameAlreadySignUp      = errors.New("该用户已经报名比赛")
	ErrGameNotSignUp          = errors.New("该用户未报名比赛")
	ErrGameIDNotProvide       = errors.New("比赛ID未提供")
)

func fromPbGame(g *game_pb.Game) *model.Game {
	return &model.Game{
		ID:           g.GameID,
		CreaterID:    g.AccountID,
		Title:        g.Title,
		Description:  g.Desc,
		StartTime:    pkg.Int64ToString(g.GetStartTime()),
		EndTime:      pkg.Int64ToString(g.GetEndTime()),
		CreateTime:   pkg.Int64ToString(g.GetCreateTime()),
		QuestionList: fromPbGameQuestionList(g.QuestionList),
	}
}
func toPbGameQuestionList(questions []*model.GameQuestionInput) []*game_pb.GameQuestion {
	pbQuestions := make([]*game_pb.GameQuestion, len(questions))
	for i, q := range questions {
		pbQuestions[i] = &game_pb.GameQuestion{
			QuestionID: q.QuestionID,
			Score:      int64(q.Score),
		}
	}
	return pbQuestions
}
func fromPbGameQuestionList(questions []*game_pb.GameQuestion) []*model.GameQuestion {
	ret := make([]*model.GameQuestion, len(questions))
	for i, q := range questions {
		ret[i] = &model.GameQuestion{
			QuestionID: q.QuestionID,
			Score:      int(q.Score),
		}
	}
	return ret
}
func findGame(client game_pb.GameServiceClient,
	ctx context.Context, gid string) (*model.Game, error) {
	resp, err := client.GetGame(ctx, &game_pb.GetGameRequest{GameID: gid})
	if err != nil {
		if status, ok := gstatus.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return nil, ErrGameNotFound
			}
		}
		return nil, ErrInternal
	}

	return &model.Game{
		ID:          resp.Game.GameID,
		Title:       resp.Game.Title,
		Description: resp.Game.Desc,
		StartTime:   pkg.Int64ToString(resp.Game.GetStartTime()),
		EndTime:     pkg.Int64ToString(resp.Game.GetEndTime()),
	}, err
}
func stringToInt64(t string) (int64, error) {
	return strconv.ParseInt(t, 10, 64)
}
