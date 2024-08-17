package svc

import (
	"context"
	"log/slog"
	game_pb "moj/apps/game/rpc"
	"moj/domain/game"
	"moj/domain/pkg/queue"
	"time"
)

var _ game_pb.GameServiceServer = (*Server)(nil)

func fromGameQuestionListPb(queList []*game_pb.GameQuestion) []game.GameQuestion {
	gameQuestionList := make([]game.GameQuestion, len(queList))
	for idx, que := range queList {
		gameQuestionList[idx] = game.GameQuestion{
			QuestionID: que.QuestionID,
			Score:      int(que.Score),
		}
	}
	return gameQuestionList
}

func toGameQuestionListPb(queList []game.GameQuestion) []*game_pb.GameQuestion {
	gameQuestionList := make([]*game_pb.GameQuestion, len(queList))
	for idx, que := range queList {
		gameQuestionList[idx] = &game_pb.GameQuestion{
			QuestionID: que.QuestionID,
			Score:      int64(que.Score),
		}
	}
	return gameQuestionList
}

func toSignUpPb(acc []game.SignUpAccount) []*game_pb.SignUpAccount {
	signUpAccountList := make([]*game_pb.SignUpAccount, len(acc))
	for idx, a := range acc {
		signUpAccountList[idx] = &game_pb.SignUpAccount{
			AccountID:  a.AccountID,
			SignUpTime: a.SignUpTime,
		}
	}
	return signUpAccountList
}

// CancelSignUpGame implements game_pb.GameServiceServer.
func (s *Server) CancelSignUpGame(ctx context.Context, req *game_pb.CancelSignUpGameRequest) (
	resp *game_pb.CancelSignUpGameResponse, err error) {
	slog.Debug("cancel sign up game", "req", req)

	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		err := s.cancelSignUpGameCmdHandler.Handle(ctx1, eq, game.CancelSignUpGameCmd{
			GameID:    req.GameID,
			AccountID: req.AccountID,
		})
		return nil, err
	})

	if err != nil {
		slog.Error("failed to invoke cancel sign up game command", "err", err)
		err = responseStatusError(err)
		return
	}

	resp = &game_pb.CancelSignUpGameResponse{
		Time: time.Now().Unix(),
	}
	return
}

// CreateGame implements game_pb.GameServiceServer.
func (s *Server) CreateGame(ctx context.Context, req *game_pb.CreateGameRequest) (
	resp *game_pb.CreateGameResponse, err error) {
	slog.Debug("create game", "req", req)

	gameQuestionList := fromGameQuestionListPb(req.QuestionList)

	cmd := game.CreateGameCmd{
		Title:        req.Title,
		Description:  req.Desc,
		AccountID:    req.AccountID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		QuestionList: gameQuestionList,
		Time:         time.Now().Unix(),
	}

	slog.Info("start to invoke create game command", "cmd", cmd)

	var gameID string
	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		gameID1, err1 := s.createGameCmdHandler.Handle(ctx1, eq, cmd)
		if gameID1 != nil {
			gameID = gameID1.(string)
		}
		return gameID1, err1
	})

	if err != nil {
		slog.Error("failed to invoke create game command", "err", err)
		err = responseStatusError(err)
		return
	}

	resp = &game_pb.CreateGameResponse{
		GameID: gameID,
		Time:   time.Now().Unix(),
	}

	return
}

// GetGame implements game_pb.GameServiceServer.
func (s *Server) GetGame(ctx context.Context, req *game_pb.GetGameRequest) (
	resp *game_pb.GetGameResponse, err error) {
	slog.Debug("get game info", "req", req)

	game, err := s.gameRepository.FindGameByID(ctx, req.GameID)
	if err != nil {
		slog.Error("failed to get game info", "err", err)
		err = responseStatusError(err)
		return
	}

	gameList := toGameQuestionListPb(game.QuestionList)
	signUpList := toSignUpPb(game.SignUpUserList)

	resp = &game_pb.GetGameResponse{
		Game: &game_pb.Game{
			GameID:            game.GameID,
			AccountID:         game.AccountID,
			Title:             game.Title,
			Desc:              game.Description,
			StartTime:         game.StartTime,
			EndTime:           game.EndTime,
			CreateTime:        game.CreateTime,
			QuestionList:      gameList,
			SignUpAccountList: signUpList,
		},
	}
	return
}

