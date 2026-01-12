// Package middleware provides HTTP middlewares for request processing.
package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const ctxKeyLang ctxKey = "lang"

// LangFromContext reads the language code from context.
func LangFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyLang).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

// MustLang returns the language code or a default.
func MustLang(ctx context.Context) string {
	if v, ok := LangFromContext(ctx); ok {
		return v
	}
	return "ru"
}

// LanguageMiddleware reads Accept-Language and stores it in context.
func LanguageMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := strings.ToLower(strings.TrimSpace(r.Header.Get("Accept-Language")))

			if !isAllowedLang(lang) {
				lang = "tg"
			}

			ctx := context.WithValue(r.Context(), ctxKeyLang, lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isAllowedLang(lang string) bool {
	switch lang {
	case "tg", "ru", "en":
		return true
	default:
		return false
	}
}
