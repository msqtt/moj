//go:build wireinject
// +build wireinject

package main

import (
	"moj/domain/judgement"
	"moj/judgement/db"
	"moj/judgement/domain"
	"moj/judgement/etc"
	"moj/judgement/listener"
	"moj/judgement/mq/consumer"
	"moj/judgement/mq/producer"
	"moj/judgement/svc"

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
	domain.NewCacheCaseReader,
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
