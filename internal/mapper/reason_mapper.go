package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToReasonResponse(reason *model.Reason) dto.ReasonResponse {
	return dto.ReasonResponse{
		ID:        reason.ID,
		Text:      reason.Text,
		CreatedAt: reason.CreatedAt,
		UpdatedAt: reason.UpdatedAt,
	}
}

func ToReasonModel(req dto.ReasonRequest) *model.Reason {
	return &model.Reason{
		Text: req.Text,
	}
}

func UpdateReasonModel(reason *model.Reason, req dto.ReasonRequest) {
	reason.Text = req.Text
}
