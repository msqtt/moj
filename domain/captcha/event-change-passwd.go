package captcha

type ChangePasswdCaptchaEvent struct {
	Code       string
	Email      string
	IpAddr     string
	Duration   int64
	CreateTime int64
}
