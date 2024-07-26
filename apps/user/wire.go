//go:build wireinject
// +build wireinject

package main

import (
	"moj/apps/user/db"
	"moj/apps/user/domain"
	"moj/apps/user/etc"
	"moj/apps/user/listener"
	"moj/apps/user/mail"
	service "moj/apps/user/svc"
	"moj/domain/account"
	"moj/domain/captcha"
	"moj/domain/policy"
	svc_account "moj/domain/service/account"

	"github.com/google/wire"
)

func ProvideEventDispatcher(
	accountViewDAO db.AccountViewDAO,
	emailServer policy.EmailService,
) domain.EventDispatcher {
	return domain.NewSyncEventDispatcher(
		listener.NewAccountViewListener(accountViewDAO),
		policy.NewSendCaptchaEmailPolicy(emailServer),
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
