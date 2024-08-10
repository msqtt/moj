package captcha

import "context"

type CaptchaRepository interface {
	FindLatestCaptcha(ctx context.Context, email string, code string, captchaType CaptchaType) (*Captcha, error)
	Save(ctx context.Context, captcha *Captcha) error
}
