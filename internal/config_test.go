package internal

import (
	"strings"
	"testing"
	"time"
)

func TestNewConfigDefaults(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error: %v", err)
	}

	if cfg.DownloadMaxSize != defaultDownloadMaxSize {
		t.Fatalf("DownloadMaxSize = %d, want %d", cfg.DownloadMaxSize, defaultDownloadMaxSize)
	}
	if cfg.DownloadTimeout != defaultDownloadTimeout {
		t.Fatalf("DownloadTimeout = %v, want %v", cfg.DownloadTimeout, defaultDownloadTimeout)
	}
	if cfg.CacheControl != defaultCacheControl {
		t.Fatalf("CacheControl = %q, want %q", cfg.CacheControl, defaultCacheControl)
	}
	if cfg.HttpPort != defaultHttpPort {
		t.Fatalf("HttpPort = %d, want %d", cfg.HttpPort, defaultHttpPort)
	}
	if len(cfg.Sources) != 0 || len(cfg.Renderers) != 0 {
		t.Fatalf("expected empty Sources/Renderers by default")
	}
}

func TestNewConfigOverridesAndLookups(t *testing.T) {
	t.Setenv("MEDIATOR_DOWNLOAD_MAX_SIZE", "1234")
	t.Setenv("MEDIATOR_DOWNLOAD_TIMEOUT", "42")
	t.Setenv("MEDIATOR_CACHE_CONTROL", "private, max-age=60")
	t.Setenv("MEDIATOR_PATH_PREFIX", "/assets")
	t.Setenv("MEDIATOR_HTTP_PORT", "18080")
	t.Setenv("MEDIATOR_HTTP_IDLE_TIMEOUT", "11")
	t.Setenv("MEDIATOR_HTTP_READ_TIMEOUT", "12")
	t.Setenv("MEDIATOR_HTTP_WRITE_TIMEOUT", "13")
	t.Setenv("MEDIATOR_SOURCES", `[{"name":"images","url":"https://cdn.example.com"}]`)
	t.Setenv("MEDIATOR_RENDERERS", `[{"name":"pdf","url":"https://renderer.example.com?url=%s"}]`)

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error: %v", err)
	}

	if cfg.DownloadMaxSize != 1234 {
		t.Fatalf("DownloadMaxSize = %d, want 1234", cfg.DownloadMaxSize)
	}
	if cfg.DownloadTimeout != 42*time.Second {
		t.Fatalf("DownloadTimeout = %v, want 42s", cfg.DownloadTimeout)
	}
	if cfg.CacheControl != "private, max-age=60" {
		t.Fatalf("CacheControl = %q", cfg.CacheControl)
	}
	if cfg.PathPrefix != "/assets" {
		t.Fatalf("PathPrefix = %q", cfg.PathPrefix)
	}
	if cfg.HttpPort != 18080 {
		t.Fatalf("HttpPort = %d", cfg.HttpPort)
	}

	if url, ok := cfg.FindSourceByName("images"); !ok || url != "https://cdn.example.com" {
		t.Fatalf("FindSourceByName(images) = (%q, %v)", url, ok)
	}
	if url, ok := cfg.FindRendererByName("pdf"); !ok || url != "https://renderer.example.com?url=%s" {
		t.Fatalf("FindRendererByName(pdf) = (%q, %v)", url, ok)
	}
}

func TestNewConfigInvalidSourceJSON(t *testing.T) {
	t.Setenv("MEDIATOR_SOURCES", "not-json")

	_, err := NewConfig()
	if err == nil {
		t.Fatalf("expected NewConfig() error")
	}
	if !strings.Contains(err.Error(), "MEDIATOR_SOURCES") {
		t.Fatalf("unexpected error: %v", err)
	}
}
