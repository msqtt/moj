package account

import (
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
		require.True(t, isURL(url), "URL: %s", url)
	}

	invalidURLs := []string{
		"example.com",
		"http://",
		"https://",
		"https://",
		"ftp://example.com",
	}

	for _, url := range invalidURLs {
		require.False(t, isURL(url), "URL: %s", url)
	}
}

func TestNickName(t *testing.T) {
	validNames := []string{"user_123", "hello", "张三", "李四王五",
		"铁臂阿童木"}
	for _, name := range validNames {
		require.True(t, isNickName(name), "Name: %s", name)
	}

	invalidNames := []string{"", "12345678901234567890", "user#",
		"user@", "user!", "一二三四五六"}
	for _, name := range invalidNames {
		require.False(t, isNickName(name), "Name: %s", name)
	}
}

func TestIsPasswd(t *testing.T) {
	// 测试有效密码
	validPasswd := "P@ssw0rd"
	require.True(t, isPasswd(validPasswd))

	// 测试密码长度不足
	shortPasswd := "pass"
	require.False(t, isPasswd(shortPasswd))

	// 测试密码长度过长
	longPasswd := "ThisIsAVeryLongPasswordThatExceedsTheLimit"
	require.False(t, isPasswd(longPasswd))

	// 测试缺少大写字母
	noUpperPasswd := "passw0rd"
	require.False(t, isPasswd(noUpperPasswd))

	// 测试缺少小写字母
	noLowerPasswd := "PASSW0RD"
	require.False(t, isPasswd(noLowerPasswd))

	// 测试缺少数字
	noDigitPasswd := "Password"
	require.False(t, isPasswd(noDigitPasswd))

	// 测试缺少特殊字符
	noSpecialPasswd := "Password123"
	require.False(t, isPasswd(noSpecialPasswd))
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
		require.True(t, isEmail(email), "Email: %s", email)
	}

	// 测试无效邮箱
	for _, email := range invalidEmails {
		require.False(t, isEmail(email), "Email: %s", email)
	}
}
