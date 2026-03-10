package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

func encodeImageForTest(t *testing.T, imageType vips.ImageType) []byte {
	t.Helper()

	src := makePNG(t, 40, 20)
	img, err := vips.LoadImageFromBuffer(src, nil)
	if err != nil {
		t.Fatalf("LoadImageFromBuffer(): %v", err)
	}
	defer img.Close()

	switch imageType {
	case vips.ImageTypeAVIF:
		params := vips.NewAvifExportParams()
		out, _, err := img.ExportAvif(params)
		if err != nil {
			t.Skipf("AVIF export not available in this environment: %v", err)
		}
		return out
	case vips.ImageTypeHEIF:
		params := vips.NewHeifExportParams()
		out, _, err := img.ExportHeif(params)
		if err != nil {
			t.Skipf("HEIF export not available in this environment: %v", err)
		}
		return out
	default:
		t.Fatalf("unsupported image type for helper: %v", imageType)
		return nil
	}
}

func TestDetectDownloadedImageTypeFallsBackToContentSniffing(t *testing.T) {
	if !vips.IsTypeSupported(vips.ImageTypeAVIF) {
		t.Skip("AVIF not supported by current libvips runtime")
	}

	avifBytes := encodeImageForTest(t, vips.ImageTypeAVIF)

	file := &DownloadedFile{ContentType: "application/octet-stream"}
	file.Buffer.Write(avifBytes)

	if got := detectDownloadedImageType(file); got != vips.ImageTypeAVIF {
		t.Fatalf("detectDownloadedImageType() = %v, want %v", got, vips.ImageTypeAVIF)
	}
}

func TestImageTransformHandlerSupportsAVIFAndHEICInput(t *testing.T) {
	cases := []struct {
		name        string
		contentType string
		imageType   vips.ImageType
	}{
		{name: "avif", contentType: "image/avif", imageType: vips.ImageTypeAVIF},
		{name: "heic", contentType: "image/heic", imageType: vips.ImageTypeHEIF},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !vips.IsTypeSupported(tc.imageType) {
				t.Skipf("%v not supported by current libvips runtime", tc.imageType)
			}

			inputBytes := encodeImageForTest(t, tc.imageType)

			srcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				_, _ = w.Write(inputBytes)
			}))
			defer srcServer.Close()

			h := NewImageTransformHandler(&Config{
				DownloadMaxSize: 1024 * 1024,
				DownloadTimeout: 2 * time.Second,
				CacheControl:    "public, max-age=60",
			})

			req := httptest.NewRequest("GET", "http://example.com/image/transform/src/img?w=10", nil)
			req = req.WithContext(setImageSource(req.Context(), &ImageSource{URL: srcServer.URL + "/img"}))
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d, body=%q", rr.Code, http.StatusOK, rr.Body.String())
			}

			outType := vips.DetermineImageType(rr.Body.Bytes())
			if outType != tc.imageType {
				t.Fatalf("output type = %v, want %v", outType, tc.imageType)
			}
		})
	}
}

func TestImageTransformHandlerSupportsAVIFOutputFormat(t *testing.T) {
	if !vips.IsTypeSupported(vips.ImageTypeAVIF) {
		t.Skip("AVIF not supported by current libvips runtime")
	}

	srcPNG := makePNG(t, 40, 20)
	srcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(srcPNG)
	}))
	defer srcServer.Close()

	h := NewImageTransformHandler(&Config{
		DownloadMaxSize: 1024 * 1024,
		DownloadTimeout: 2 * time.Second,
		CacheControl:    "public, max-age=60",
	})

	req := httptest.NewRequest("GET", "http://example.com/image/transform/src/img?w=10&format=avif", nil)
	req = req.WithContext(setImageSource(req.Context(), &ImageSource{URL: srcServer.URL + "/img"}))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%q", rr.Code, http.StatusOK, rr.Body.String())
	}

	if got := rr.Header().Get("Content-Type"); got != "image/avif" {
		t.Fatalf("Content-Type = %q, want image/avif", got)
	}

	if outType := vips.DetermineImageType(rr.Body.Bytes()); outType != vips.ImageTypeAVIF {
		t.Fatalf("output type = %v, want %v", outType, vips.ImageTypeAVIF)
	}
}
