package oss

import (
	"io"
	"path/filepath"

	"golang.org/x/exp/slog"
)

type MinioOssUploader struct {
}

// Upload implements Uploader.
func (m *MinioOssUploader) Upload(prefixPath string, file io.Reader) (url string, err error) {
	slog.Warn("MinioOssUploader unimplemented")
	path := filepath.Join(prefixPath, "avatar.jpg")
	return path, nil
}

func NewMinioOssUploader() Uploader {
	return &MinioOssUploader{}
}
