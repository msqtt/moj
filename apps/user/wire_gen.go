// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"moj/apps/user/db"
	"moj/apps/user/domain"
	"moj/apps/user/etc"
	"moj/apps/user/listener"
	"moj/apps/user/svc"
	"moj/domain/account"
	"moj/domain/captcha"
	account2 "moj/domain/service/account"
)

// Injectors from wire.go:

func InitializeApplication() *App {
	config := etc.NewAppConfig()
	mongoDB := db.NewMongoDB(config)
	accountRepository := domain.NewMongoDBAccountRepository(config, mongoDB)
	loginAccountCmdHandler := account.NewLoginAccountCmdHandler(accountRepository)
	cryptor := domain.NewBCryptor()
	createAccountCmdHandler := account.NewCreateAccountCmdHandler(accountRepository, cryptor)
	captchaRepository := domain.NewMongoDBCaptchaRepository(config, mongoDB)
	accountRegisterService := account2.NewAccountRegisterService(createAccountCmdHandler, captchaRepository, accountRepository)
	createChangePasswdCaptchaCmdHandler := captcha.NewCreateChangePasswdCaptchaCmdHandler(captchaRepository)
	createRegisterCaptchaCmdHandler := captcha.NewCreateRegisterCaptchaCmdHandler(captchaRepository)
	transactionManager := db.NewMongoDBTransactionManager(mongoDB)
	accountViewDAO := db.NewMongoDBAccountViewDAO(config, mongoDB)
	eventDispatcher := ProvideEventDispatcher(accountViewDAO)
	commandInvoker := domain.NewTransactionCommandInvoker(transactionManager, eventDispatcher)
	server := svc.NewServer(loginAccountCmdHandler, accountRegisterService, createChangePasswdCaptchaCmdHandler, createRegisterCaptchaCmdHandler, commandInvoker, config, accountViewDAO)
	app := NewApp(server, mongoDB, config)
	return app
}

// wire.go:

func ProvideEventDispatcher(accountViewDAO db.AccountViewDAO) domain.EventDispatcher {
	return domain.NewSyncEventDispatcher(listener.NewAccountViewListener(accountViewDAO))
}

var (
	serverSet = wire.NewSet(svc.NewServer, account.NewLoginAccountCmdHandler, account.NewCreateAccountCmdHandler, account2.NewAccountRegisterService, captcha.NewCreateChangePasswdCaptchaCmdHandler, captcha.NewCreateRegisterCaptchaCmdHandler)
	dbSet     = wire.NewSet(db.NewMongoDB, db.NewMongoDBTransactionManager, domain.NewMongoDBAccountRepository, domain.NewMongoDBCaptchaRepository, db.NewMongoDBAccountViewDAO)
	otherSet  = wire.NewSet(domain.NewBCryptor, domain.NewSimpleEventQueue, domain.NewTransactionCommandInvoker, etc.NewAppConfig, ProvideEventDispatcher)
)

var providers = wire.NewSet(serverSet, dbSet, otherSet)
