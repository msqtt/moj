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
	err := s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
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
	err := s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
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

// DeleteUser implements user_pb.UserServiceServer.
func (s *Server) DeleteUser(ctx context.Context, req *user_pb.DeleteUserRequest) (
	resp *user_pb.DeleteUserResponse, err error) {
	slog.Debug("delete account request", "req", req)

	cmd := account.DeleteAccountCmd{
		AccountID: req.AccountID,
		Time:      time.Now().Unix(),
	}

	slog.Info("invoking delete account command", "cmd", cmd)
	err = s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
		return s.deleteAccountCmdHandler.Handle(eq, cmd)
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
	err = s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
		return s.setAdminAccountCmdHandler.Handle(eq, cmd)
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

	err = s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
		return s.setStatusAccountCmdHandler.Handle(eq, cmd)
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
	err = s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
		return s.modifyInfoAccountCmdHandler.Handle(eq, cmd)
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
	err = s.commandInvoker.Invoke(func(eq queue.EventQueue) error {
		return s.changePasswdService.Handle(eq, cmd)
	})
	if err != nil {
		slog.Error("failed to invoke change user password command", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.ChangeUserPasswordResponse{Time: time.Now().Unix()}
	return
}

func (s *Server) GetUserInfo(ctx context.Context,
	req *user_pb.GetUserInfoRequest) (resp *user_pb.GetUserInfoResponse, err error) {
	slog.Debug("get user info request", "req", req)

	view, err := s.accountViewDAO.FindByAccountID(req.AccountID)
	if err != nil {
		slog.Error("failed to find user info", "err", err)
		err = responseStatusError(err)
	}
	resp = &user_pb.GetUserInfoResponse{
		UserInfo: &user_pb.UserInfo{
			AccountID:            view.AccountID,
			Email:                view.Email,
			AvatarLink:           view.AvatarLink,
			NickName:             view.NickName,
			Enabled:              view.Enabled,
			IsAdmin:              view.IsAdmin,
			LastLoginTime:        view.LastLoginTime.Unix(),
			LastLoginIPAddr:      view.LastLoginIPAddr,
			LastLoginDevice:      view.LastLoginDevice,
			LastPasswdChangeTime: view.LastPasswdChangeTime.Unix(),
			RegisterTime:         view.RegisterTime.Unix(),
			DeleteTime:           view.DeleteTime.Unix(),
		},
	}
	return

}

func accountViewModelToUserView(m *db.AccountViewModel) *user_pb.UserView {
	return &user_pb.UserView{
		AccountID:    m.AccountID,
		Email:        m.Email,
		AvatarLink:   m.AvatarLink,
		NickName:     m.NickName,
		Enabled:      m.Enabled,
		IsAdmin:      m.IsAdmin,
		RegisterTime: m.RegisterTime.Unix(),
		DeleteTime:   m.DeleteTime.Unix(),
	}

}

func (s *Server) GetUserPage(ctx context.Context,
	req *user_pb.GetUserPageRequest) (resp *user_pb.GetUserPageResponse, err error) {
	slog.Debug("get user page request", "req", req)

	views, err := s.accountViewDAO.FindByPage(int(req.GetPageSize()), req.GetCursor(), req.GetWord())
	if err != nil {
		slog.Error("failed to find user page", "err", err)
		err = responseStatusError(err)
	}

	userViews := make([]*user_pb.UserView, 0)
	for _, view := range views {
		userViews = append(userViews, accountViewModelToUserView(view))
	}
	var nexCursor string
	if len(userViews) > 0 {
		nexCursor = views[len(userViews)-1].AccountID
	}
	resp = &user_pb.GetUserPageResponse{
		UserView:   userViews,
		NextCursor: nexCursor,
	}
	return
}

var _ user_pb.UserServiceServer = (*Server)(nil)
