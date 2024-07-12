package captcha

type CaptchaRepository interface {
	findLatestCaptcha(accountID int, content string, captchaType CaptchaType) (*Captcha, error)
	save(captcha *Captcha) error
}
