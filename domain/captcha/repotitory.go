package captcha

type CaptchaRepository interface {
	FindLatestCaptcha(email string, code string, captchaType CaptchaType) (*Captcha, error)
	Save(captcha *Captcha) error
}
