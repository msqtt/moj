package policy

import (
	"errors"
	"moj/domain/captcha"
	domain_err "moj/domain/pkg/error"
)

type SendCaptchaEmailPolicy struct {
	emailService EmailService
}

var ErrFailedToSendEmail = errors.New("failed to send email")

func NewSendCaptchaEmailPolicy(emailService EmailService) *SendCaptchaEmailPolicy {
	return &SendCaptchaEmailPolicy{
		emailService: emailService,
	}
}

func (p *SendCaptchaEmailPolicy) OnEvent(event any) (err error) {
	switch v := event.(type) {
	case captcha.ChangePasswdCaptchaEvent:
		cmd := &CaptchaEmailCmd{
			Email:    v.Email,
			IpAddr:   v.IpAddr,
			Time:     v.CreateTime,
			Duration: v.Duration,
		}
		err = p.emailService.SendChangePasswordEmail(cmd)
	case captcha.RegisterCaptchaEvent:
		cmd := &CaptchaEmailCmd{
			Email:    v.Email,
			IpAddr:   v.IpAddr,
			Time:     v.CreateTime,
			Duration: v.Duration,
		}
		err = p.emailService.SendRegisterEmail(cmd)
	default:
		return domain_err.ErrEventTypeInvalid
	}
	if err != nil {
		err = errors.Join(ErrFailedToSendEmail, err)
	}
	return
}
