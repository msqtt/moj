package common_test

import (
	"github.com/msqtt/moj/domain/pkg/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsURL(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"http://example.com/page",
		"https://example.com/page?param=value",
	}

	for _, url := range validURLs {
		require.True(t, common.IsURL(url), "URL: %s", url)
	}

	invalidURLs := []string{
		"example.com",
		"http://",
		"https://",
		"https://",
		"ftp://example.com",
	}

	for _, url := range invalidURLs {
		require.False(t, common.IsURL(url), "URL: %s", url)
	}
}

func TestIsEmail(t *testing.T) {
	// 有效邮箱
	validEmails := []string{
		"example@example.com",
		"test.user@example.org",
	}

	// 无效邮箱
	invalidEmails := []string{
		"invalid.email",
		"@example.com",
		"test@",
	}

	// 测试有效邮箱
	for _, email := range validEmails {
		require.True(t, common.IsEmail(email), "Email: %s", email)
	}

	// 测试无效邮箱
	for _, email := range invalidEmails {
		require.False(t, common.IsEmail(email), "Email: %s", email)
	}
}

func TestSha1(t *testing.T) {
	t.Log(common.Sha1("hello world"))

	require.Equal(t, common.Sha1("hello"), common.Sha1("hello"))
	require.NotEqual(t, common.Sha1("hello"), common.Sha1("Bye"))
}
