package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultRouteHandlerRoot(t *testing.T) {
	h := NewDefaultRouteHandler()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if body["message"] != "hello, world" {
		t.Fatalf("message = %q", body["message"])
	}
}

func TestDefaultRouteHandlerNotFound(t *testing.T) {
	h := NewDefaultRouteHandler()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/not-found", nil)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}
