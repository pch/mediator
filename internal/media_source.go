package internal

import (
	"context"
	"fmt"
	"net/http"
)

type MediaSource struct {
	Source string
	Path   string
	URL    string
}

func (mediaSource *MediaSource) URLWithQueryString(r *http.Request) string {
	if r.URL.RawQuery == "" {
		return mediaSource.URL
	}
	return mediaSource.URL + "?" + r.URL.RawQuery
}

func NewMediaSourceFromHttpRequest(r *http.Request, config *Config) (*MediaSource, error) {
	source := r.PathValue("source")
	path := r.PathValue("path")

	baseURL, exists := config.Sources[source]
	if !exists {
		return nil, fmt.Errorf("source not found: %s", source)
	}
	url := baseURL + "/" + escapeURLPath(path)

	return &MediaSource{
		Source: source,
		Path:   path,
		URL:    url,
	}, nil
}

type contextKey string

func setMediaSource(ctx context.Context, mediaSource *MediaSource) context.Context {
	key := contextKey("mediaSource")
	return context.WithValue(ctx, key, mediaSource)
}

func getMediaSource(ctx context.Context) *MediaSource {
	key := contextKey("mediaSource")
	return ctx.Value(key).(*MediaSource)
}
