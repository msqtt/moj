package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"moj/apps/web-bff/etc"
	"moj/apps/web-bff/oss"
	"moj/apps/web-bff/pkg"
	"moj/apps/web-bff/token"
	"net/http"
	"path/filepath"
	"time"
)

type CaseFileHandler struct {
	uploader       oss.Uploader
	conf           *etc.Config
	sessionManager *token.SessionManager
}

// ServeHTTP implements http.Handler.
func (c *CaseFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid, err := checkLogin(r.Context(), c.sessionManager)
	if err != nil {
		slog.Error("upload case file: failed to check login", "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("user not login"))
	}

	err = r.ParseMultipartForm(c.conf.CaseFileSizeLimit)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(
			fmt.Sprintf("uploaded file exceeds case file size limit %dkb.",
				c.conf.CaseFileSizeLimit)))
		return
	}

	file, _, err := r.FormFile("case")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get avatar from form file"))
		return
	}
	defer file.Close()

	path := filepath.Join(c.conf.CaseFilePrefixPath, uid,
		pkg.Int64ToString(time.Now().Unix()))

	caseFileURL, err := c.uploader.Upload(path, file)
	if err != nil {
		slog.Error("failed to upload case file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to upload case file"))
		return
	}

	resp := &UploadResponse{
		Path: caseFileURL,
		Time: time.Now().Unix(),
	}

	jsonText, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonText)
}

func NewCaseFileHandler(uploader oss.Uploader) *CaseFileHandler {
	return &CaseFileHandler{
		uploader: uploader,
	}
}
