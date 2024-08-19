package svc

import (
	"moj/game/db"
	"moj/game/domain"
	game_pb "moj/game/rpc"
	"moj/domain/game"
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
	}
}
