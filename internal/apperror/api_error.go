package appError

import "fmt"

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openwa: status %d: %s", e.StatusCode, e.Body)
}
