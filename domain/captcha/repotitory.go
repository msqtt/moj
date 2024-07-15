package captcha

type CaptchaRepository interface {
	FindLatestCaptcha(email string, content string, captchaType CaptchaType) (*Captcha, error)
	Save(captcha *Captcha) error
}
