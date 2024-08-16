package token

import (
	"moj/apps/web-bff/etc"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	// Create a PasetoTokener instance
	p := NewPasetoTokener(&etc.Config{
		AppPort:              0,
		KeyFile:              "",
		CertFile:             "",
		TLS:                  false,
		Debug:                false,
		AccessTokenDuration:  0,
		RefreshTokenDuration: 0,
		SymmetricKey:         "12345678901234567890123456789012",
		UserRPCAddr:          "",
		CaptchaRPCAddr:       "",
		QuestionRPCAddr:      "",
		GameRPCAddr:          "",
		RecordRPCAddr:        "",
	})

	// Test case 1: Normal encryption
	plain1 := map[string]string{
		"sub":  "testSubject1",
		"data": "testData1",
	}
	expireDur1 := time.Hour
	token1, err1 := p.Encrypt(plain1, expireDur1)
	require.NoError(t, err1)
	require.NotEmpty(t, token1)
	t.Log(token1)

	// Test case 2: Empty plain data
	plain2 := map[string]string{}
	expireDur2 := time.Minute
	token2, err2 := p.Encrypt(plain2, expireDur2)
	require.NoError(t, err2)
	require.NotEmpty(t, token2)
	t.Log(token2)

	// Test case 3: Expired token
	plain3 := map[string]string{
		"sub": "testSubject3",
	}
	expireDur3 := -time.Hour
	token3, err3 := p.Encrypt(plain3, expireDur3)
	require.NoError(t, err3)
	require.NotEmpty(t, token3)
	t.Log(token3)
}
