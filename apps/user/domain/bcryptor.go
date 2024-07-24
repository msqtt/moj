package domain

import (
	"errors"
	inter_error "moj/apps/user/pkg/app_err"
	"moj/domain/pkg/crypt"

	"golang.org/x/crypto/bcrypt"
)

type BCryptor struct {
}

// Encrypt implements crypt.Cryptor.
func (b *BCryptor) Encrypt(raw string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal, err)
	}
	return string(s), err
}

// Valid implements crypt.Cryptor.
func (b *BCryptor) Valid(raw string, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal, err)
	}
	return err
}

func NewBCryptor() crypt.Cryptor {
	return &BCryptor{}
}

var _ crypt.Cryptor = (*BCryptor)(nil)
