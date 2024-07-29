//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/question/db"
	"moj/apps/question/domain"
	"moj/apps/question/etc"
	"moj/apps/question/svc"
	"moj/domain/question"

	"github.com/google/wire"
)

var providers = wire.NewSet(
	svc.NewServer,
	question.NewCreateQuestionCmdHandler,
	question.NewModifyQuestionCmdHandler,
	domain.NewMongoDBQuestionRepository,
	db.NewMongoDBQuestionDAO,

	etc.NewAppConfig,
	db.NewMongoDB,
)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
