package internal

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

var ImageExportMap = map[vips.ImageType]ImageExport{
	vips.ImageTypeJPEG: ExportJPEG,
	vips.ImageTypePNG:  ExportPNG,
	vips.ImageTypeWEBP: ExportWEBP,
	vips.ImageTypeGIF:  ExportGIF,
}

type ImageExport func(*vips.ImageRef, *ImageOptions) ([]byte, error)

func ExportImage(image *vips.ImageRef, imageOptions *ImageOptions) (*ProcessedImage, error) {
	if imageOptions.AutoRotate {
		image.AutoRotate()
	}

	if imageOptions.Format == vips.ImageTypeUnknown {
		imageOptions.Format = vips.ImageTypeJPEG
	}

	exportFunc, exists := ImageExportMap[imageOptions.Format]
	if !exists {
		return nil, fmt.Errorf("format not supported: %s", fmt.Sprint(imageOptions.Format))
	}

	fileBytes, err := exportFunc(image, imageOptions)
	if err != nil {
		return nil, err
	}

	mime := MimeTypeFromImageType(vips.DetermineImageType(fileBytes))

	return &ProcessedImage{
		Bytes: fileBytes,
		Mime:  mime,
		Size:  len(fileBytes),
	}, nil
}

func ExportJPEG(image *vips.ImageRef, imageOptions *ImageOptions) ([]byte, error) {
	ep := vips.NewJpegExportParams()
	ep.StripMetadata = imageOptions.StripMetadata
	ep.Quality = imageOptions.Quality
	ep.OptimizeCoding = true
	ep.SubsampleMode = vips.VipsForeignSubsampleAuto
	ep.TrellisQuant = true
	ep.OvershootDeringing = true
	ep.OptimizeScans = true
	ep.QuantTable = 3

	fileBytes, _, err := image.ExportJpeg(ep)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func ExportPNG(image *vips.ImageRef, imageOptions *ImageOptions) ([]byte, error) {
	ep := vips.NewPngExportParams()
	ep.StripMetadata = imageOptions.StripMetadata
	ep.Quality = imageOptions.Quality

	fileBytes, _, err := image.ExportPng(ep)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func ExportWEBP(image *vips.ImageRef, imageOptions *ImageOptions) ([]byte, error) {
	ep := vips.NewWebpExportParams()
	ep.StripMetadata = imageOptions.StripMetadata
	ep.Quality = imageOptions.Quality

	fileBytes, _, err := image.ExportWebp(ep)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func ExportGIF(image *vips.ImageRef, imageOptions *ImageOptions) ([]byte, error) {
	ep := vips.NewGifExportParams()
	ep.StripMetadata = imageOptions.StripMetadata
	ep.Quality = imageOptions.Quality

	fileBytes, _, err := image.ExportGIF(ep)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
