package v1

import (
	"account-management/internal/service"
	"context"
	"log"
	"net/http"
	"strings"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			log.Print("AuthMiddleware.UserIdentity: bearerToken")
			http.Error(w, "Invalid auth header", http.StatusUnauthorized)
			return
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			log.Printf("AuthMiddleware.UserIdentity: h.authService.ParseToken: %v", err)
			http.Error(w, "can not parse token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
