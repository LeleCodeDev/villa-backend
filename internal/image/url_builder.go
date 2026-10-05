package image

import "github.com/lelecodedev/villa-backend/internal/config"

func BuildURL(path string) string {
	if path == "" {
		return ""
	}
	return config.Env.BaseURL + "/api/" + path
}

func BuildURLPtr(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	url := BuildURL(*path)
	return &url
}
