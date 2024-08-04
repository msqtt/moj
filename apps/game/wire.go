//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/game/db"
	"moj/apps/game/domain"
	"moj/apps/game/etc"
	"moj/apps/game/listener"
	"moj/apps/game/mq"
	"moj/apps/game/svc"
	"moj/domain/game"

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
