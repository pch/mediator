package internal

import (
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
)

func TestImageType(t *testing.T) {
	cases := map[string]vips.ImageType{
		"jpeg": vips.ImageTypeJPEG,
		"JPEG": vips.ImageTypeJPEG,
		"png":  vips.ImageTypePNG,
		"webp": vips.ImageTypeWEBP,
		"gif":  vips.ImageTypeGIF,
		"avif": vips.ImageTypeAVIF,
		"heif": vips.ImageTypeHEIF,
		"heic": vips.ImageTypeHEIF,
		"pdf":  vips.ImageTypePDF,
	}

	for input, want := range cases {
		if got := ImageType(input); got != want {
			t.Fatalf("ImageType(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestImageTypeFromMimeTypeAndAccept(t *testing.T) {
	if got := ImageTypeFromMimeType("image/jpeg; charset=binary"); got != vips.ImageTypeJPEG {
		t.Fatalf("ImageTypeFromMimeType() = %v", got)
	}
	if got := ImageTypeFromMimeType("application/octet-stream"); got != vips.ImageTypeUnknown {
		t.Fatalf("ImageTypeFromMimeType(unknown) = %v", got)
	}

	if got := ImageTypeFromMimeType("image/heic"); got != vips.ImageTypeHEIF {
		t.Fatalf("ImageTypeFromMimeType(heic) = %v", got)
	}
	if got := ImageTypeFromMimeType("image/avif"); got != vips.ImageTypeAVIF {
		t.Fatalf("ImageTypeFromMimeType(avif) = %v", got)
	}

	if got := ImageTypeFromAccept("image/avif,image/webp,image/jpeg"); got != vips.ImageTypeAVIF {
		t.Fatalf("ImageTypeFromAccept() = %v, want AVIF", got)
	}
	if got := ImageTypeFromAccept("text/html"); got != vips.ImageTypeUnknown {
		t.Fatalf("ImageTypeFromAccept(unknown) = %v", got)
	}
}

func TestMimeTypeFromImageType(t *testing.T) {
	if got := MimeTypeFromImageType(vips.ImageTypePNG); got != "image/png" {
		t.Fatalf("MimeTypeFromImageType(PNG) = %q", got)
	}
	if got := MimeTypeFromImageType(vips.ImageTypeAVIF); got != "image/avif" {
		t.Fatalf("MimeTypeFromImageType(AVIF) = %q", got)
	}
	if got := MimeTypeFromImageType(vips.ImageTypeHEIF); got != "image/heif" {
		t.Fatalf("MimeTypeFromImageType(HEIF) = %q", got)
	}
	if got := MimeTypeFromImageType(vips.ImageTypeWEBP); got != "image/webp" {
		t.Fatalf("MimeTypeFromImageType(WEBP) = %q", got)
	}
	if got := MimeTypeFromImageType(vips.ImageTypeGIF); got != "image/gif" {
		t.Fatalf("MimeTypeFromImageType(GIF) = %q", got)
	}
	if got := MimeTypeFromImageType(vips.ImageTypeUnknown); got != "image/jpeg" {
		t.Fatalf("MimeTypeFromImageType(default) = %q", got)
	}
}
