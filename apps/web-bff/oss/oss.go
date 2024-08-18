package oss

import (
	"io"
)

type Uploader interface {
	Upload(prefixPath string, file io.Reader) (url string, err error)
}
