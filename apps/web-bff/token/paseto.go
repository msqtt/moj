package token

import (
	"encoding/json"
	"errors"
	"moj/web-bff/etc"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoTokener struct {
	conf *etc.Config
}

const footerMessage = "moj's foot"

// Encrypt implements Tokener.
func (p *PasetoTokener) Encrypt(plain map[string]string,
	expireDur time.Duration) (token string, err error) {
	now := time.Now()
	jsonToken := paseto.JSONToken{
		Audience:   "any",
		Issuer:     "moj",
		Jti:        "moj identifier",
		Subject:    plain["sub"],
		Expiration: now.Add(expireDur),
		IssuedAt:   now,
		NotBefore:  now,
	}

	for k, v := range plain {
		jsonToken.Set(k, v)
	}

	return paseto.NewV2().Encrypt([]byte(p.conf.SymmetricKey), jsonToken, footerMessage)

}

// Decrypt implements Tokener.
func (p *PasetoTokener) Decrypt(token string) (plain map[string]string, err error) {
	var jsonToken paseto.JSONToken
	var footer string
	err = paseto.NewV2().Decrypt(token, []byte(p.conf.SymmetricKey), &jsonToken, &footer)
	if err != nil {
		return nil, err
	}
	if footer != footerMessage {
		return nil, errors.New("invalid footer")
	}
	mByte, err := jsonToken.MarshalJSON()
	if err != nil {
		return nil, errors.New("cannot marshal payload")
	}
	plain = make(map[string]string)
	err = json.Unmarshal(mByte, &plain)
	if err != nil {
		return nil, errors.New("cannot unmarshal payload to map")
	}
	return
}

func NewPasetoTokener(conf *etc.Config) Tokener {
	return &PasetoTokener{
		conf: conf,
	}
}
