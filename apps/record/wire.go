//go:build wireinject
// +build wireinject

package main

import (
	"moj/record/db"
	"moj/record/domain"
	"moj/record/etc"
	"moj/record/listener"
	"moj/record/mq/consumer"
	"moj/record/mq/producer"
	"moj/record/schedule"
	"moj/record/svc"
	"moj/domain/question"
	"moj/domain/record"
	"time"

	"github.com/google/wire"
)

func provideDispatcher(
	conf *etc.Config,
	dailyTaskViewDao db.DailyTaskViewDao,
	passedQuestionViewDao db.PassQuestionViewDao,
	questionRepository question.QuestionRepository,
) domain.EventDispatcher {
	return domain.NewSyncAndAsyncEventDispatcher(
		[]listener.Listener{
			listener.NewDailyTaskViewListener(dailyTaskViewDao),
			listener.NewPassedQuestionViewListener(passedQuestionViewDao, questionRepository),
		},
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
	domain.NewRPCQuestionRepository,
	producer.NewRecordGameScoreProducer,

	provideDispatcher,
	provideScheduleTask,

	consumer.NewNsqFinishRecordConsumer,

	db.NewMongoDBTransactionManager,
	db.NewMongoDBRecordViewDao,
	db.NewMongoDBPassedQuestionViewDao,
	db.NewMongoDBDayTaskViewDao,
	db.NewMongoDB,
	etc.NewAppConfig,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
