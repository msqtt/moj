package token

import "time"

type Tokener interface {
	Encrypt(plain map[string]string, expireDur time.Duration) (token string, err error)
	Decrypt(token string) (plain map[string]string, err error)
}
