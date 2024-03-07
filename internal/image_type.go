package internal

import (
	"mime"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

func ImageType(name string) vips.ImageType {
	switch strings.ToLower(name) {
	case "jpeg":
		return vips.ImageTypeJPEG
	case "png":
		return vips.ImageTypePNG
	case "webp":
		return vips.ImageTypeWEBP
	case "gif":
		return vips.ImageTypeGIF
	default:
		return vips.ImageTypeUnknown
	}
}

func ImageTypeFromMimeType(mimeType string) vips.ImageType {
	mediaType, _, _ := mime.ParseMediaType(mimeType)
	switch mediaType {
	case "image/jpeg":
		return vips.ImageTypeJPEG
	case "image/png":
		return vips.ImageTypePNG
	case "image/webp":
		return vips.ImageTypeWEBP
	case "image/gif":
		return vips.ImageTypeGIF
	default:
		return vips.ImageTypeUnknown
	}
}

func ImageTypeFromAccept(accept string) vips.ImageType {
	for _, v := range strings.Split(accept, ",") {
		mediaType, _, _ := mime.ParseMediaType(v)
		switch mediaType {
		case "image/webp":
			return vips.ImageTypeWEBP
		case "image/png":
			return vips.ImageTypePNG
		case "image/jpeg":
			return vips.ImageTypeJPEG
		}
	}

	return vips.ImageTypeUnknown
}

func MimeTypeFromImageType(code vips.ImageType) string {
	switch code {
	case vips.ImageTypePNG:
		return "image/png"
	case vips.ImageTypeWEBP:
		return "image/webp"
	case vips.ImageTypeGIF:
		return "image/gif"
	default:
		return "image/jpeg"
	}
}
