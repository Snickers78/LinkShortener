package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Token := r.Header.Get("Authorization")
		updatedToken := strings.SplitAfterN(Token, " ", 2)
		fmt.Println(updatedToken[1])
		next.ServeHTTP(w, r)
	})
}
