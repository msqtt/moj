package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type authTokenKey string

const authTokenKey1 authTokenKey = "auth_token"

var ErrTokenInvalid = errors.New("invalid auth token")

func WithAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			r = r.WithContext(context.WithValue(r.Context(), authTokenKey1, token))
		}
		next.ServeHTTP(w, r)
	})
}

func GetAuthTokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(authTokenKey1).(string)
	if !ok {
		return "", errors.New("auth token not found")
	}
	if !strings.Contains(token, "Bearer ") {
		return "", ErrTokenInvalid
	}
	tmp := strings.Split(token, "Bearer ")
	if len(tmp) < 2 {
		return "", ErrTokenInvalid
	}
	token = tmp[1]
	return token, nil
}
