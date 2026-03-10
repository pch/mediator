package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCheckFileSize(t *testing.T) {
	if _, err := checkFileSize(10, ""); err != nil {
		t.Fatalf("empty content length should not error: %v", err)
	}
	if _, err := checkFileSize(10, "invalid"); err != nil {
		t.Fatalf("invalid content length should not error: %v", err)
	}
	if _, err := checkFileSize(10, "11"); err == nil {
		t.Fatalf("expected oversize error")
	}
}

func TestCopyWithSizeLimit(t *testing.T) {
	out := &bytes.Buffer{}
	in := strings.NewReader("12345")

	n, err := copyWithSizeLimit(out, in, 5)
	if err != nil {
		t.Fatalf("copyWithSizeLimit() error: %v", err)
	}
	if n != 5 || out.String() != "12345" {
		t.Fatalf("unexpected copy result: n=%d out=%q", n, out.String())
	}

	out.Reset()
	in = strings.NewReader("123456")
	if _, err := copyWithSizeLimit(out, in, 5); err == nil {
		t.Fatalf("expected oversize error")
	}
}

func TestDownloadFileRejectsOversizeWithoutContentLength(t *testing.T) {
	upload := bytes.Repeat([]byte("x"), 2048)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		_, _ = w.Write(upload)
	}))
	defer srv.Close()

	_, err := DownloadFile(srv.URL, 1024, 2*time.Second)
	if err == nil {
		t.Fatalf("expected oversize error")
	}
	if !strings.Contains(err.Error(), "file too big") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDownloadFileSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	file, err := DownloadFile(srv.URL, 1024, 2*time.Second)
	if err != nil {
		t.Fatalf("DownloadFile() error: %v", err)
	}
	if file.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", file.StatusCode)
	}
	if file.Size() != 5 {
		t.Fatalf("size = %d, want 5", file.Size())
	}
}

func TestProxyFileNon200MarksResponseNoStore(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("missing"))
	}))
	defer upstream.Close()

	req := httptest.NewRequest("GET", "http://origin/render/pdf/payload", nil)
	rr := httptest.NewRecorder()

	// Mimic headers set by render handler before proxying.
	rr.Header().Set("ETag", "\"etag\"")
	rr.Header().Set("Cache-Control", "public, max-age=31536000")

	_, err := ProxyFile(upstream.URL, 1024, 2*time.Second, req, rr)
	if err != nil {
		t.Fatalf("ProxyFile() error: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	if got := rr.Header().Get("ETag"); got != "" {
		t.Fatalf("ETag should be cleared for non-200, got %q", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q", got)
	}
}

func TestProxyFileForwardsRequestHeaders(t *testing.T) {
	var authHeader string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		_, _ = fmt.Fprint(w, "ok")
	}))
	defer upstream.Close()

	req := httptest.NewRequest("GET", "http://origin/render/pdf/payload", nil)
	req.Header.Set("Authorization", "Bearer abc")
	rr := httptest.NewRecorder()

	_, err := ProxyFile(upstream.URL, 1024, 2*time.Second, req, rr)
	if err != nil {
		t.Fatalf("ProxyFile() error: %v", err)
	}

	if authHeader != "Bearer abc" {
		t.Fatalf("expected Authorization header to be forwarded, got %q", authHeader)
	}
}
