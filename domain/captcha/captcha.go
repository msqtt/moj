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
	Content    string
	Email      string
	Type       CaptchaType
	IpAddr     string
	Duration   int64
	CreateTime int64
	ExpireTime int64
}

func NewCaptcha(
	accountID int,
	email string, captchaType CaptchaType, ipAddr string, duration int64, time int64) (*Captcha, error) {

	if common.IsEmail(email) {
		return nil, ErrInValidEmail
	}

	if !captchaType.IsValid() {
		return nil, ErrInValidCaptchaType
	}

	content := common.RandomStr(6)

	return &Captcha{
		AccountID:  accountID,
		Content:    content,
		Email:      email,
		Type:       captchaType,
		IpAddr:     ipAddr,
		Duration:   duration,
		CreateTime: time,
		ExpireTime: time + duration,
	}, nil
}

func (c *Captcha) IsExpired(time int64) bool {
	return c.ExpireTime >= time
}
