package internal

import (
	"log/slog"
	"net/http"
)

type ImageSourceMiddleware struct {
	config *Config
	next   http.Handler
}

func NewImageSourceMiddleware(config *Config, next http.Handler) *ImageSourceMiddleware {
	return &ImageSourceMiddleware{config, next}
}

func (h *ImageSourceMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageSource, err := NewImageSourceFromHttpRequest(r, h.config)

	if err != nil {
		slog.Error("Image source not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := setImageSource(r.Context(), imageSource)
	newReq := r.WithContext(ctx)
	h.next.ServeHTTP(w, newReq)
}
