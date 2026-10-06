package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToAttractionResponse(attraction *model.Attraction, transports []model.Transport) dto.AttractionResponse {
	transportResponses := make([]dto.TransportResponse, 0, len(transports))
	for _, transport := range transports {
		transportResponses = append(transportResponses, ToTransportResponse(&transport))
	}

	return dto.AttractionResponse{
		ID:         attraction.ID,
		Title:      attraction.Title,
		Subtitle:   attraction.Subtitle,
		Price:      attraction.Price,
		Transports: transportResponses,
		SortOrder:  attraction.SortOrder,
		Image:      image.BuildURLPtr(attraction.Image),
		CreatedAt:  attraction.CreatedAt,
		UpdatedAt:  attraction.UpdatedAt,
	}
}

func ToAttractionModel(req dto.AttractionRequest, sortOrder int, imagePath *string) *model.Attraction {
	return &model.Attraction{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		Price:     req.Price,
		SortOrder: sortOrder,
		Image:     imagePath,
	}
}

func UpdateAttractionModel(attraction *model.Attraction, req dto.AttractionRequest, sortOrder int, imagePath *string) {
	attraction.Title = req.Title
	attraction.Subtitle = req.Subtitle
	attraction.Price = req.Price
	attraction.Image = imagePath
	attraction.SortOrder = sortOrder
}
