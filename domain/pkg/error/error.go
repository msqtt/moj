package domain_err

import "errors"

var (
	ErrInValided  = errors.New("invalided argument")
	ErrDuplicated = errors.New("duplicated operation")
	ErrExpired    = errors.New("expired operation")
)
