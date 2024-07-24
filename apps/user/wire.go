//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/user/db"
	"moj/apps/user/domain"
	"moj/apps/user/etc"
	"moj/apps/user/listener"
	service "moj/apps/user/svc"
	"moj/domain/account"
	"moj/domain/captcha"
	svc_account "moj/domain/service/account"

	"github.com/google/wire"
)

func ProvideEventDispatcher(accountViewDAO db.AccountViewDAO) domain.EventDispatcher {
	return domain.NewSyncEventDispatcher(listener.NewAccountViewListener(accountViewDAO))
}

var (
	serverSet = wire.NewSet(
		service.NewServer,
		account.NewLoginAccountCmdHandler,
		account.NewCreateAccountCmdHandler,
		svc_account.NewAccountRegisterService,
		captcha.NewCreateChangePasswdCaptchaCmdHandler,
		captcha.NewCreateRegisterCaptchaCmdHandler,
	)
	dbSet = wire.NewSet(
		db.NewMongoDB,
		db.NewMongoDBTransactionManager,
		domain.NewMongoDBAccountRepository,
		domain.NewMongoDBCaptchaRepository,
		db.NewMongoDBAccountViewDAO,
	)
	otherSet = wire.NewSet(
		domain.NewBCryptor,
		domain.NewSimpleEventQueue,
		domain.NewTransactionCommandInvoker,
		etc.NewAppConfig,
		ProvideEventDispatcher,
	)
)

var providers = wire.NewSet(serverSet, dbSet, otherSet)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
