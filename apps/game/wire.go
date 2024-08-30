//go:build wireinject
// +build wireinject

package main

import (
	"moj/domain/game"
	"moj/game/db"
	"moj/game/domain"
	"moj/game/etc"
	"moj/game/listener"
	"moj/game/mq"
	"moj/game/svc"

	"github.com/google/wire"
)

func provideDispatcher(supDao db.SignUpScoreDao) domain.EventDispatcher {
	return domain.NewSyncEventDispatcher(listener.NewSignUpScoreLisener(supDao))
}

var providers = wire.NewSet(
	svc.NewServer,
	game.NewSignUpGameCmdHandler,
	game.NewCancelSignUpGameCmdHandler,
	game.NewCreateGameCmdHandler,
	game.NewModifyGameCmdHandler,
	game.NewCalculateScoreCmdHandler,

	mq.NewNsqCalculateScoreConsumer,

	provideDispatcher,
	domain.NewSimpleEventQueue,
	domain.NewTransactionCommandInvoker,
	domain.NewMongoDBGameRepository,
	domain.NewRPCRecordRepository,
	db.NewMongoDBTransactionManager,
	db.NewMongoDBSignUpScoreDao,
	db.NewMongoDBGameViewDao,
	db.NewMongoDB,
	etc.NewAppConfig,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
