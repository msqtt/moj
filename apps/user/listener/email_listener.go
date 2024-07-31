package listener

import "moj/domain/policy"

type EmailListener struct {
	*policy.SendCaptchaEmailPolicy
}

func NewEmailListener(policy *policy.SendCaptchaEmailPolicy) *EmailListener {
	return &EmailListener{
		SendCaptchaEmailPolicy: policy,
	}
}
