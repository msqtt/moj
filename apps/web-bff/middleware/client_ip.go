package middleware

import (
	"context"
	"errors"
	"net/http"
)

type clientIpKey string

const clientIpKey1 clientIpKey = "client_ip"

func WithClientIp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp := r.Header.Get("X-Forwarded-For")
		if clientIp == "" {
			clientIp = r.RemoteAddr
		}
		ctx := context.WithValue(r.Context(), clientIpKey1, clientIp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClientIpFromContext(ctx context.Context) (string, error) {
	clientIp, ok := ctx.Value(clientIpKey1).(string)
	if !ok {
		return "", errors.New("client ip not found in context")
	}
	return clientIp, nil
}
