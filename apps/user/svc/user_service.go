package svc

import (
	"context"
	"log/slog"
	"time"

	user_pb "moj/apps/user/rpc"
	"moj/domain/account"
	"moj/domain/pkg/queue"
	svc_account "moj/domain/service/account"
)

// Login implements user_pb.UserServiceServer.
func (s *Server) Login(ctx context.Context, req *user_pb.LoginRequest) (
	*user_pb.LoginResponse, error) {
	slog.Debug("login request", "req", req)
	cmd := account.LoginAccountCmd{
		AccountID: req.GetAccountID(),
		Device:    req.GetDevice(),
		IPAddr:    req.GetIpAddr(),
		Time:      time.Now().Unix(),
	}

	slog.Info("invoking login command", "cmd", cmd)
	err := s.commandInvoker.Invoker(func(eq queue.EventQueue) error {
		return s.loginAccountCmdHandler.Handle(eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke login command", "err", err)
		return nil, responseStatusError(err)
	}
	return &user_pb.LoginResponse{
		Time: time.Now().Unix(),
	}, nil
}

// Register implements user_pb.UserServiceServer.
func (s *Server) Register(ctx context.Context, req *user_pb.RegisterRequest) (
	*user_pb.RegisterResponse, error) {
	slog.Debug("register request", "req", req)

	cmd := svc_account.RegisterCmd{
		Email:    req.Email,
		NickName: req.NickName,
		Password: req.Password,
		Captcha:  req.Captcha,
		Time:     time.Now().Unix(),
	}

	slog.Info("invoking register command", "cmd", cmd)
	err := s.commandInvoker.Invoker(func(eq queue.EventQueue) error {
		return s.accountRegisterService.Handle(eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke register command", "err", err)
		return nil, responseStatusError(err)
	}
	return &user_pb.RegisterResponse{
		Time: time.Now().Unix(),
	}, nil
}

var _ user_pb.UserServiceServer = (*Server)(nil)
