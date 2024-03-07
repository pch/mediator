package internal

import (
	"context"
	"fmt"
	"net/http"
)

type ImageSource struct {
	Source string
	Path   string
	URL    string
}

func (imageSource *ImageSource) URLWithQueryString(r *http.Request) string {
	if r.URL.RawQuery == "" {
		return imageSource.URL
	}
	return imageSource.URL + "?" + r.URL.RawQuery
}

func NewImageSourceFromHttpRequest(r *http.Request, config *Config) (*ImageSource, error) {
	source := r.PathValue("source")
	path := r.PathValue("path")

	baseURL, exists := config.Sources[source]
	if !exists {
		return nil, fmt.Errorf("source not found: %s", source)
	}
	url := baseURL + "/" + escapeURLPath(path)

	return &ImageSource{
		Source: source,
		Path:   path,
		URL:    url,
	}, nil
}

type contextKey string

func setImageSource(ctx context.Context, imageSource *ImageSource) context.Context {
	key := contextKey("imageSource")
	return context.WithValue(ctx, key, imageSource)
}

func getImageSource(ctx context.Context) *ImageSource {
	key := contextKey("imageSource")
	return ctx.Value(key).(*ImageSource)
}
