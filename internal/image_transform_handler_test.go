package internal

import (
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
)

func TestGenerateImageETagIncludesAllOutputAffectingOptions(t *testing.T) {
	base := &ImageOptions{
		Operations:      []string{"fit"},
		Width:           100,
		Height:          50,
		Quality:         80,
		StripMetadata:   true,
		Format:          vips.ImageTypeJPEG,
		RequestedFormat: "jpeg",
		PixelateFactor:  20,
		Page:            1,
	}

	etagBase := generateImageETag("https://cdn.example.com/file.jpg", base)

	withPage := *base
	withPage.Page = 2
	if generateImageETag("https://cdn.example.com/file.jpg", &withPage) == etagBase {
		t.Fatalf("etag should change when page changes")
	}

	withPixelate := *base
	withPixelate.PixelateFactor = 99
	if generateImageETag("https://cdn.example.com/file.jpg", &withPixelate) == etagBase {
		t.Fatalf("etag should change when pixelate factor changes")
	}

	withRequestedFormat := *base
	withRequestedFormat.RequestedFormat = "auto"
	if generateImageETag("https://cdn.example.com/file.jpg", &withRequestedFormat) == etagBase {
		t.Fatalf("etag should change when requested format changes")
	}
}
