package captcha

type CreateRegisterCaptchaCmd struct {
	Email    string
	IpAddr   string
	Time     int64
	Duration int64
}

type CreateRegisterCaptchaCmdHandler struct {
	repo CaptchaRepository
}

func NewCreateRegisterCaptchaCmdHandler(repo CaptchaRepository) *CreateRegisterCaptchaCmdHandler {
	return &CreateRegisterCaptchaCmdHandler{
		repo: repo,
	}
}

func (h *CreateRegisterCaptchaCmdHandler) Handle(cmd *CreateRegisterCaptchaCmd) error {
	cap, err := NewCaptcha(0, cmd.Email, CaptchaTypeRegister, cmd.IpAddr, cmd.Duration, cmd.Time)
	if err != nil {
		return err
	}
	return h.repo.Save(cap)
}
