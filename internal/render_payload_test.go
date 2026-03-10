package internal

import (
	"encoding/base64"
	"testing"
)

func TestDecodePayloadFromBase64Valid(t *testing.T) {
	encoded := mustEncodeRenderPayload(t, RenderPayload{URL: "https://example.com", Filename: "invoice.pdf"})

	payload, err := DecodePayloadFromBase64(encoded)
	if err != nil {
		t.Fatalf("DecodePayloadFromBase64() error: %v", err)
	}

	if payload.URL != "https://example.com" || payload.Filename != "invoice.pdf" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestDecodePayloadFromBase64Errors(t *testing.T) {
	if _, err := DecodePayloadFromBase64(""); err == nil {
		t.Fatalf("expected missing payload error")
	}

	if _, err := DecodePayloadFromBase64("%%%not-base64"); err == nil {
		t.Fatalf("expected invalid base64 error")
	}

	invalidJSON := base64.URLEncoding.EncodeToString([]byte(`{"url":`))
	if _, err := DecodePayloadFromBase64(invalidJSON); err == nil {
		t.Fatalf("expected invalid json error")
	}

	missingURL := base64.URLEncoding.EncodeToString([]byte(`{"filename":"x"}`))
	if _, err := DecodePayloadFromBase64(missingURL); err == nil {
		t.Fatalf("expected missing url error")
	}
}
