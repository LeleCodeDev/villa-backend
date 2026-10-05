// Package pagination
package pagination

import (
	"math"

	"github.com/lelecodedev/villa-backend/internal/response"
)

func BuildPagination(page int, size int, total int64) response.PaginationData {
	totalPages := 0
	if page > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(size)))
	}

	return response.PaginationData{
		Page:       page,
		Size:       size,
		TotalItems: total,
		TotalPages: totalPages,
	}
}
