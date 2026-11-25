package media

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
	appError "github.com/lelecodedev/villa-backend/internal/apperror"
)

type Rules struct {
	MaxSizeBytes int64
	AllowedExts  []string
}

func Validate(field string, fh *multipart.FileHeader, r Rules) (string, error) {
	if fh.Size > r.MaxSizeBytes {
		return "", appError.BadRequest(fmt.Sprintf("%s must be under %d MB", field, r.MaxSizeBytes/1024/1024))
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fh.Filename), "."))
	if !slices.Contains(r.AllowedExts, ext) {
		return "", appError.BadRequest(fmt.Sprintf("%s must be: %s", field, strings.Join(r.AllowedExts, ", ")))
	}
	return ext, nil
}

func UniqueName(ext string) string {
	return uuid.NewString() + "." + ext
}

func Write(dir, name string, write func(io.Writer) error) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, name)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	if err := write(dst); err != nil {
		dst.Close()
		os.Remove(path)
		return "", err
	}

	if err := dst.Close(); err != nil {
		os.Remove(path)
		return "", err
	}

	return path, nil
}

func SaveRaw(dir string, fh *multipart.FileHeader, ext string) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	return Write(dir, UniqueName(ext), func(w io.Writer) error {
		_, err := io.Copy(w, src)
		return err
	})
}

func Delete(path string) error {
	if path == "" {
		return nil
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
