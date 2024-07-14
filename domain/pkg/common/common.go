package common

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"regexp"
)

// IsEmail checks whether given email is a valid email.
func IsEmail(email string) bool {
	regex := regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
	return regex.MatchString(email)
}

// IsURL checks whether given link is a valid URL.
func IsURL(link string) bool {
	regex := regexp.MustCompile(
		`^(http|https)://[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`)
	return regex.MatchString(link)
}

func RandomStr(n int) (ret string) {
	alphbet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for range n {
		ret += string(alphbet[rand.Intn(len(alphbet))])
	}
	return
}

func Sha1(str string) string {
	hashBytes := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", hashBytes)
}
