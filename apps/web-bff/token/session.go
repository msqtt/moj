package token

import (
	"errors"
	"moj/web-bff/etc"
)

type SessionManager struct {
	tokener Tokener
	conf    *etc.Config
}

var (
	ErrAccessTokenGenerate  = errors.New("access token generate error")
	ErrRefreshTokenGenerate = errors.New("refresh token generate error")
	ErrTokenInvalid         = errors.New("token invalid")
)

func (s *SessionManager) GenerateAccessToken(accountID string) (token string, err error) {
	m := make(map[string]string)
	m["sub"] = "access token"
	m["account_id"] = accountID
	token, err = s.tokener.Encrypt(m, s.conf.AccessTokenDuration)
	if err != nil {
		err = errors.Join(ErrAccessTokenGenerate, err)
	}
	return
}

func (s *SessionManager) GenerateRefreshToken(accountID string) (token string, err error) {
	m := make(map[string]string)
	m["sub"] = "refresh token"
	m["account_id"] = accountID
	token, err = s.tokener.Encrypt(m, s.conf.RefreshTokenDuration)
	if err != nil {
		err = errors.Join(ErrRefreshTokenGenerate, err)
	}
	return
}

func (s *SessionManager) ValidRefreshToken(token string) (accountID string, err error) {
	m, err := s.tokener.Decrypt(token)
	accountID = m["account_id"]
	sub := m["sub"]
	if err != nil || sub != "refresh token" {
		err = errors.Join(ErrTokenInvalid, err)
	}
	return
}

func (s *SessionManager) ValidAccessToken(token string) (accountID string, err error) {
	m, err := s.tokener.Decrypt(token)
	accountID = m["account_id"]
	sub := m["sub"]
	if err != nil || sub != "access token" {
		err = errors.Join(ErrTokenInvalid, err)
	}
	return
}

func NewSessionManager(conf *etc.Config, tokener Tokener) *SessionManager {
	return &SessionManager{
		tokener: tokener,
		conf:    conf,
	}
}
