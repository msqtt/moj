package captcha

import "github.com/msqtt/moj/domain/pkg/queue"

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

func (h *CreateRegisterCaptchaCmdHandler) Handle(queue queue.EventQueue, cmd *CreateRegisterCaptchaCmd) error {
	cap, err := NewCaptcha("", cmd.Email, CaptchaTypeRegister, cmd.IpAddr, cmd.Duration, cmd.Time)
	if err != nil {
		return err
	}
	err = h.repo.Save(cap)
	if err != nil {
		return err
	}
	return cap.sendEmail(queue)
}
