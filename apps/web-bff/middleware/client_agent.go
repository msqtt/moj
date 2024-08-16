package middleware

import (
	"context"
	"errors"
	"net/http"
)

type clientAgentKey string

const clientAgentKey1 clientAgentKey = "client_agent"

func WithClientAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agent := r.Header.Get("User-Agent")
		r = r.WithContext(context.WithValue(r.Context(), clientAgentKey1, agent))
		next.ServeHTTP(w, r)
	})
}

func GetClientAgentFromContext(ctx context.Context) (string, error) {
	agent, ok := ctx.Value(clientAgentKey1).(string)
	if !ok {
		return "", errors.New("client agent not found in context")
	}
	return agent, nil
}
