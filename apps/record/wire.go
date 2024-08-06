//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/record/db"
	"moj/apps/record/domain"
	"moj/apps/record/etc"
	"moj/apps/record/listener"
	"moj/apps/record/mq/consumer"
	"moj/apps/record/mq/producer"
	"moj/apps/record/schedule"
	"moj/apps/record/svc"
	"moj/domain/record"
	"time"

	"github.com/google/wire"
)

func provideDispatcher(conf *etc.Config) domain.EventDispatcher {
	return domain.NewSyncAndAsyncEventDispatcher(
		[]listener.Listener{},
		[]producer.Producer{
			producer.NewExecuteJudgeProducer(conf),
			producer.NewRecordGameScoreProducer(conf),
		})
}

func provideScheduleTask(conf *etc.Config, dao db.RecordViewDao) []*schedule.TikerTasker {
	return []*schedule.TikerTasker{
		schedule.NewTickerTasker(
			time.NewTicker(conf.ScheduleRedoSubmitDuration),
			schedule.NewExecuteJudgementTask(dao, producer.NewExecuteJudgeProducer(conf)),
		),
	}
}

var providers = wire.NewSet(
	svc.NewServer,

	record.NewModifyRecordCmdHandler,
	record.NewSubmitRecordCmdHandler,
	domain.NewMongoDBRecordRepository,
	domain.NewTransactionCommandInvoker,
	producer.NewRecordGameScoreProducer,

	provideDispatcher,
	provideScheduleTask,

	consumer.NewNsqFinishRecordConsumer,

	db.NewMongoDBTransactionManager,
	db.NewMongoDBRecordViewDao,
	db.NewMongoDB,
	etc.NewAppConfig,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
