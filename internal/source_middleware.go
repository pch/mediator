package internal

import (
	"log/slog"
	"net/http"
)

type SourceMiddleware struct {
	config *Config
	next   http.Handler
}

func NewSourceMiddleware(config *Config, next http.Handler) *SourceMiddleware {
	return &SourceMiddleware{config, next}
}

func (h *SourceMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mediaSource, err := NewMediaSourceFromHttpRequest(r, h.config)

	if err != nil {
		slog.Error("Source not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := setMediaSource(r.Context(), mediaSource)
	newReq := r.WithContext(ctx)
	h.next.ServeHTTP(w, newReq)
}
