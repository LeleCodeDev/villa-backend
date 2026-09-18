package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToFaqResponse(faq *model.Faq) dto.FaqResponse {
	return dto.FaqResponse{
		ID:        faq.ID,
		Question:  faq.Question,
		Answer:    faq.Answer,
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}
}
