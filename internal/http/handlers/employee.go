package handlers

import (
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/http/middleware"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetEmployee(w http.ResponseWriter, r *http.Request) {
	uniqueID, err := h.readStringQueryParam(r, "uniqueID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	langCode := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeByUniqueID(r.Context(), uniqueID, langCode)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) GetEmployeeProfilePicture(w http.ResponseWriter, r *http.Request) {
	uid := strings.TrimSpace(chi.URLParam(r, "uid"))
	if uid == "" {
		WriteError(w, r, apperr.Validation("missing uid", map[string]string{
			"uid": "required",
		}))
		return
	}

	// Basic path traversal guard; uid must be a simple filename.
	if filepath.Base(uid) != uid || strings.Contains(uid, "..") {
		WriteError(w, r, apperr.Validation("invalid uid", map[string]string{
			"uid": "invalid",
		}))
		return
	}

	baseDir := filepath.Join("storage", "employee", "profile_picture")
	filePath := filepath.Join(baseDir, uid)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			WriteError(w, r, apperr.NotFound("profile picture not found"))
			return
		}
		WriteError(w, r, apperr.Internal("internal error", err))
		return
	}
	if info.IsDir() {
		WriteError(w, r, apperr.NotFound("profile picture not found"))
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *Handlers) UpdateEmployeeProfilePicture(w http.ResponseWriter, r *http.Request) {
	uid := strings.TrimSpace(chi.URLParam(r, "uid"))
	if uid == "" {
		WriteError(w, r, apperr.Validation("missing uid", map[string]string{
			"uid": "required",
		}))
		return
	}

	if filepath.Base(uid) != uid || strings.Contains(uid, "..") {
		WriteError(w, r, apperr.Validation("invalid uid", map[string]string{
			"uid": "invalid",
		}))
		return
	}

	const maxUploadSize = 10 << 20 // 10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		WriteError(w, r, apperr.Validation("invalid upload", map[string]string{
			"file": "invalid or too large",
		}))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		WriteError(w, r, apperr.Validation("missing file", map[string]string{
			"file": "required",
		}))
		return
	}
	defer file.Close()

	baseDir := filepath.Join("storage", "employee", "profile_picture")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		WriteError(w, r, apperr.Internal("internal error", err))
		return
	}

	pattern := filepath.Join(baseDir, uid) + ".*"
	matches, _ := filepath.Glob(pattern)
	for _, path := range matches {
		_ = os.Remove(path)
	}
	_ = os.Remove(filepath.Join(baseDir, uid))

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}
	filePath := filepath.Join(baseDir, uid+ext)

	dst, err := os.Create(filePath)
	if err != nil {
		WriteError(w, r, apperr.Internal("internal error", err))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		WriteError(w, r, apperr.Internal("internal error", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
