//go:build wireinject
// +build wireinject

package main

import (
	"moj/question/db"
	"moj/question/domain"
	"moj/question/etc"
	"moj/question/svc"
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
