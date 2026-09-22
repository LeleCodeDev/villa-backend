package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToTestimonyResponse(testimony *model.Testimony) dto.TestimonyResponse {
	return dto.TestimonyResponse{
		ID:        testimony.ID,
		Star:      testimony.Star,
		Comment:   testimony.Comment,
		Username:  testimony.Username,
		SortOrder: testimony.SortOrder,
		CreatedAt: testimony.CreatedAt,
		UpdatedAt: testimony.UpdatedAt,
	}
}

func ToTestimonyModel(req dto.TestimonyRequest, sortOrder int) *model.Testimony {
	return &model.Testimony{
		Star:      req.Star,
		Comment:   req.Comment,
		Username:  req.Username,
		SortOrder: sortOrder,
	}
}

func UpdateTestimonyModel(testimony *model.Testimony, req dto.TestimonyResponse, sortOrder int) {
	testimony.Star = req.Star
	testimony.Comment = req.Comment
	testimony.Username = req.Username
	testimony.SortOrder = sortOrder
}
