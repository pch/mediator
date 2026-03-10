package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
)

func makePNG(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode(): %v", err)
	}

	return buf.Bytes()
}

func decodeImageSize(t *testing.T, data []byte) (int, int) {
	t.Helper()

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("image.DecodeConfig(): %v", err)
	}

	return cfg.Width, cfg.Height
}

func TestTransformImageFitWidthOnly(t *testing.T) {
	src := makePNG(t, 40, 20)
	opts := &ImageOptions{
		Operations:      []string{"fit"},
		Width:           10,
		Height:          0,
		Quality:         80,
		StripMetadata:   true,
		Format:          vips.ImageTypePNG,
		RequestedFormat: "png",
		AutoRotate:      true,
		PixelateFactor:  20,
		Page:            1,
	}

	out, err := TransformImage(src, opts)
	if err != nil {
		t.Fatalf("TransformImage() error: %v", err)
	}

	w, h := decodeImageSize(t, out.Bytes)
	if w != 10 || h != 5 {
		t.Fatalf("result size = %dx%d, want 10x5", w, h)
	}
}

func TestTransformImageFitHeightOnly(t *testing.T) {
	src := makePNG(t, 40, 20)
	opts := &ImageOptions{
		Operations:      []string{"fit"},
		Width:           0,
		Height:          10,
		Quality:         80,
		StripMetadata:   true,
		Format:          vips.ImageTypePNG,
		RequestedFormat: "png",
		AutoRotate:      true,
		PixelateFactor:  20,
		Page:            1,
	}

	out, err := TransformImage(src, opts)
	if err != nil {
		t.Fatalf("TransformImage() error: %v", err)
	}

	w, h := decodeImageSize(t, out.Bytes)
	if w != 20 || h != 10 {
		t.Fatalf("result size = %dx%d, want 20x10", w, h)
	}
}

func TestTransformImageUnknownOperation(t *testing.T) {
	src := makePNG(t, 10, 10)
	opts := &ImageOptions{Operations: []string{"unknown"}, Format: vips.ImageTypePNG, RequestedFormat: "png"}

	if _, err := TransformImage(src, opts); err == nil {
		t.Fatalf("expected unsupported operation error")
	}
}

func TestImportPageForImageType(t *testing.T) {
	if got := importPageForImageType(1, vips.ImageTypeAVIF); got != 0 {
		t.Fatalf("AVIF page 1 should map to 0, got %d", got)
	}
	if got := importPageForImageType(2, vips.ImageTypeHEIF); got != 1 {
		t.Fatalf("HEIF page 2 should map to 1, got %d", got)
	}
	if got := importPageForImageType(0, vips.ImageTypeHEIF); got != 0 {
		t.Fatalf("HEIF page 0 should remain 0, got %d", got)
	}
	if got := importPageForImageType(1, vips.ImageTypeGIF); got != 1 {
		t.Fatalf("GIF page should not be remapped, got %d", got)
	}
}
