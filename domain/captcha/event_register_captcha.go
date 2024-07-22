package captcha

type RegisterCaptchaEvent struct {
	Code       string
	Email      string
	IpAddr     string
	Duration   int64
	CreateTime int64
}
