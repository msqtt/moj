package svc

import (
	"moj/user/db"
	"moj/user/domain"
	"moj/user/etc"
	user_pb "moj/user/rpc"
	"moj/domain/account"
	"moj/domain/captcha"
	"moj/domain/pkg/crypt"
	svc_account "moj/domain/service/account"
)

type Server struct {
	user_pb.UnimplementedUserServiceServer
	user_pb.UnimplementedCaptchaServiceServer
	conf                                *etc.Config
	commandInvoker                      domain.CommandInvoker
	cryptor                             crypt.Cryptor
	accountViewDAO                      db.AccountViewDAO
	accountRepository                   account.AccountRepository
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
	cryptor crypt.Cryptor,
	accountViewDAO db.AccountViewDAO,
	accountRepository account.AccountRepository,
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
		cryptor:                             cryptor,
		accountViewDAO:                      accountViewDAO,
		accountRepository:                   accountRepository,
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
