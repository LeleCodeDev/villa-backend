package image

import (
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"

	"github.com/disintegration/imaging"
	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/media"
	_ "golang.org/x/image/webp"
)

type SaveOptions struct {
	media.Rules
	MaxWidth    int
	JPEGQuality int
}

func DefaultOptions() SaveOptions {
	return SaveOptions{
		Rules: media.Rules{
			MaxSizeBytes: 5 * 1024 * 1024,
			AllowedExts:  []string{"jpeg", "jpg", "png", "webp"},
		},
		MaxWidth:    1280,
		JPEGQuality: 75,
	}
}

func decode(fh *multipart.FileHeader) (image.Image, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, appError.BadRequest("Uploaded file is not a valid image")
	}
	return img, nil
}

func Save(basePath string, fh *multipart.FileHeader, opts SaveOptions) (string, error) {
	if _, err := media.Validate("image", fh, opts.Rules); err != nil {
		return "", err
	}

	img, err := decode(fh)
	if err != nil {
		return "", err
	}

	if img.Bounds().Dx() > opts.MaxWidth {
		img = imaging.Resize(img, opts.MaxWidth, 0, imaging.Lanczos)
	}

	return media.Write(basePath, media.UniqueName("jpg"), func(w io.Writer) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: opts.JPEGQuality})
	})
}
