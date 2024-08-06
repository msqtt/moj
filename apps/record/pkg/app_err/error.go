package app_err

import "errors"

var (
	ErrModelNotFound  = errors.New("model not found")
	ErrServerInternal = errors.New("server internal error")
)
