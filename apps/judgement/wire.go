//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/judgement/db"
	"moj/apps/judgement/domain"
	"moj/apps/judgement/etc"
	"moj/apps/judgement/listener"
	"moj/apps/judgement/mq/consumer"
	"moj/apps/judgement/mq/producer"
	"moj/apps/judgement/svc"
	"moj/domain/judgement"

	"github.com/google/wire"
)

func provideDispatcher(conf *etc.Config) domain.EventDispatcher {
	return domain.NewSyncAndAsyncEventDispatcher(
		[]listener.Listener{},
		[]producer.Producer{producer.NewNsqModifyRecordProducer(conf)})
}

var providers = wire.NewSet(
	svc.NewServer,
	domain.NewTransactionCommandInvoker,
	domain.NewMongoDBJudementRepository,
	domain.NewSbJudger,
	domain.NewRPCQuestionRepository,
	domain.NewMinioCaseReader,
	judgement.NewExecutionCmdHandler,
	consumer.NewNsqExecuteJudgementConsumer,

	provideDispatcher,

	db.NewMongoDBTransactionManager,
	db.NewMongoDB,
	etc.NewAppConfig,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
