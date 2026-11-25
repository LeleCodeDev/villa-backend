package audio

import (
	"mime/multipart"

	"github.com/lelecodedev/villa-backend/internal/media"
)

func DefaultOptions() media.Rules {
	return media.Rules{
		MaxSizeBytes: 10 * 1024 * 1024,
		AllowedExts:  []string{"mp3", "wav"},
	}
}

func Save(basePath string, fh *multipart.FileHeader, rules media.Rules) (string, error) {
	ext, err := media.Validate("music", fh, rules)
	if err != nil {
		return "", err
	}

	return media.SaveRaw(basePath, fh, ext)
}