// GetGamePage implements game_pb.GameServiceServer.
func (s *Server) GetGamePage(ctx context.Context, req *game_pb.GetGamePageRequest) (
	resp *game_pb.GetGamePageResponse, err error) {
	slog.Debug("get game page", "req", req)

	f := make(map[string]any)
	if req.FilterOptions != nil {
		if req.FilterOptions.Word != nil {
			f["word"] = req.FilterOptions.Word
		}
		if req.FilterOptions.Time != nil {
			f["time"] = time.Unix(*req.FilterOptions.Time, 0)
		}
	}

	games, err := s.gameViewDao.FindGamePage(ctx, req.Cursor, int(req.PageSize), f)
	if err != nil {
		slog.Error("get question page error", "error", err)
		err = responseStatusError(err)
		return
	}
	var nextCursor string
	gameView := make([]*game_pb.Game, len(games))
	for i, q := range games {
		gameView[i] = &game_pb.Game{
			GameID:            q.ID.Hex(),
			AccountID:         q.AccountID,
			Title:             q.Title,
			Desc:              q.Description,
			StartTime:         q.StartTime.Unix(),
			EndTime:           q.EndTime.Unix(),
			CreateTime:        q.CreateTime.Unix(),
			QuestionList:      toGameQuestionListPb(q.ToAggreation().QuestionList),
			SignUpAccountList: toSignUpPb(q.ToAggreation().SignUpUserList),
		}
	}
	if len(gameView) > 0 {
		nextCursor = gameView[len(gameView)-1].GameID
	}

	resp = &game_pb.GetGamePageResponse{
		Games:  gameView,
		Cursor: nextCursor,
	}
	return
}

// GetScore implements game_pb.GameServiceServer.
func (s *Server) GetScore(ctx context.Context, req *game_pb.GetScoreRequest) (
	resp *game_pb.GetScoreResponse, err error) {
	slog.Debug("get score", "req", req)

	score, err := s.signUpScoreDao.FindByID(ctx, req.GameID, req.AccountID)
	if err != nil {
		slog.Error("failed to get score", "err", err)
		err = responseStatusError(err)
	}

	resp = &game_pb.GetScoreResponse{
		Score: &game_pb.Score{
			AccountID:  score.AccountID,
			Score:      int64(score.Score),
			SignUpTime: score.SignUpTime.Unix(),
		},
	}
	return
}

// GetScorePage implements game_pb.GameServiceServer.
func (s *Server) GetScorePage(ctx context.Context, req *game_pb.GetScorePageRequest) (
	resp *game_pb.GetScorePageResponse, err error) {
	slog.Debug("get score page", "req", req)

	scores, total, err := s.signUpScoreDao.FindPage(ctx, req.GameID, int(req.Page), int(req.PageSize))
	if err != nil {
		slog.Error("failed to get score page", "err", err)
		err = responseStatusError(err)
	}

	scoreView := make([]*game_pb.Score, len(scores))
	for i, score := range scores {
		scoreView[i] = &game_pb.Score{
			AccountID:  score.AccountID,
			Score:      int64(score.Score),
			SignUpTime: score.SignUpTime.Unix(),
		}
	}

	resp = &game_pb.GetScorePageResponse{
		Scores: scoreView,
		Total:  total,
	}
	return
}

// SignUpGame implements game_pb.GameServiceServer.
func (s *Server) SignUpGame(ctx context.Context, req *game_pb.SignUpGameRequest) (
	resp *game_pb.SignUpGameResponse, err error) {
	slog.Debug("sign up game", "req", req)

	cmd := game.SignUpGameCmd{
		GameID:    req.GameID,
		AccountID: req.AccountID,
		Time:      time.Now().Unix(),
	}

	slog.Info("start to invoke sign up game command", "cmd", cmd)

	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		return nil, s.signUpGameCmdHandler.Handle(ctx1, eq, cmd)
	})

	if err != nil {
		slog.Error("failed to invoke sign up game command", "err", err)
		err = responseStatusError(err)
	}
	resp = &game_pb.SignUpGameResponse{
		Time: time.Now().Unix(),
	}

	return
}

// UpdateGame implements game_pb.GameServiceServer.
func (s *Server) UpdateGame(ctx context.Context, req *game_pb.UpdateGameRequest) (
	resp *game_pb.UpdateGameResponse, err error) {
	slog.Debug("update game", "req", req)

	gameQuestionList := fromGameQuestionListPb(req.QuestionList)
	cmd := game.ModifyGameCmd{
		GameID:       req.GameID,
		Title:        req.Title,
		Descirption:  req.Desc,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		QuestionList: gameQuestionList,
	}

	slog.Info("start to invoke update game command", "cmd", cmd)

	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) (any, error) {
		return s.modifyGameCmdHandler.Handle(ctx1, eq, cmd)
	})

	if err != nil {
		slog.Error("failed to invoke update game command", "err", err)
		err = responseStatusError(err)
		return
	}
	resp = &game_pb.UpdateGameResponse{
		Time: time.Now().Unix(),
	}
	return
}

// DeleteGame implements game_pb.GameServiceServer.
func (s *Server) DeleteGame(ctx context.Context, req *game_pb.DeleteGameRequest) (
	resp *game_pb.DeleteGameResponse, err error) {
	slog.Debug("delete game", "req", req)

	err = s.gameViewDao.DeleteGame(ctx, req.GetGameID())
	if err != nil {
		slog.Error("failed to delete game", "err", err)
		err = responseStatusError(err)
		return
	}
	resp = &game_pb.DeleteGameResponse{
		Time: time.Now().Unix(),
	}
	return
}
