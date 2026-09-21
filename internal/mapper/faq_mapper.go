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
		SortOrder: faq.SortOrder,
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}
}

func ToFaqModel(req dto.FaqRequest, sortOrder int) *model.Faq {
	return &model.Faq{
		Question:  req.Question,
		Answer:    req.Answer,
		SortOrder: sortOrder,
	}
}

func UpdateFaqModel(faq *model.Faq, req dto.FaqRequest, sortOrder int) {
	faq.Question = req.Question
	faq.Answer = req.Answer
	faq.SortOrder = sortOrder
}
