// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
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
)

// Injectors from wire.go:

func InitializeApplication() *App {
	config := etc.NewAppConfig()
	mongoDB := db.NewMongoDB(config)
	transactionManager := db.NewMongoDBTransactionManager(mongoDB)
	eventDispatcher := provideDispatcher(config)
	commandInvoker := domain.NewTransactionCommandInvoker(transactionManager, eventDispatcher)
	recordRepository := domain.NewMongoDBRecordRepository(mongoDB)
	modifyRecordCmdHandler := record.NewModifyRecordCmdHandler(recordRepository)
	submitRecordCmdHandler := record.NewSubmitRecordCmdHandler(recordRepository)
	recordViewDao := db.NewMongoDBRecordViewDao(mongoDB)
	server := svc.NewServer(commandInvoker, modifyRecordCmdHandler, submitRecordCmdHandler, recordRepository, recordViewDao)
	v := provideScheduleTask(config, recordViewDao)
	nsqFinishRecordConsumer := consumer.NewNsqFinishRecordConsumer(config, modifyRecordCmdHandler, eventDispatcher)
	app := NewApp(server, v, mongoDB, config, nsqFinishRecordConsumer)
	return app
}

// wire.go:

func provideDispatcher(conf *etc.Config) domain.EventDispatcher {
	return domain.NewSyncAndAsyncEventDispatcher(
		[]listener.Listener{},
		[]producer.Producer{producer.NewExecuteJudgeProducer(conf), producer.NewRecordGameScoreProducer(conf)})
}

func provideScheduleTask(conf *etc.Config, dao db.RecordViewDao) []*schedule.TikerTasker {
	return []*schedule.TikerTasker{schedule.NewTickerTasker(time.NewTicker(conf.ScheduleRedoSubmitDuration), schedule.NewExecuteJudgementTask(dao, producer.NewExecuteJudgeProducer(conf))),
	}
}

var providers = wire.NewSet(svc.NewServer, record.NewModifyRecordCmdHandler, record.NewSubmitRecordCmdHandler, domain.NewMongoDBRecordRepository, domain.NewTransactionCommandInvoker, producer.NewRecordGameScoreProducer, provideDispatcher,
	provideScheduleTask, consumer.NewNsqFinishRecordConsumer, db.NewMongoDBTransactionManager, db.NewMongoDBRecordViewDao, db.NewMongoDB, etc.NewAppConfig,
)