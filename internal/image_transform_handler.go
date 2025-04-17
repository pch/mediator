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
}

func NewImageTransformHandler(config *Config) *ImageTransformHandler {
	return &ImageTransformHandler{config}
}

func generateImageETag(sourceURL string, imageOptions *ImageOptions) string {
	h := sha1.New()

	io.WriteString(h, sourceURL)
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Width))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Height))
	io.WriteString(h, fmt.Sprintf("%d", imageOptions.Quality))
	io.WriteString(h, fmt.Sprintf("%v", imageOptions.StripMetadata))
	io.WriteString(h, fmt.Sprintf("%v", imageOptions.Format))
	io.WriteString(h, strings.Join(imageOptions.Operations, ","))

	return fmt.Sprintf("\"%x\"", h.Sum(nil))
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

	downloadedFile, err := DownloadFile(imageSource.URL, h.config.DownloadMaxSize, h.config.DownloadTimeout)
	if err != nil {
		slog.Error("Download error", "error", err)
		http.Error(w, "Download error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	downloadedImageType := ImageTypeFromMimeType(downloadedFile.ContentType)
	if downloadedImageType == vips.ImageTypeUnknown {
		slog.Error("Unsupported image format: "+downloadedFile.ContentType, "error", err)
		http.Error(w, "Unsupported image format: "+downloadedFile.ContentType, http.StatusUnprocessableEntity)
		return
	}

	if imageOptions.Format == vips.ImageTypeUnknown {
		imageOptions.Format = downloadedImageType
	}

	processedImage, err := TransformImage(downloadedFile.Buffer.Bytes(), imageOptions)
	if err != nil {
		slog.Error("TransformImage error", "error", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
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
