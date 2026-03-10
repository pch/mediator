package internal

import (
	"crypto/tls"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestQueryHelpers(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/?name=alice&n=7&flag=true&flag=false", nil)

	if got, ok := getQueryParam("name", r); !ok || got != "alice" {
		t.Fatalf("getQueryParam(name) = (%q, %v)", got, ok)
	}
	if got, ok := getQueryParam("missing", r); ok || got != "" {
		t.Fatalf("getQueryParam(missing) = (%q, %v)", got, ok)
	}

	if got := getQueryParamWithDefault("missing", "default", r); got != "default" {
		t.Fatalf("getQueryParamWithDefault = %q", got)
	}

	if got, ok := getQueryParamInt("n", r); !ok || got != 7 {
		t.Fatalf("getQueryParamInt(n) = (%d, %v)", got, ok)
	}
	if got := getQueryParamIntWithDefault("missing", 99, r); got != 99 {
		t.Fatalf("getQueryParamIntWithDefault = %d", got)
	}

	if got, ok := getQueryParamBool("flag", r); !ok || !got {
		t.Fatalf("getQueryParamBool(flag) = (%v, %v)", got, ok)
	}
	if got := getQueryParamBoolWithDefault("missing", true, r); !got {
		t.Fatalf("getQueryParamBoolWithDefault = %v", got)
	}
}

func TestRemoveParamFromQuery(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/?a=1&s=abc&b=2", nil)
	removeParamFromQuery(r, "s")

	if got := r.URL.Query().Get("s"); got != "" {
		t.Fatalf("signature param not removed: %q", got)
	}
	if got := r.URL.Query().Get("a"); got != "1" {
		t.Fatalf("unexpected a = %q", got)
	}
}

func TestCurrentRequestHost(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/path", nil)
	r.Host = "example.com"
	if got := currentRequestHost(r); got != "http://example.com" {
		t.Fatalf("currentRequestHost(http) = %q", got)
	}

	r.Header.Set("X-Forwarded-Proto", "https")
	if got := currentRequestHost(r); got != "https://example.com" {
		t.Fatalf("currentRequestHost(forwarded https) = %q", got)
	}

	r.TLS = &tls.ConnectionState{}
	if got := currentRequestHost(r); got != "https://example.com" {
		t.Fatalf("currentRequestHost(tls) = %q", got)
	}
}

func TestPathHelpers(t *testing.T) {
	if got := escapeURLPath("a b/c+d"); got != "a%20b/c+d" {
		t.Fatalf("escapeURLPath = %q", got)
	}

	if got := ensureValidPathPrefixFormat("assets/"); got != "/assets" {
		t.Fatalf("ensureValidPathPrefixFormat = %q", got)
	}
	if got := ensureValidPathPrefixFormat(""); got != "" {
		t.Fatalf("ensureValidPathPrefixFormat(empty) = %q", got)
	}
}

func TestMergeRequestQueryParams(t *testing.T) {
	query := url.Values{}
	query.Add("b", "2")
	query.Add("b", "3")
	query.Add("a", "1")

	merged, err := mergeRequestQueryParams("https://example.com/render?x=y", query)
	if err != nil {
		t.Fatalf("mergeRequestQueryParams() error: %v", err)
	}

	parsed, err := url.Parse(merged)
	if err != nil {
		t.Fatalf("url.Parse() error: %v", err)
	}

	if parsed.Query().Get("x") != "y" {
		t.Fatalf("missing existing query param")
	}
	if parsed.Query().Get("a") != "1" {
		t.Fatalf("missing merged query param a")
	}
	if got := parsed.Query()["b"]; len(got) != 2 || got[0] != "2" || got[1] != "3" {
		t.Fatalf("unexpected merged b values: %#v", got)
	}
}
