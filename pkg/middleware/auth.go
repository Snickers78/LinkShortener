package middleware

import (
	"GoAdvanced/configs"
	"GoAdvanced/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextKeyEmail key = "email"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Token := r.Header.Get("Authorization")
		if !strings.HasPrefix(Token, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		updatedToken := strings.TrimPrefix(Token, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(updatedToken)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyEmail, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
