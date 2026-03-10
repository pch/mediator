package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRenderHandlerInvalidPayloadIsNoStore(t *testing.T) {
	h := NewRenderHandler(&Config{CacheControl: "public, max-age=31536000"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/render/pdf/notbase64", nil)
	req.SetPathValue("renderer", "pdf")
	req.SetPathValue("payloadBase64", "notbase64")

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q", got)
	}
	if got := rr.Header().Get("ETag"); got != "" {
		t.Fatalf("ETag should be cleared, got %q", got)
	}
}

func TestRenderHandlerProxyFailureIsNoStore(t *testing.T) {
	cfg := &Config{
		DownloadMaxSize: 1024,
		DownloadTimeout: 500 * time.Millisecond,
		Renderers:       []SourceConfig{{Name: "pdf", URL: "http://127.0.0.1:1/render?url=%s"}},
		CacheControl:    "public, max-age=31536000",
	}
	h := NewRenderHandler(cfg)

	payload := mustEncodeRenderPayload(t, RenderPayload{URL: "https://example.com"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/render/pdf/"+payload, nil)
	req.SetPathValue("renderer", "pdf")
	req.SetPathValue("payloadBase64", payload)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q", got)
	}
	if got := rr.Header().Get("ETag"); got != "" {
		t.Fatalf("ETag should be cleared, got %q", got)
	}
}

func TestRenderHandlerSuccessAndNon200ProxyBehavior(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("url") == "missing" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("pdf bytes"))
	}))
	defer upstream.Close()

	cfg := &Config{
		DownloadMaxSize: 1024 * 1024,
		DownloadTimeout: 2 * time.Second,
		Renderers:       []SourceConfig{{Name: "pdf", URL: fmt.Sprintf("%s/render?url=%%s", upstream.URL)}},
		CacheControl:    "public, max-age=3600",
	}
	h := NewRenderHandler(cfg)

	payloadOK := mustEncodeRenderPayload(t, RenderPayload{URL: "ok", Filename: "file.pdf"})
	rrOK := httptest.NewRecorder()
	reqOK := httptest.NewRequest("GET", "http://example.com/render/pdf/"+payloadOK, nil)
	reqOK.SetPathValue("renderer", "pdf")
	reqOK.SetPathValue("payloadBase64", payloadOK)

	h.ServeHTTP(rrOK, reqOK)

	if rrOK.Code != http.StatusOK {
		t.Fatalf("success status = %d, want %d", rrOK.Code, http.StatusOK)
	}
	if got := rrOK.Header().Get("Cache-Control"); got != "public, max-age=3600" {
		t.Fatalf("success Cache-Control = %q", got)
	}
	if got := rrOK.Header().Get("Content-Disposition"); got != `inline; filename="file.pdf"` {
		t.Fatalf("Content-Disposition = %q", got)
	}

	payloadMissing := mustEncodeRenderPayload(t, RenderPayload{URL: "missing"})
	rrMissing := httptest.NewRecorder()
	reqMissing := httptest.NewRequest("GET", "http://example.com/render/pdf/"+payloadMissing, nil)
	reqMissing.SetPathValue("renderer", "pdf")
	reqMissing.SetPathValue("payloadBase64", payloadMissing)

	h.ServeHTTP(rrMissing, reqMissing)

	if rrMissing.Code != http.StatusNotFound {
		t.Fatalf("non-200 status = %d, want %d", rrMissing.Code, http.StatusNotFound)
	}
	if got := rrMissing.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("non-200 Cache-Control = %q", got)
	}
	if got := rrMissing.Header().Get("ETag"); got != "" {
		t.Fatalf("non-200 ETag should be cleared, got %q", got)
	}
}
