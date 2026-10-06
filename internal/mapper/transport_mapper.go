package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToTransportResponse(transport *model.Transport) dto.TransportResponse {
	return dto.TransportResponse{
		ID:           transport.ID,
		Vehicle:      transport.Vehicle,
		Distance:     transport.Distance,
		TimeRequired: transport.TimeRequired,
		SortOrder:    transport.SortOrder,
		CreatedAt:    transport.CreatedAt,
		UpdatedAt:    transport.UpdatedAt,
	}
}

func ToTransportModel(req dto.TransportRequest, attraction model.Attraction, sortOrder int) *model.Transport {
	return &model.Transport{
		Vehicle:      req.Vehicle,
		Distance:     req.Distance,
		TimeRequired: req.TimeRequired,
		Attraction:   attraction,
		AttractionID: attraction.ID,
		SortOrder:    sortOrder,
	}
}

func UpdateTransportModel(transport *model.Transport, attraction model.Attraction, req dto.TransportRequest, sortOrder int) {
	transport.Vehicle = req.Vehicle
	transport.Distance = req.Distance
	transport.TimeRequired = req.TimeRequired
	transport.AttractionID = attraction.ID
	transport.Attraction = attraction
	transport.SortOrder = sortOrder
}
