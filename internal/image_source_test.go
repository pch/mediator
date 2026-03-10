package internal

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestNewImageSourceFromHttpRequest(t *testing.T) {
	cfg := &Config{Sources: []SourceConfig{{Name: "images", URL: "https://cdn.example.com"}}}

	req := httptest.NewRequest("GET", "http://example.com/image/transform/images/folder/file%20name.jpg", nil)
	req.SetPathValue("source", "images")
	req.SetPathValue("path", "folder/file name.jpg")

	source, err := NewImageSourceFromHttpRequest(req, cfg)
	if err != nil {
		t.Fatalf("NewImageSourceFromHttpRequest() error: %v", err)
	}

	if source.Source != "images" {
		t.Fatalf("Source = %q", source.Source)
	}
	if source.Path != "folder/file name.jpg" {
		t.Fatalf("Path = %q", source.Path)
	}
	if source.URL != "https://cdn.example.com/folder/file%20name.jpg" {
		t.Fatalf("URL = %q", source.URL)
	}
}

func TestNewImageSourceFromHttpRequestMissingSource(t *testing.T) {
	cfg := &Config{Sources: []SourceConfig{{Name: "images", URL: "https://cdn.example.com"}}}

	req := httptest.NewRequest("GET", "http://example.com/image/transform/unknown/path", nil)
	req.SetPathValue("source", "unknown")
	req.SetPathValue("path", "path")

	if _, err := NewImageSourceFromHttpRequest(req, cfg); err == nil {
		t.Fatalf("expected source-not-found error")
	}
}

func TestImageSourceURLWithQueryString(t *testing.T) {
	source := &ImageSource{URL: "https://cdn.example.com/path/file.jpg"}
	req := httptest.NewRequest("GET", "http://example.com/?w=100&h=80", nil)

	if got := source.URLWithQueryString(req); got != "https://cdn.example.com/path/file.jpg?w=100&h=80" {
		t.Fatalf("URLWithQueryString = %q", got)
	}
}

func TestImageSourceContextHelpers(t *testing.T) {
	source := &ImageSource{Source: "images", URL: "https://cdn.example.com/x.jpg"}
	ctx := setImageSource(context.Background(), source)

	if got := getImageSource(ctx); got != source {
		t.Fatalf("unexpected image source in context: %+v", got)
	}
}
