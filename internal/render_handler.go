package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type RenderHandler struct {
	config *Config
}

func NewRenderHandler(config *Config) *RenderHandler {
	return &RenderHandler{config}
}

func assembleRendererURL(baseURL string, capturedURL string, queryParams url.Values) (string, error) {
	escapedURL := url.QueryEscape(capturedURL)
	finalURL := fmt.Sprintf(baseURL, escapedURL)

	finalURL, err := mergeRequestQueryParams(finalURL, queryParams)
	if err != nil {
		slog.Error("Failed to merge query parameters", "error", err)
		return "", err
	}

	return finalURL, nil
}

func (h *RenderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	renderer := r.PathValue("renderer")
	payloadBase64 := r.PathValue("payloadBase64")

	payload, err := DecodePayloadFromBase64(payloadBase64)
	if err != nil {
		slog.Error("Failed to decode payload", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("Rendering file", "renderer", renderer, "payload", payload)

	rendererURL, exists := h.config.Renderers[renderer]
	if !exists {
		slog.Error("Renderer not supported", "renderer", renderer)
		http.Error(w, "Renderer not supported", http.StatusBadRequest)
		return
	}

	finalURL, err := assembleRendererURL(rendererURL, payload.URL, r.URL.Query())
	if err != nil {
		slog.Error("Failed to assemble renderer URL", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug("Rendering file", "finalURL", finalURL)

	if payload.Filename != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", payload.Filename))
	}

	_, err = ProxyFile(finalURL, h.config.DownloadMaxSize, h.config.DownloadTimeout, r, w)
	if err != nil {
		http.Error(w, "Error when downloading the file", http.StatusInternalServerError)
		return
	}
}
