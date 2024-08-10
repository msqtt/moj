package svc

import (
	"context"
	"log/slog"
	"time"

	"moj/apps/user/db"
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
	err := s.commandInvoker.Invoke(ctx, func(ctx context.Context, eq queue.EventQueue) error {
		return s.loginAccountCmdHandler.Handle(ctx, eq, cmd)
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
	err := s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) error {
		return s.accountRegisterService.Handle(ctx1, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke register command", "err", err)
		return nil, responseStatusError(err)
	}
	return &user_pb.RegisterResponse{
		Time: time.Now().Unix(),
	}, nil
}

// DeleteUser implements user_pb.UserServiceServer.
func (s *Server) DeleteUser(ctx context.Context, req *user_pb.DeleteUserRequest) (
	resp *user_pb.DeleteUserResponse, err error) {
	slog.Debug("delete account request", "req", req)

	cmd := account.DeleteAccountCmd{
		AccountID: req.AccountID,
		Time:      time.Now().Unix(),
	}

	slog.Info("invoking delete account command", "cmd", cmd)
	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) error {
		return s.deleteAccountCmdHandler.Handle(ctx1, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke delete account command", "err", err)
		err = responseStatusError(err)
		return
	}
	resp = &user_pb.DeleteUserResponse{Time: time.Now().Unix()}
	return
}

// SetAdmin implements user_pb.UserServiceServer.
func (s *Server) SetAdmin(ctx context.Context, req *user_pb.SetAdminRequest) (
	resp *user_pb.SetAdminResponse, err error) {
	slog.Debug("set admin request", "req", req)

	cmd := account.SetAdminAccountCmd{
		AccountID: req.AccountID,
		IsAdmin:   req.IsAdmin,
	}

	slog.Info("invoking set admin command", "cmd", cmd)
	err = s.commandInvoker.Invoke(ctx, func(ctx1 context.Context, eq queue.EventQueue) error {
		return s.setAdminAccountCmdHandler.Handle(ctx1, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke set admin command", "err", err)
		err = responseStatusError(err)
		return
	}
	resp = &user_pb.SetAdminResponse{Time: time.Now().Unix()}
	return
}

// SetStatus implements user_pb.UserServiceServer.
func (s *Server) SetStatus(ctx context.Context, req *user_pb.SetStatusRequest) (
	resp *user_pb.SetStatusResponse, err error) {
	slog.Debug("set status request", "req", req)

	cmd := account.SetStatusAccountCmd{
		AccountID: req.AccountID,
		Enabled:   req.Enabled,
	}

	slog.Info("invoking set status command", "cmd", cmd)

	err = s.commandInvoker.Invoke(ctx, func(ctx context.Context, eq queue.EventQueue) error {
		return s.setStatusAccountCmdHandler.Handle(ctx, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke set status command", "err", err)
		err = responseStatusError(err)
		return
	}
	resp = &user_pb.SetStatusResponse{Time: time.Now().Unix()}
	return
}

// UpdateUserInfo implements user_pb.UserServiceServer.
func (s *Server) UpdateUserInfo(ctx context.Context,
	req *user_pb.UpdateUserInfoRequest) (resp *user_pb.UpdateUserInfoResponse, err error) {
	slog.Debug("update user info request", "req", req)

	cmd := account.ModifyInfoAccountCmd{
		AccountID:  req.AccountID,
		NickName:   req.NickName,
		AvatarLink: req.AvatarLink,
	}

	slog.Info("invoking update user info command", "cmd", cmd)
	err = s.commandInvoker.Invoke(ctx, func(ctx context.Context, eq queue.EventQueue) error {
		return s.modifyInfoAccountCmdHandler.Handle(ctx, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke update user info command", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.UpdateUserInfoResponse{Time: time.Now().Unix()}
	return
}

func (s *Server) ChangeUserPassword(ctx context.Context,
	req *user_pb.ChangeUserPasswordRequest) (resp *user_pb.ChangeUserPasswordResponse, err error) {
	slog.Debug("change user password request", "req", req)

	cmd := svc_account.ChangePasswdCmd{
		AccountID: req.GetAccountID(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		Captcha:   req.GetCaptcha(),
		Time:      time.Now().Unix(),
	}

	slog.Info("invoking change user password command", "cmd", cmd)
	err = s.commandInvoker.Invoke(ctx, func(ctx context.Context, eq queue.EventQueue) error {
		return s.changePasswdService.Handle(ctx, eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke change user password command", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.ChangeUserPasswordResponse{Time: time.Now().Unix()}
	return
}

func (s *Server) GetUser(ctx context.Context,
	req *user_pb.GetUserRequest) (resp *user_pb.GetUserResponse, err error) {
	slog.Debug("get user info request", "req", req)

	view, err := s.accountViewDAO.FindByAccountID(ctx, req.AccountID)
	if err != nil {
		slog.Error("failed to find user info", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.GetUserResponse{
		User: accountViewModelToUsers(view)}
	return

}

func accountViewModelToUsers(m *db.AccountViewModel) *user_pb.User {
	return &user_pb.User{
		AccountID:            m.AccountID,
		Email:                m.Email,
		AvatarLink:           m.AvatarLink,
		NickName:             m.NickName,
		Enabled:              m.Enabled,
		IsAdmin:              m.IsAdmin,
		LastLoginTime:        m.LastLoginTime.Unix(),
		LastLoginIPAddr:      m.LastLoginIPAddr,
		LastLoginDevice:      m.LastLoginDevice,
		LastPasswdChangeTime: m.LastPasswdChangeTime.Unix(),
		RegisterTime:         m.RegisterTime.Unix(),
		DeleteTime:           m.DeleteTime.Unix(),
	}

}

func (s *Server) GetUserPage(ctx context.Context,
	req *user_pb.GetUserPageRequest) (resp *user_pb.GetUserPageResponse, err error) {
	slog.Debug("get user page request", "req", req)

	m := make(map[string]any)
	if req.FilterOptions != nil {
		if req.FilterOptions.Word != nil {
			m["word"] = req.FilterOptions.Word
		}
		if req.FilterOptions.Enabled != nil {
			m["enabled"] = req.FilterOptions.Enabled
		}
		if req.FilterOptions.IsAdmin != nil {
			m["is_admin"] = req.FilterOptions.IsAdmin
		}
	}

	views, err := s.accountViewDAO.FindByPage(ctx, int(req.GetPageSize()), req.GetCursor(), m)
	if err != nil {
		slog.Error("failed to find user page", "err", err)
		err = responseStatusError(err)
	}

	users := make([]*user_pb.User, 0)
	for _, view := range views {
		users = append(users, accountViewModelToUsers(view))
	}
	var nexCursor string
	if len(users) > 0 {
		nexCursor = views[len(users)-1].AccountID
	}
	resp = &user_pb.GetUserPageResponse{
		Users:      users,
		NextCursor: nexCursor,
	}
	return
}

var _ user_pb.UserServiceServer = (*Server)(nil)
