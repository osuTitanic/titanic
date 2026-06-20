// MIT License: https://github.com/prplecake/go-thumbnail
// Package thumbnail provides a method to create thumbnails from images.
package media

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"golang.org/x/image/draw"
)

// An Image is an image and information about it.
type Image struct {
	// Path is a path to an image.
	Path string

	// ContentType is the content type of the image.
	ContentType string

	// Data is the image data in a byte-array
	Data []byte

	// Size is the length of Data
	Size int

	// Current stores the existing image's dimensions
	Current ImageDimensions

	// Future store the new thumbnail dimensions.
	Future ImageDimensions
}

// ImageDimensions stores dimensional information for an Image.
type ImageDimensions struct {
	// Width is the width of an image in pixels.
	Width int

	// Height is the height on an image in pixels.
	Height int

	// X is the right-most X-coordinate.
	X int

	// Y is the top-most Y-coordinate.
	Y int
}

var (
	// ErrInvalidMimeType is returned when a non-image content type is
	// detected.
	ErrInvalidMimeType = errors.New("invalid mimetype")

	// ErrInvalidScaler is returned when an unrecognized scaler is
	// passed to the Generator.
	ErrInvalidScaler = errors.New("invalid scaler")
)

// NewImageGenerator returns an instance of a thumbnail generator with a
// given configuration.
func NewImageGenerator(c ImageGenerator) *ImageGenerator {
	return &ImageGenerator{
		Width:             300,
		Height:            300,
		DestinationPath:   c.DestinationPath,
		DestinationPrefix: c.DestinationPrefix,
		Scaler:            c.Scaler,
	}
}

// NewImageFromFile reads in an image file from the file system and
// populates an Image object. That new Image object is returned along
// with any errors that occur during the operation.
func (gen *ImageGenerator) NewImageFromFile(path string) (*Image, error) {
	imageBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentType := detectContentType(imageBytes)
	return &Image{
		Path:        path,
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
		Current: ImageDimensions{
			Width:  0,
			Height: 0,
		},
		Future: ImageDimensions{
			Width:  gen.Width,
			Height: gen.Height,
		},
	}, nil
}

// NewImageFromByteArray reads in an image from a byte array and
// populates an Image object. That new Image object is returned along
// with any errors that occur during the operation.
func (gen *ImageGenerator) NewImageFromByteArray(imageBytes []byte) (*Image, error) {
	contentType := detectContentType(imageBytes)
	return &Image{
		ContentType: contentType,
		Data:        imageBytes,
		Size:        len(imageBytes),
		Current: ImageDimensions{
			Width:  0,
			Height: 0,
		},
		Future: ImageDimensions{
			Width:  gen.Width,
			Height: gen.Height,
		},
	}, nil
}

// ImageGenerator registers a generator configuration to be used when
// creating thumbnails.
type ImageGenerator struct {
	// Width is the destination thumbnail width.
	Width int

	// Height is the destination thumbnail height.
	Height int

	// DestinationPath is the destination thumbnail path.
	DestinationPath string

	// DestinationPrefix is the prefix for the destination thumbnail
	// filename.
	DestinationPrefix string

	// Scaler is the scaler to be used when generating thumbnails.
	Scaler string
}

// ResizeImage resizes the given image data to the provided
// dimensions and returns the encoded result.
func ResizeImage(data []byte, width int, height int) ([]byte, error) {
	generator := &ImageGenerator{
		Width:  width,
		Height: height,
		Scaler: "CatmullRom",
	}

	image, err := generator.NewImageFromByteArray(data)
	if err != nil {
		return nil, err
	}

	return generator.CreateThumbnail(image)
}

// CreateThumbnail generates a thumbnail.
func (gen *ImageGenerator) CreateThumbnail(i *Image) ([]byte, error) {
	if i.ContentType == "application/octet-stream" {
		return nil, ErrInvalidMimeType
	}

	dst, err := gen.createRect(i)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	switch i.ContentType {
	case "image/jpeg":
		err = jpeg.Encode(&buffer, dst, nil)
	case "image/png":
		err = png.Encode(&buffer, dst)
	}
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (gen *ImageGenerator) createRect(i *Image) (*image.RGBA, error) {
	img, _, err := image.Decode(bytes.NewReader(i.Data))
	if err != nil {
		return nil, err
	}
	rect := image.Rect(0, 0, gen.Width, gen.Height)
	dst := image.NewRGBA(rect)
	var scaler draw.Interpolator
	switch scalerChoice := gen.Scaler; scalerChoice {
	case "NearestNeighbor":
		scaler = draw.NearestNeighbor
	case "ApproxBiLinear":
		scaler = draw.ApproxBiLinear
	case "BiLinear":
		scaler = draw.BiLinear
	case "CatmullRom":
		scaler = draw.CatmullRom
	}
	if scaler == nil {
		return nil, ErrInvalidScaler
	}
	scaler.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst, nil

}

// detectContentType from
// https://golangcode.com/get-the-content-type-of-file/
func detectContentType(fb []byte) string {
	// Only the first 512 bytes are used to sniff the content type.
	// Use the net/http package's handy DetectContentType function.
	// Always seems to return a valid content-type by returning
	// "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(fb[:512])
}
