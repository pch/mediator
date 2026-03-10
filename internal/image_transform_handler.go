package internal

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

type ImageTransformHandler struct {
	config *Config
	sem    chan struct{}
}

func NewImageTransformHandler(config *Config) *ImageTransformHandler {
	maxConcurrent := config.MaxConcurrentTransforms
	if maxConcurrent <= 0 {
		maxConcurrent = defaultMaxConcurrentTransforms
	}

	return &ImageTransformHandler{
		config: config,
		sem:    make(chan struct{}, maxConcurrent),
	}
}

func generateImageETag(sourceURL string, imageOptions *ImageOptions) string {
	h := sha1.New()

	io.WriteString(h, sourceURL)
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Width))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Height))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Quality))
	io.WriteString(h, fmt.Sprintf("%v", imageOptions.StripMetadata))
	io.WriteString(h, fmt.Sprintf("%v", imageOptions.Format))
	io.WriteString(h, fmt.Sprintf("%s", imageOptions.RequestedFormat))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.PixelateFactor))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Page))
	io.WriteString(h, strings.Join(imageOptions.Operations, ","))

	return fmt.Sprintf("\"%x\"", h.Sum(nil))
}

func detectDownloadedImageType(downloadedFile *DownloadedFile) vips.ImageType {
	imageType := ImageTypeFromMimeType(downloadedFile.ContentType)
	if imageType != vips.ImageTypeUnknown {
		return imageType
	}

	return vips.DetermineImageType(downloadedFile.Buffer.Bytes())
}

func (h *ImageTransformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageSource := getImageSource(r.Context())
	imageOptions := NewImageOptionsFromRequest(r)

	etag := generateImageETag(imageSource.URL, imageOptions)
	w.Header().Set("ETag", etag)

	if ifNoneMatch := r.Header.Get("If-None-Match"); ifNoneMatch != "" {
		if strings.Contains(ifNoneMatch, etag) {
			slog.Info("ETag match", "etag", etag, "ifNoneMatch", ifNoneMatch)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	select {
	case h.sem <- struct{}{}:
		defer func() { <-h.sem }()
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusServiceUnavailable)
		return
	}

	downloadedFile, err := DownloadFile(imageSource.URL, h.config.DownloadMaxSize, h.config.DownloadTimeout)
	if err != nil {
		slog.Error("Download error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	downloadedImageType := detectDownloadedImageType(downloadedFile)
	if downloadedImageType == vips.ImageTypeUnknown || !vips.IsTypeSupported(downloadedImageType) {
		slog.Error("Unsupported image format: "+downloadedFile.ContentType, "error", err)
		http.Error(w, "Unsupported image format: "+downloadedFile.ContentType, http.StatusUnprocessableEntity)
		return
	}

	if imageOptions.Format == vips.ImageTypeUnknown {
		if IsImageExportSupported(downloadedImageType) {
			imageOptions.Format = downloadedImageType
		} else {
			imageOptions.Format = vips.ImageTypeJPEG
		}
	}

	processedImage, err := TransformImage(downloadedFile.Buffer.Bytes(), imageOptions)
	if err != nil {
		slog.Error("TransformImage error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if imageOptions.RequestedFormat == "auto" {
		w.Header().Set("Vary", "Accept")
	}

	w.Header().Set("Cache-Control", h.config.CacheControl)
	w.Header().Set("Content-Type", processedImage.Mime)
	w.Header().Set("Content-Length", strconv.Itoa(processedImage.Size))
	w.WriteHeader(http.StatusOK)
	w.Write(processedImage.Bytes)
}
