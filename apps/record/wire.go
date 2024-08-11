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
	"moj/domain/question"
	"moj/domain/record"
	"time"

	"github.com/google/wire"
)

func provideDispatcher(
	conf *etc.Config,
	dailyTaskViewDao db.DailyTaskViewDao,
	passedQuestionViewDao db.PassedQuestionViewDao,
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
