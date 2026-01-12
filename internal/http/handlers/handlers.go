// Package handlers hosts HTTP handlers and shared helpers.
package handlers

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handlers wires HTTP handlers to the service layer.
type Handlers struct {
	service *service.Service
}

// New constructs a Handlers instance.
func New(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// ErrorResponse is the standard error payload.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Code   apperr.Code       `json:"code"`
	Fields map[string]string `json:"fields,omitempty"`
}

func statusFromCode(code apperr.Code) int {
	switch code {
	case apperr.CodeValidation:
		return http.StatusBadRequest // or 422
	case apperr.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperr.CodeForbidden:
		return http.StatusForbidden
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeConflict:
		return http.StatusConflict
	case apperr.CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// WriteError serializes an application error to the response.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	// If the client disconnected/canceled the request, usually don't write response.
	if errors.Is(err, context.Canceled) {
		return
	}

	code := apperr.CodeInternal
	msg := "internal error"
	var fields map[string]string

	if ae, ok := apperr.As(err); ok {
		code = ae.Code()
		msg = ae.Msg()
		fields = ae.Fields()
	}

	status := statusFromCode(code)

	// Request ID (chi middleware.RequestID sets X-Request-ID)
	reqID := r.Header.Get("X-Request-ID")

	// One log per failed request (centralized)
	slog.Error("request failed",
		"ts", time.Now().Format(time.RFC3339Nano),
		"req_id", reqID,
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"code", code,
		"err", err, // full wrapped chain
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:  msg,
		Code:   code,
		Fields: fields,
	})
}

// WriteJSON writes a JSON response with a status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// DecodeJSON decodes a JSON body into dst with strict field checks.
func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // helps catch client typos early

	if err := dec.Decode(dst); err != nil {
		return apperr.Validation("invalid json body", map[string]string{
			"body": "invalid json",
		})
	}

	// optional: reject trailing garbage after JSON object
	if dec.More() {
		return apperr.Validation("invalid json body", map[string]string{
			"body": "unexpected trailing data",
		})
	}

	return nil
}

func (h *Handlers) readInt64URLParam(r *http.Request, name string) (int64, error) {
	raw := strings.TrimSpace(chi.URLParam(r, name))

	if raw == "" {
		return 0, apperr.Validation("missing path parameter", map[string]string{
			name: "required",
		})
	}

	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperr.Validation("invalid path parameter", map[string]string{
			name: "must be an integer",
		})
	}
	return v, nil
}

func (h *Handlers) readInt64QueryParam(r *http.Request, name string) (int64, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		// keep consistent with your existing error style
		return 0, apperr.Validation("query parameters required", map[string]string{
			name: "query parameter is required",
		})
	}

	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperr.Validation("invalid query parameter", map[string]string{
			name: "query parameters must be integer",
		})
	}

	return v, nil
}

func (h *Handlers) readStringQueryParam(r *http.Request, name string) (string, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(name))
	if raw == "" {
		// keep consistent with your existing error style
		return "", apperr.Validation("query parameters required", map[string]string{
			name: "query parameter is required",
		})
	}

	return raw, nil
}
