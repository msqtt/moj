//go:build wireinject
// +build wireinject

package main

import (
	"moj/user/db"
	"moj/user/domain"
	"moj/user/etc"
	"moj/user/listener"
	"moj/user/mail"
	service "moj/user/svc"
	"moj/domain/account"
	"moj/domain/captcha"
	"moj/domain/policy"
	svc_account "moj/domain/service/account"

	"github.com/google/wire"
)

func ProvideEventDispatcher(
	emailPolicy *policy.SendCaptchaEmailPolicy,
	avDao db.AccountViewDAO,
) domain.EventDispatcher {
	return domain.NewSyncEventDispatcher(
		listener.NewEmailListener(emailPolicy),
		listener.NewAccountViewListener(avDao),
	)
}

var (
	serverSet = wire.NewSet(
		service.NewServer,
		account.NewLoginAccountCmdHandler,
		account.NewCreateAccountCmdHandler,
		account.NewSetAdminAccountCmdHandler,
		account.NewSetStatusAccountCmdHandler,
		account.NewDeleteAccountCmdHandler,
		account.NewModifyInfoAccountCmdHandler,
		account.NewChangePasswdAccountCmdHandler,
		svc_account.NewAccountRegisterService,
		svc_account.NewChangePasswdService,
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
		policy.NewSendCaptchaEmailPolicy,
		mail.NewEmailServer,
		etc.NewAppConfig,
		ProvideEventDispatcher,
	)
)

var providers = wire.NewSet(serverSet, dbSet, otherSet)

func InitializeApplication() *App {
	wire.Build(NewApp, providers)
	return nil
}
