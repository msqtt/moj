package svc

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/user/db"
	user_pb "moj/apps/user/rpc"
	"moj/domain/captcha"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
	"time"
)

var ErrAlreadyRegistered = errors.Join(domain_err.ErrDuplicated,
	errors.New("the email already been registered"))

func (s *Server) SendRegisterCaptcha(ctx context.Context, req *user_pb.SendRegisterCaptchaRequest) (
	resp *user_pb.SendRegisterCaptchaResponse, err error) {
	slog.Debug("send register request", "req", req)

	// check latest account by email
	_, err = s.accountViewDAO.FindLatestByEmail(req.GetEmail())
	if err == nil {
		slog.Info("the email already been registered", "email", req.Email)
		return nil, responseStatusError(ErrAlreadyRegistered)
	}
	if !errors.Is(err, db.ErrAccountViewNotFound) {
		slog.Error("failed to find account by email", "err", err)
		return nil, responseStatusError(err)
	}
	// account not exist

	cmd := &captcha.CreateRegisterCaptchaCmd{
		Email:    req.GetEmail(),
		IpAddr:   req.GetIpAddr(),
		Time:     time.Now().Unix(),
		Duration: s.conf.CaptchaLiveDuration,
	}
	slog.Info("invoking send register captcha command", "cmd", cmd)
	err = s.commandInvoker.Invoker(func(eq queue.EventQueue) error {
		return s.createRegisterCaptchaCmdHandler.Handle(eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke send register captcha command", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.SendRegisterCaptchaResponse{
		Time: time.Now().Unix(),
	}
	return
}

func (s *Server) SendChangePasswdCaptcha(ctx context.Context, req *user_pb.SendChangePasswdCaptchaRequest) (
	resp *user_pb.SendChangePasswdCaptchaResponse, err error) {
	slog.Debug("send change password request", "req", req)

	// check the account by email
	_, err = s.accountViewDAO.FindLatestByEmail(req.GetEmail())
	if err != nil {
		slog.Error("failed to find account by email", "err", err)
		if errors.Is(err, db.ErrAccountViewNotFound) {
			return nil, responseStatusError(err)
		}
		return nil, responseStatusError(err)
	}
	// account exist

	cmd := &captcha.CreateChangePasswdCaptchaCmd{
		AccountID: req.GetAccountID(),
		Email:     req.GetEmail(),
		IpAddr:    req.GetIpAddr(),
		Time:      time.Now().Unix(),
		Duration:  s.conf.CaptchaLiveDuration,
	}
	slog.Info("invoking send change password captcha command", "cmd", cmd)
	err = s.commandInvoker.Invoker(func(eq queue.EventQueue) error {
		return s.createChangePasswdCaptchaCmdHandler.Handle(eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke send change password  captcha command", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.SendChangePasswdCaptchaResponse{
		Time: time.Now().Unix(),
	}
	return
}

var _ user_pb.CaptchaServiceServer = (*Server)(nil)
