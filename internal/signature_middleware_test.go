package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComputeHmacAndMatch(t *testing.T) {
	sig := computeHmac("/image/transform/source/file.jpg?w=10", "my-secret")
	if len(sig) != 64 {
		t.Fatalf("unexpected signature length: %d", len(sig))
	}
	if !signaturesMatch(sig, sig) {
		t.Fatalf("expected signatures to match")
	}
	if signaturesMatch("deadbeef", sig) {
		t.Fatalf("expected signatures not to match")
	}
}

func TestSignatureMiddlewareNoSecretSkipsValidation(t *testing.T) {
	called := false
	mw := NewSignatureMiddleware("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.URL.Query().Get(SignatureParam); got != "" {
			t.Fatalf("signature param should be removed before next handler, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/render/pdf/payload?a=1&s=abc", nil)
	mw.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("next handler was not called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestSignatureMiddlewareRejectsInvalidSignature(t *testing.T) {
	mw := NewSignatureMiddleware("my-secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/image/transform/images/cat.jpg?w=200&s=bad", nil)
	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestSignatureMiddlewareAcceptsValidSignature(t *testing.T) {
	const secret = "my-secret"
	unsignedReq := httptest.NewRequest("GET", "http://example.com/image/transform/images/cat.jpg?w=200", nil)
	sig := computeHmac(unsignedReq.URL.String(), secret)

	called := false
	mw := NewSignatureMiddleware(secret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.URL.Query().Get(SignatureParam); got != "" {
			t.Fatalf("signature param should not be forwarded, got %q", got)
		}
		if got := r.URL.Query().Get("w"); got != "200" {
			t.Fatalf("expected query param w=200, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/image/transform/images/cat.jpg?w=200&s="+sig, nil)
	mw.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("next handler was not called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}
