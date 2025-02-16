package midleware

import (
	"context"
	"fmt"
	"go/avito-test/config"
	"go/avito-test/pkg/jwt"
	"go/avito-test/pkg/response"
	"log"
	"net/http"
	"strings"
)

type key string

const (
	ContextUserNameKey key = "ContextUserNameKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	resp := map[string]string{
		"errors": http.StatusText(http.StatusUnauthorized),
	}
	response.JsonResponse(w, resp, http.StatusUnauthorized)
}

func CheckAuthed(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println("In checked auth", authHeader)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.Split(authHeader, " ")[1] // Remove "Bearer:"
		data, isValid := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		//log.Println(data, isValid)
		r.Context()
		ctx := context.WithValue(r.Context(), ContextUserNameKey, data.Username)
		req := r.WithContext(ctx)
		log.Println(authHeader, "Token:", token)
		next.ServeHTTP(w, req)
	})
}
