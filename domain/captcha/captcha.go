package captcha

import (
	"errors"

	"moj/domain/pkg/common"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
)

type CaptchaType string

const (
	CaptchaTypeRegister     CaptchaType = "register"
	CaptchaTypeChangePasswd CaptchaType = "change password"
)

func (c CaptchaType) IsValid() bool {
	switch c {
	case CaptchaTypeRegister, CaptchaTypeChangePasswd:
		return true
	default:
		return false
	}
}

var (
	ErrInValidEmail       = errors.Join(domain_err.ErrInValided, errors.New("invalid email"))
	ErrInValidCaptchaType = errors.Join(domain_err.ErrInValided, errors.New("invalid captcha type"))
	ErrExpiredCaptcha     = errors.Join(domain_err.ErrExpired, errors.New("captcha expired"))
)

type Captcha struct {
	AccountID  string
	Code       string
	Email      string
	Type       CaptchaType
	IpAddr     string
	Enabled    bool
	Duration   int64
	CreateTime int64
	ExpireTime int64
}

func NewCaptcha(
	accountID string,
	email string, captchaType CaptchaType, ipAddr string, duration int64, time int64) (*Captcha, error) {

	if !common.IsEmail(email) {
		return nil, ErrInValidEmail
	}

	if !captchaType.IsValid() {
		return nil, ErrInValidCaptchaType
	}

	content := generateRandomCaptcha()

	return &Captcha{
		AccountID:  accountID,
		Code:       content,
		Email:      email,
		Type:       captchaType,
		IpAddr:     ipAddr,
		Enabled:    true,
		Duration:   duration,
		CreateTime: time,
		ExpireTime: time + duration,
	}, nil
}

func (c *Captcha) IsExpired(time int64) bool {
	if !c.Enabled {
		return false
	}
	return c.ExpireTime < time
}

func generateRandomCaptcha() string {
	return common.RandomStr(6)
}

func (c *Captcha) SetDisable() {
	c.Enabled = false
}

func (c *Captcha) sendEmail(queue queue.EventQueue) error {
	if !c.Enabled {
		return ErrExpiredCaptcha
	}

	var event any
	switch c.Type {
	case CaptchaTypeRegister:
		event = RegisterCaptchaEvent{
			Code:       c.Code,
			Email:      c.Email,
			IpAddr:     c.IpAddr,
			Duration:   c.Duration,
			CreateTime: c.CreateTime,
		}
	case CaptchaTypeChangePasswd:
		event = RegisterCaptchaEvent{
			Code:       c.Code,
			Email:      c.Email,
			IpAddr:     c.IpAddr,
			Duration:   c.Duration,
			CreateTime: c.CreateTime,
		}
	default:
		return ErrInValidCaptchaType
	}

	return queue.EnQueue(event)
}
