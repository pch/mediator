package internal

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func mustEncodeRenderPayload(t *testing.T, payload RenderPayload) string {
	t.Helper()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	return base64.URLEncoding.EncodeToString(payloadJSON)
}
