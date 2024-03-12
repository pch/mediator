package internal

import (
	"fmt"
	"math"

	"github.com/davidbyttow/govips/v2/vips"
)

type ProcessedImage struct {
	Bytes []byte
	Mime  string
	Size  int
}

var ImageOperationsMap = map[string]ImageOperation{
	"fit":       FitImage,
	"smartcrop": SmartCropImage,
	"pixelate":  PixelateImage,
}

type ImageOperation func(*vips.ImageRef, *ImageOptions) error

func TransformImage(imageBytes []byte, imageOptions *ImageOptions) (*ProcessedImage, error) {
	if len(imageOptions.Operations) == 0 {
		return nil, fmt.Errorf("no operations specified")
	}

	params := vips.NewImportParams()
	params.Page.Set(imageOptions.Page)

	image, err := vips.LoadImageFromBuffer(imageBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	for _, operation := range imageOptions.Operations {
		operationFunc, exists := ImageOperationsMap[operation]

		if !exists {
			return nil, fmt.Errorf("operation not supported: %s", operation)
		}

		err := operationFunc(image, imageOptions)
		if err != nil {
			return nil, err
		}
	}

	return ExportImage(image, imageOptions)
}

func FitImage(image *vips.ImageRef, imageOptions *ImageOptions) error {
	if imageOptions.Width == 0 && imageOptions.Height == 0 {
		// No width or height specified, nothing to do here - and it's not really an error
		return nil
	}

	originalWidth := image.Width()
	originalHeight := image.Height()

	if originalWidth == 0 || originalHeight == 0 {
		return fmt.Errorf("invalid image size")
	}

	finalWidth, finalHeight := fitSizeToLimits(image, imageOptions)

	imageOptions.Width = finalWidth
	imageOptions.Height = finalHeight

	err := image.Thumbnail(imageOptions.Width, imageOptions.Height, vips.InterestingNone)
	if err != nil {
		return err
	}

	return nil
}

func SmartCropImage(image *vips.ImageRef, imageOptions *ImageOptions) error {
	if imageOptions.Width == 0 || imageOptions.Height == 0 {
		return fmt.Errorf("width and height must be specified for smartcrop")
	}

	if imageOptions.Width > image.Width() || imageOptions.Height > image.Height() {
		scale := math.Min(float64(imageOptions.Width)/float64(image.Width()), float64(imageOptions.Height)/float64(image.Height()))
		imageOptions.Width = int(float64(imageOptions.Width) * scale)
		imageOptions.Height = int(float64(imageOptions.Height) * scale)
	}

	err := image.SmartCrop(imageOptions.Width, imageOptions.Height, vips.InterestingAttention)
	if err != nil {
		return err
	}

	return nil
}

func PixelateImage(image *vips.ImageRef, imageOptions *ImageOptions) error {
	if imageOptions.PixelateFactor == 0 {
		return fmt.Errorf("pixelate factor must be specified (non-zero)")
	}
	err := vips.Pixelate(image, float64(imageOptions.PixelateFactor))
	if err != nil {
		return err
	}

	return nil
}

func fitSizeToLimits(image *vips.ImageRef, imageOptions *ImageOptions) (int, int) {
	var originalWidth, originalHeight, fitWidth, fitHeight int

	// Orientation values
	// 0: no EXIF orientation
	// 1: CW 0
	// 2: CW 0, flip horizontal
	// 3: CW 180
	// 4: CW 180, flip horizontal
	// 5: CW 90, flip horizontal
	// 6: CW 270
	// 7: CW 270, flip horizontal
	// 8: CW 90
	if !imageOptions.AutoRotate || image.Orientation() <= 4 {
		originalWidth = image.Width()
		originalHeight = image.Height()
		fitWidth = imageOptions.Width
		fitHeight = imageOptions.Height
	} else {
		originalWidth = image.Height()
		originalHeight = image.Width()
		fitWidth = imageOptions.Height
		fitHeight = imageOptions.Width
	}

	if originalWidth*fitHeight > fitWidth*originalHeight {
		fitHeight = int(math.Round(float64(fitWidth) * float64(originalHeight) / float64(originalWidth)))
	} else {
		fitWidth = int(math.Round(float64(fitHeight) * float64(originalWidth) / float64(originalHeight)))
	}

	return fitWidth, fitHeight
}
