package svc

import (
	"moj/apps/user/db"
	"moj/apps/user/domain"
	"moj/apps/user/etc"
	user_pb "moj/apps/user/rpc"
	"moj/domain/account"
	"moj/domain/captcha"
	svc_account "moj/domain/service/account"
)

type Server struct {
	user_pb.UnimplementedUserServiceServer
	user_pb.UnimplementedCaptchaServiceServer
	loginAccountCmdHandler              *account.LoginAccountCmdHandler
	accountRegisterService              *svc_account.AccountRegisterService
	createChangePasswdCaptchaCmdHandler *captcha.CreateChangePasswdCaptchaCmdHandler
	createRegisterCaptchaCmdHandler     *captcha.CreateRegisterCaptchaCmdHandler
	commandInvoker                      domain.CommandInvoker
	conf                                *etc.Config
	accountViewDAO                      db.AccountViewDAO
}

func NewServer(
	loginAccountCmdHandler *account.LoginAccountCmdHandler,
	accountRegisterService *svc_account.AccountRegisterService,
	createChangePasswdCaptchaCmdHandler *captcha.CreateChangePasswdCaptchaCmdHandler,
	createRegisterCaptchaCmdHandler *captcha.CreateRegisterCaptchaCmdHandler,
	commandInvoker domain.CommandInvoker,
	conf *etc.Config,
	accountViewDAO db.AccountViewDAO,
) *Server {
	return &Server{
		loginAccountCmdHandler:              loginAccountCmdHandler,
		accountRegisterService:              accountRegisterService,
		createChangePasswdCaptchaCmdHandler: createChangePasswdCaptchaCmdHandler,
		createRegisterCaptchaCmdHandler:     createRegisterCaptchaCmdHandler,
		commandInvoker:                      commandInvoker,
		conf:                                conf,
		accountViewDAO:                      accountViewDAO,
	}
}
