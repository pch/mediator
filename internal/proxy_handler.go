package internal

import (
	"log/slog"
	"net/http"
)

type ProxyHandler struct {
	config *Config
}

func NewProxyHandler(config *Config) *ProxyHandler {
	return &ProxyHandler{config}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageSource := getImageSource(r.Context())

	slog.Debug("Proxying image", "url", imageSource.URLWithQueryString(r))

	_, err := ProxyFile(imageSource.URLWithQueryString(r), h.config.DownloadMaxSize, h.config.DownloadTimeout, r, w)
	if err != nil {
		http.Error(w, "Error when downloading the file", http.StatusInternalServerError)
	}
}
