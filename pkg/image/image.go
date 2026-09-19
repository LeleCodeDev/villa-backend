package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"uuid"

	"github.com/disintegration/imaging"
	"github.com/lelecodedev/villa-backend/pkg/errors"
)

const (
	maxWidth    = 1280
	jpegQuality = 75
)

func SaveImage(basePath string, file *multipart.FileHeader) (*string, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, err
	}

	allowedExts := []string{"jpeg", "png", "jpg"}
	ext := filepath.Ext(file.Filename)
	if !slices.Contains(allowedExts, strings.TrimPrefix(ext, ".")) {
		return nil, errors.BadRequest(fmt.Sprintf("File must be: %v", strings.Join(allowedExts, ", ")))
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, errors.BadRequest("Uploaded file is not a valid image")
	}

	if img.Bounds().Dx() > maxWidth {
		img = imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
	}

	filename := uuid.New().String() + ".jpg"
	path := filepath.Join(basePath, filename)

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if err := jpeg.Encode(dst, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, err
	}

	return &path, nil
}

func DeleteImage(path string) {
	if path != "" {
		os.Remove(path)
	}
}
