package image

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"

	"github.com/lelecodedev/villa-backend/pkg/errors"
)

type SaveOptions struct {
	MaxSizeBytes int64
	MaxWidth     int
	AllowedExts  []string
	JPEGQuality  int
}

func DefaultOptions() SaveOptions {
	return SaveOptions{
		MaxSizeBytes: 5 * 1024 * 1024, // 5MB
		MaxWidth:     1280,
		AllowedExts:  []string{"jpeg", "jpg", "png", "webp"},
		JPEGQuality:  75,
	}
}

func ValidateImage(file *multipart.FileHeader, opts SaveOptions) (image.Image, error) {
	if file.Size > opts.MaxSizeBytes {
		return nil, errors.BadRequest(fmt.Sprintf("File must be under %d MB", opts.MaxSizeBytes/1024/1024))
	}

	ext := filepath.Ext(file.Filename)
	if !slices.Contains(opts.AllowedExts, strings.TrimPrefix(ext, ".")) {
		return nil, errors.BadRequest(fmt.Sprintf("File must be: %v", strings.Join(opts.AllowedExts, ", ")))
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	limitedReader := io.LimitReader(src, opts.MaxSizeBytes)

	img, _, err := image.Decode(limitedReader)
	if err != nil {
		return nil, errors.BadRequest("Uploaded file is not a valid image")
	}

	return img, nil
}

// always save file into jpg for smaller image size
func SaveImage(basePath string, file *multipart.FileHeader, opts SaveOptions) (string, error) {
	img, err := ValidateImage(file, opts)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return "", err
	}

	if img.Bounds().Dx() > opts.MaxWidth {
		img = imaging.Resize(img, opts.MaxWidth, 0, imaging.Lanczos)
	}

	filename := uuid.New().String() + ".jpg"
	path := filepath.Join(basePath, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if err := jpeg.Encode(dst, img, &jpeg.Options{Quality: opts.JPEGQuality}); err != nil {
		return "", err
	}

	return path, nil
}

func DeleteImage(path string) {
	if path != "" {
		os.Remove(path)
	}
}
