package captcha

import "moj/domain/pkg/queue"

type CreateChangePasswdCaptchaCmd struct {
	AccountID string
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

func (h *CreateChangePasswdCaptchaCmdHandler) Handle(queue queue.EventQueue, cmd *CreateChangePasswdCaptchaCmd) error {
	cap, err := NewCaptcha(cmd.AccountID,
		cmd.Email, CaptchaTypeChangePasswd, cmd.IpAddr, cmd.Duration, cmd.Time)
	if err != nil {
		return err
	}
	err = h.repo.Save(cap)
	if err != nil {
		return err
	}
	return cap.sendEmail(queue)
}
