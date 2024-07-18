package captcha

import (
	"errors"
	"moj/domain/pkg/common"
)

type CaptchaType int

const (
	CaptchaTypeRegister CaptchaType = iota
	CaptchaTypeChangePasswd
)

func (c CaptchaType) IsValid() bool {
	return c >= CaptchaTypeRegister && c <= CaptchaTypeChangePasswd
}

var (
	ErrInValidEmail       = errors.New("invalid email")
	ErrInValidCaptchaType = errors.New("invalid captcha type")
)

type Captcha struct {
	AccountID  int
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
	accountID int,
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
