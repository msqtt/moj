package listener

import (
	"errors"
	domain_err "moj/domain/pkg/error"
	"moj/domain/policy"
)

type EmailListener struct {
	sendCaptchaEmailPolicy *policy.SendCaptchaEmailPolicy
}

// OnEvent implements Listener.
func (e *EmailListener) OnEvent(event any) error {
	err := e.sendCaptchaEmailPolicy.OnEvent(event)
	if errors.Is(err, domain_err.ErrEventTypeInvalid) {
		return nil
	}
	return err
}

func NewEmailListener(policy *policy.SendCaptchaEmailPolicy) *EmailListener {
	return &EmailListener{
		sendCaptchaEmailPolicy: policy,
	}
}
