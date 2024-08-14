package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNickName(t *testing.T) {
	validNames := []string{"user_123", "hello", "张三三", "李四王五",
		"铁臂阿童木123", "张three", "傻逼2-ai次元"}
	for _, name := range validNames {
		require.True(t, isNickName(name), "Name: %s", name)
	}

	invalidNames := []string{"", "12345678901234567890", "user#",
		"user@", "user!", "一二三四五六七八九十个"}
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
