package svc

import (
	"moj/domain/game"
	"moj/domain/record"
	"moj/game/db"
	"moj/game/domain"
	game_pb "moj/game/rpc"
)

type Server struct {
	game_pb.UnimplementedGameServiceServer
	gameViewDao                db.GameViewDao
	signUpScoreDao             db.SignUpScoreDao
	commandInvoker             domain.CommandInvoker
	gameRepository             game.GameRepository
	createGameCmdHandler       *game.CreateGameCmdHandler
	modifyGameCmdHandler       *game.ModifyGameCmdHandler
	signUpGameCmdHandler       *game.SignupGameCmdHandler
	cancelSignUpGameCmdHandler *game.CancelSignUpGameCmdHandler
	recordRepository           record.RecordRepository
}

func NewServer(
	gameViewDao db.GameViewDao,
	signUpScoreDao db.SignUpScoreDao,
	commandInvoker domain.CommandInvoker,
	gameRepository game.GameRepository,
	createGameCmdHandler *game.CreateGameCmdHandler,
	modifyGameCmdHandler *game.ModifyGameCmdHandler,
	signUpGameCmdHandler *game.SignupGameCmdHandler,
	cancelSignUpGameCmdHandler *game.CancelSignUpGameCmdHandler,
	recordRepository record.RecordRepository,
) *Server {
	return &Server{
		gameViewDao:                gameViewDao,
		signUpScoreDao:             signUpScoreDao,
		commandInvoker:             commandInvoker,
		gameRepository:             gameRepository,
		createGameCmdHandler:       createGameCmdHandler,
		modifyGameCmdHandler:       modifyGameCmdHandler,
		signUpGameCmdHandler:       signUpGameCmdHandler,
		cancelSignUpGameCmdHandler: cancelSignUpGameCmdHandler,
		recordRepository:           recordRepository,
	}
}
