package internal

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	next http.Handler
}

func NewLoggingMiddleware(next http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{next}
}

func (h *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newWriter := newResponseWriter(w)

	started := time.Now()
	h.next.ServeHTTP(newWriter, r)
	elapsed := time.Since(started)

	userAgent := r.Header.Get("User-Agent")
	remoteAddr := r.Header.Get("X-Forwarded-For")
	queryString := r.URL.Query().Encode()
	respContent := newWriter.Header().Get("Content-Type")

	slog.Info("Request",
		"method", r.Method,
		"path", r.URL.Path,
		"query", queryString,
		"remote_add", remoteAddr,
		"user_agent", userAgent,
		"resp_content_type", respContent,
		"duration", elapsed,
	)
}
