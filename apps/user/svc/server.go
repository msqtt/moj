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
	conf                                *etc.Config
	commandInvoker                      domain.CommandInvoker
	accountViewDAO                      db.AccountViewDAO
	loginAccountCmdHandler              *account.LoginAccountCmdHandler
	setAdminAccountCmdHandler           *account.SetAdminAccountCmdHandler
	setStatusAccountCmdHandler          *account.SetStatusAccountCmdHandler
	deleteAccountCmdHandler             *account.DeleteAccountCmdHandler
	modifyInfoAccountCmdHandler         *account.ModifyInfoAccountCmdHandler
	accountRegisterService              *svc_account.AccountRegisterService
	changePasswdService                 *svc_account.ChangePasswdService
	createChangePasswdCaptchaCmdHandler *captcha.CreateChangePasswdCaptchaCmdHandler
	createRegisterCaptchaCmdHandler     *captcha.CreateRegisterCaptchaCmdHandler
}

func NewServer(
	conf *etc.Config,
	commandInvoker domain.CommandInvoker,
	accountViewDAO db.AccountViewDAO,
	loginAccountCmdHandler *account.LoginAccountCmdHandler,
	setAdminAccountCmdHandler *account.SetAdminAccountCmdHandler,
	setStatusAccountCmdHandler *account.SetStatusAccountCmdHandler,
	deleteAccountCmdHandler *account.DeleteAccountCmdHandler,
	modifyInfoAccountCmdHandler *account.ModifyInfoAccountCmdHandler,
	accountRegisterService *svc_account.AccountRegisterService,
	changePasswdService *svc_account.ChangePasswdService,
	createChangePasswdCaptchaCmdHandler *captcha.CreateChangePasswdCaptchaCmdHandler,
	createRegisterCaptchaCmdHandler *captcha.CreateRegisterCaptchaCmdHandler,
) *Server {
	return &Server{
		conf:                                conf,
		commandInvoker:                      commandInvoker,
		accountViewDAO:                      accountViewDAO,
		loginAccountCmdHandler:              loginAccountCmdHandler,
		setAdminAccountCmdHandler:           setAdminAccountCmdHandler,
		setStatusAccountCmdHandler:          setStatusAccountCmdHandler,
		deleteAccountCmdHandler:             deleteAccountCmdHandler,
		modifyInfoAccountCmdHandler:         modifyInfoAccountCmdHandler,
		accountRegisterService:              accountRegisterService,
		changePasswdService:                 changePasswdService,
		createChangePasswdCaptchaCmdHandler: createChangePasswdCaptchaCmdHandler,
		createRegisterCaptchaCmdHandler:     createRegisterCaptchaCmdHandler,
	}
}
