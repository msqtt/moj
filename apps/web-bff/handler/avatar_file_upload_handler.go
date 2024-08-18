package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"moj/apps/web-bff/etc"
	"moj/apps/web-bff/middleware"
	"moj/apps/web-bff/oss"
	"moj/apps/web-bff/pkg"
	"moj/apps/web-bff/token"
	"net/http"
	"path/filepath"
	"time"
)

type AvatarFileHandler struct {
	uploader       oss.Uploader
	conf           *etc.Config
	sessionManager *token.SessionManager
}

type UploadResponse struct {
	Path string `json:"path"`
	Time int64  `json:"time"`
}

func checkLogin(ctx context.Context, smg *token.SessionManager) (uid string, err error) {
	token, err := middleware.GetAuthTokenFromContext(ctx)
	if err != nil {
		err = errors.Join(errors.New("failed to get token from token"), err)
		return
	}
	uid, err = smg.ValidAccessToken(token)
	if err != nil {
		err = errors.Join(errors.New("failed to get uid from token"), err)
		return
	}
	return
}

// ServeHTTP implements http.Handler.
func (f *AvatarFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid, err := checkLogin(r.Context(), f.sessionManager)
	if err != nil {
		slog.Error("upload avatar: failed to check login", "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("user not login"))
	}

	err = r.ParseMultipartForm(f.conf.AvatarFileSizeLimit)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(
			fmt.Sprintf("uploaded file exceeds avatar size limit %dkb.",
				f.conf.AvatarFileSizeLimit)))
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get avatar from form file"))
		return
	}
	defer avatar.Close()

	path := filepath.Join(f.conf.AvatarFilePrefixPath, uid,
		pkg.Int64ToString(time.Now().Unix()))

	avatarURL, err := f.uploader.Upload(path, avatar)
	if err != nil {
		slog.Error("failed to upload avatar", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to upload avatar"))
		return
	}

	resp := &UploadResponse{
		Path: avatarURL,
		Time: time.Now().Unix(),
	}

	jsonText, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonText)
}

func NewAvatarHandler(
	uploader oss.Uploader,
	conf *etc.Config,
	sessionManager *token.SessionManager,
) *AvatarFileHandler {
	return &AvatarFileHandler{
		uploader:       uploader,
		conf:           conf,
		sessionManager: sessionManager,
	}
}
