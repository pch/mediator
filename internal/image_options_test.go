package internal

import (
	"net/http/httptest"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
)

func TestNewImageOptionsFromRequestDefaults(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/image/transform/images/file.jpg", nil)
	opts := NewImageOptionsFromRequest(req)

	if len(opts.Operations) != 1 || opts.Operations[0] != defaultOperation {
		t.Fatalf("Operations = %#v", opts.Operations)
	}
	if opts.Quality != defaultQuality {
		t.Fatalf("Quality = %d", opts.Quality)
	}
	if !opts.StripMetadata {
		t.Fatalf("StripMetadata = false, want true")
	}
	if opts.PixelateFactor != defaultPixelateFactor {
		t.Fatalf("PixelateFactor = %d", opts.PixelateFactor)
	}
	if opts.Page != defaultPage {
		t.Fatalf("Page = %d", opts.Page)
	}
}

func TestNewImageOptionsFromRequestParsesValues(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/?op=fit,pixelate&w=120&h=80&q=72&strip=false&format=auto&pixelatefactor=17&page=3", nil)
	req.Header.Set("Accept", "image/avif,image/webp,image/jpeg")

	opts := NewImageOptionsFromRequest(req)

	if len(opts.Operations) != 2 || opts.Operations[0] != "fit" || opts.Operations[1] != "pixelate" {
		t.Fatalf("Operations = %#v", opts.Operations)
	}
	if opts.Width != 120 || opts.Height != 80 {
		t.Fatalf("size = %dx%d", opts.Width, opts.Height)
	}
	if opts.Quality != 72 {
		t.Fatalf("Quality = %d", opts.Quality)
	}
	if opts.StripMetadata {
		t.Fatalf("StripMetadata = true, want false")
	}
	if opts.AutoRotate {
		t.Fatalf("AutoRotate should follow StripMetadata")
	}
	if opts.RequestedFormat != "auto" {
		t.Fatalf("RequestedFormat = %q", opts.RequestedFormat)
	}
	if opts.Format != vips.ImageTypeAVIF {
		t.Fatalf("Format = %v, want AVIF", opts.Format)
	}
	if opts.PixelateFactor != 17 || opts.Page != 3 {
		t.Fatalf("pixelate/page = %d/%d", opts.PixelateFactor, opts.Page)
	}
}
