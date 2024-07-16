package captcha

type CreateChangePasswdCaptchaCmd struct {
	AccountID int
	Email     string
	Time      int64
	IpAddr    string
	Duration  int64
}

type CreateChangePasswdCaptchaCmdHandler struct {
	repo CaptchaRepository
}

func NewCreateChangePasswdCaptchaCmdHandler(repo CaptchaRepository) *CreateChangePasswdCaptchaCmdHandler {
	return &CreateChangePasswdCaptchaCmdHandler{
		repo: repo,
	}
}

func (h *CreateChangePasswdCaptchaCmdHandler) Handle(cmd *CreateChangePasswdCaptchaCmd) error {
	cap, err := NewCaptcha(cmd.AccountID,
		cmd.Email, CaptchaTypeChangePasswd, cmd.IpAddr, cmd.Duration, cmd.Time)
	if err != nil {
		return err
	}
	return h.repo.Save(cap)
}
