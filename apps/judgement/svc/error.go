package svc

import (
	"errors"
	inter_error "moj/apps/judgement/pkg/app_err"
	domain_err "moj/domain/pkg/error"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func responseStatusError(err error) error {
	var code codes.Code
	switch {
	case errors.Is(err, inter_error.ErrServerInternal):
		err = unwrapAndRmErr(err, inter_error.ErrServerInternal)
		code = codes.Internal
	case errors.Is(err, inter_error.ErrModelNotFound):
		err = unwrapAndRmErr(err, inter_error.ErrModelNotFound)
		code = codes.NotFound
	case errors.Is(err, domain_err.ErrInValided):
		err = unwrapAndRmErr(err, domain_err.ErrInValided)
		code = codes.InvalidArgument
	case errors.Is(err, domain_err.ErrExpired):
		err = unwrapAndRmErr(err, domain_err.ErrExpired)
		code = codes.DeadlineExceeded
	case errors.Is(err, domain_err.ErrDuplicated):
		err = unwrapAndRmErr(err, domain_err.ErrDuplicated)
		code = codes.AlreadyExists
	default:
		code = codes.Unknown
	}
	return status.Error(code, err.Error())
}

func unwrapAll(err error) []error {
	var errs []error
	if err == nil {
		return errs
	}

	// 检查错误是否实现了 Unwrap() []error 方法
	type unwrapper interface {
		Unwrap() []error
	}
	if uw, ok := err.(unwrapper); ok {
		for _, e := range uw.Unwrap() {
			errs = append(errs, unwrapAll(e)...)
		}
		return errs
	}

	// 如果错误没有实现 Unwrap() []error 方法，返回单个错误
	return []error{err}
}

func filter(errs []error, check func(error) bool) []error {
	ret := make([]error, 0)
	for _, err := range errs {
		if check(err) {
			ret = append(ret, err)
		}
	}
	return ret
}

func unwrapAndRmErr(src, tgt error) error {
	return errors.Join(
		filter(unwrapAll(src),
			func(err error) bool {
				return !errors.Is(err, tgt)
			})...)
}
