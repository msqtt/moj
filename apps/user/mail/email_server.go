package mail

import (
	"log/slog"
	"moj/user/etc"
	"moj/domain/policy"
)

type EmailServer struct {
	conf *etc.Config
}

func NewEmailServer(conf *etc.Config) policy.EmailService {
	return &EmailServer{
		conf: conf,
	}
}

// SendChangePasswordEmail implements policy.EmailService.
func (e *EmailServer) SendChangePasswordEmail(cmd *policy.CaptchaEmailCmd) error {
	slog.Warn("email server unimplemented")
	return nil
}

// SendRegisterEmail implements policy.EmailService.
func (e *EmailServer) SendRegisterEmail(cmd *policy.CaptchaEmailCmd) error {
	slog.Warn("email server unimplemented")
	return nil
}

var _ policy.EmailService = (*EmailServer)(nil)
