package captcha

type CaptchaRepo interface {
	findLatestCaptcha(accountID int, content string, captchaType CaptchaType) (*Captcha, error)
	save(captcha *Captcha) error
}
