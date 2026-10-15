// BSD 3-Clause License: https://github.com/prplecake/go-thumbnail
package media

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

var (
	// ErrInvalidMimeType is returned when a non-image content type is detected.
	ErrInvalidMimeType = errors.New("invalid mimetype")

	// ErrInvalidScaler is returned when an unrecognized scaler is configured.
	ErrInvalidScaler = errors.New("invalid scaler")

	// ErrInvalidDimensions is returned when either target dimension is not positive.
	ErrInvalidDimensions = errors.New("invalid image dimensions")
)

// ImageGenerator configures thumbnail generation.
type ImageGenerator struct {
	// Width is the destination thumbnail width.
	Width int

	// Height is the destination thumbnail height.
	Height int

	// Scaler selects the interpolation algorithm.
	Scaler string

	// CropToFit preserves the source aspect ratio by cropping the image from
	// its center until it matches the destination aspect ratio.
	CropToFit bool
}

// NewImageGenerator applies defaults to a thumbnail generator configuration.
func NewImageGenerator(config ImageGenerator) *ImageGenerator {
	if config.Width == 0 {
		config.Width = 300
	}
	if config.Height == 0 {
		config.Height = 300
	}
	if config.Scaler == "" {
		config.Scaler = "CatmullRom"
	}
	return &config
}

// ResizeImage resizes an encoded jpg / png from source into target.
func ResizeImage(source io.Reader, target io.Writer, width int, height int) error {
	generator := &ImageGenerator{
		Width:  width,
		Height: height,
		Scaler: "CatmullRom",
	}
	return generator.CreateThumbnail(source, target, "")
}

// ResizeAndCropImage resizes an encoded jpg / png while preserving its aspect ratio.
// The source is cropped from its center to fill the target.
func ResizeAndCropImage(source io.Reader, target io.Writer, width int, height int) error {
	generator := &ImageGenerator{
		Width:     width,
		Height:    height,
		Scaler:    "CatmullRom",
		CropToFit: true,
	}
	return generator.CreateThumbnail(source, target, "")
}

// ResizeImageToPng resizes an encoded image into target and always
// encodes the result as png, regardless of the source format.
func ResizeImageToPng(source io.Reader, target io.Writer, width int, height int) error {
	generator := &ImageGenerator{
		Width:  width,
		Height: height,
		Scaler: "CatmullRom",
	}
	return generator.CreateThumbnail(source, target, "png")
}

// ResizeImageBytes resizes a jpg / png byte slice and
// returns the result in the source format.
func ResizeImageBytes(source []byte, width int, height int) ([]byte, error) {
	var target bytes.Buffer
	if err := ResizeImage(bytes.NewReader(source), &target, width, height); err != nil {
		return nil, err
	}
	return target.Bytes(), nil
}

// ResizeImageBytesToPng resizes an encoded image byte slice
// and always returns the result encoded as png.
func ResizeImageBytesToPng(source []byte, width int, height int) ([]byte, error) {
	var target bytes.Buffer
	if err := ResizeImageToPng(bytes.NewReader(source), &target, width, height); err != nil {
		return nil, err
	}
	return target.Bytes(), nil
}

// CreateThumbnail decodes an image from source, resizes it,
// and writes the encoded result to target.
func (generator *ImageGenerator) CreateThumbnail(source io.Reader, target io.Writer, outputFormat string) error {
	input, inputFormat, err := image.Decode(source)
	if errors.Is(err, image.ErrFormat) {
		return ErrInvalidMimeType
	}
	if err != nil {
		return err
	}

	thumbnail, err := generator.createRect(input)
	if err != nil {
		return err
	}
	if outputFormat == "" {
		outputFormat = inputFormat
	}

	switch outputFormat {
	case "jpeg":
		return jpeg.Encode(target, thumbnail, nil)
	case "png":
		return png.Encode(target, thumbnail)
	default:
		return ErrInvalidMimeType
	}
}

func (generator *ImageGenerator) createRect(input image.Image) (*image.RGBA, error) {
	if generator.Width <= 0 || generator.Height <= 0 {
		return nil, ErrInvalidDimensions
	}

	scaler := resolveImageScaler(generator.Scaler)
	if scaler == nil {
		return nil, ErrInvalidScaler
	}

	targetRect := image.Rect(0, 0, generator.Width, generator.Height)
	target := image.NewRGBA(targetRect)

	sourceRect := input.Bounds()
	if generator.CropToFit {
		sourceRect = cropToAspectRatio(sourceRect, generator.Width, generator.Height)
	}
	scaler.Scale(target, targetRect, input, sourceRect, draw.Over, nil)
	return target, nil
}

func resolveImageScaler(name string) draw.Interpolator {
	switch name {
	case "NearestNeighbor":
		return draw.NearestNeighbor
	case "ApproxBiLinear":
		return draw.ApproxBiLinear
	case "BiLinear":
		return draw.BiLinear
	case "CatmullRom":
		return draw.CatmullRom
	default:
		return nil
	}
}

func cropToAspectRatio(source image.Rectangle, targetWidth int, targetHeight int) image.Rectangle {
	sourceWidth := source.Dx()
	sourceHeight := source.Dy()

	// The source is wider than the target -> crop its left and right edges
	if int64(sourceWidth)*int64(targetHeight) > int64(sourceHeight)*int64(targetWidth) {
		width := int(int64(sourceHeight) * int64(targetWidth) / int64(targetHeight))
		offset := (sourceWidth - width) / 2
		return image.Rect(source.Min.X+offset, source.Min.Y, source.Min.X+offset+width, source.Max.Y)
	}

	// The source is taller than the target -> crop its top and bottom edges
	height := int(int64(sourceWidth) * int64(targetHeight) / int64(targetWidth))
	offset := (sourceHeight - height) / 2
	return image.Rect(source.Min.X, source.Min.Y+offset, source.Max.X, source.Min.Y+offset+height)
}
