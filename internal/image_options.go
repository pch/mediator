package internal

import (
	"net/http"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

// Allowed URL parameters
const (
	ParamOperations     = "op"
	ParamWidth          = "w"
	ParamHeight         = "h"
	ParamQuality        = "q"
	ParamStripMetadata  = "strip"
	ParamFormat         = "format"
	ParamPixelateFactor = "pixelatefactor"
)

type ImageOptions struct {
	Operations      []string
	Width           int
	Height          int
	Quality         int
	StripMetadata   bool
	Format          vips.ImageType
	RequestedFormat string
	AutoRotate      bool
	PixelateFactor  int
	Page            int
}

const (
	defaultOperation      = "fit"
	defaultQuality        = 80
	defaultStripMetadata  = true
	defaultPixelateFactor = 20
	defaultPage           = 1

	formatAuto = "auto"
)

func NewImageOptionsFromRequest(r *http.Request) *ImageOptions {
	operations := strings.Split(getQueryParamWithDefault(ParamOperations, defaultOperation, r), ",")
	width := getQueryParamIntWithDefault(ParamWidth, 0, r)
	height := getQueryParamIntWithDefault(ParamHeight, 0, r)
	quality := getQueryParamIntWithDefault(ParamQuality, defaultQuality, r)
	stripMetadata := getQueryParamBoolWithDefault(ParamStripMetadata, defaultStripMetadata, r)
	format := getQueryParamWithDefault(ParamFormat, "", r)
	pixelateFactor := getQueryParamIntWithDefault(ParamPixelateFactor, defaultPixelateFactor, r)
	page := getQueryParamIntWithDefault("page", defaultPage, r)

	var imageType vips.ImageType

	if format == formatAuto {
		imageType = ImageTypeFromAccept(r.Header.Get("Accept"))
	} else {
		imageType = ImageType(format)
	}

	return &ImageOptions{
		Operations:      operations,
		Width:           width,
		Height:          height,
		Quality:         quality,
		StripMetadata:   stripMetadata,
		Format:          imageType,
		RequestedFormat: format,
		AutoRotate:      stripMetadata,
		PixelateFactor:  pixelateFactor,
		Page:            page,
	}
}
