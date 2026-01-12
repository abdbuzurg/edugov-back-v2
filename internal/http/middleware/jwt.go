package middleware

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const ctxKeyUserID ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUserID).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				writeUnauthorized(w, "missing authorization header")
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				writeUnauthorized(w, "invalid authorization header")
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if tokenString == "" {
				writeUnauthorized(w, "invalid authorization header")
				return
			}

			signingSecret := viper.GetString("jwt.secret")
			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
				}
				return []byte(signingSecret), nil
			})
			if err != nil || !token.Valid {
				writeUnauthorized(w, "invalid token")
				return
			}

			if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
				writeUnauthorized(w, "token expired")
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUserID, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type errorResponse struct {
	Error  string            `json:"error"`
	Code   apperr.Code       `json:"code"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: msg,
		Code:  apperr.CodeUnauthorized,
	})
}
