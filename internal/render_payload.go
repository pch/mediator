package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type RenderPayload struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func DecodePayloadFromBase64(payloadBase64 string) (*RenderPayload, error) {
	if payloadBase64 == "" {
		return nil, fmt.Errorf("missing payload")
	}

	decodedPayload, err := base64.URLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 payload: %w", err)
	}

	var payload RenderPayload
	if err := json.Unmarshal(decodedPayload, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (p *RenderPayload) Validate() error {
	if p.URL == "" {
		return fmt.Errorf("missing required field: url")
	}
	return nil
}
