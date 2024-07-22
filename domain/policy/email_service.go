package policy

type CaptchaEmailCmd struct {
	Email    string
	IpAddr   string
	Time     int64
	Duration int64
}

type EmailService interface {
	SendRegisterEmail(CaptchaEmailCmd) error
	SendChangePassword(CaptchaEmailCmd) error
}
