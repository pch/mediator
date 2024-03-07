package internal

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

type AuthMiddleware struct {
	config *Config
	next   http.Handler
}

func NewAuthMiddleware(config *Config, next http.Handler) *AuthMiddleware {
	return &AuthMiddleware{config, next}
}

func (h *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.config.AuthToken == "" {
		h.next.ServeHTTP(w, r)
		return
	}

	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	if len(authHeader) <= len(prefix) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := authHeader[len(prefix):]

	hashedToken := sha256.Sum256([]byte(token))
	hashedAuthToken := sha256.Sum256([]byte(h.config.AuthToken))

	if subtle.ConstantTimeCompare(hashedToken[:], hashedAuthToken[:]) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.next.ServeHTTP(w, r)
}
