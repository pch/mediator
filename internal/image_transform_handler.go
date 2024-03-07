package internal

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/davidbyttow/govips/v2/vips"
)

type ImageTransformHandler struct {
	config *Config
}

func NewImageTransformHandler(config *Config) *ImageTransformHandler {
	return &ImageTransformHandler{config}
}

func (h *ImageTransformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageSource := getImageSource(r.Context())

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

	imageOptions := NewImageOptionsFromRequest(r)
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

	w.Header().Set("Content-Type", processedImage.Mime)
	w.Header().Set("Content-Length", strconv.Itoa(processedImage.Size))
	w.WriteHeader(http.StatusOK)
	w.Write(processedImage.Bytes)
}
