package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddlewareNoTokenConfigured(t *testing.T) {
	called := false
	mw := NewAuthMiddleware(&Config{}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com", nil)
	mw.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("next handler was not called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestAuthMiddlewareRejectsMalformedOrWrongToken(t *testing.T) {
	mw := NewAuthMiddleware(&Config{AuthToken: "secret-token"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	cases := []string{
		"",
		"Basic abc",
		"Bearer",
		"Bearer wrong",
		"Bearer wrong token",
	}

	for _, auth := range cases {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("Authorization", auth)
		mw.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("Authorization %q: status = %d, want %d", auth, rr.Code, http.StatusUnauthorized)
		}
	}
}

func TestAuthMiddlewareAcceptsValidBearerCaseInsensitive(t *testing.T) {
	called := false
	mw := NewAuthMiddleware(&Config{AuthToken: "secret-token"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "bearer secret-token")
	mw.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("next handler was not called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}
