package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToVillaRuleResponse(villaRule *model.VillaRule) dto.VillaRuleResponse {
	return dto.VillaRuleResponse{
		ID:        villaRule.ID,
		Text:      villaRule.Text,
		SortOrder: villaRule.SortOrder,
		CreatedAt: villaRule.CreatedAt,
		UpdatedAt: villaRule.UpdatedAt,
	}
}

func ToVillaRuleModel(req dto.VillaRuleRequest, sortOrder int) *model.VillaRule {
	return &model.VillaRule{
		Text:      req.Text,
		SortOrder: sortOrder,
	}
}

func UpdateVillaRuleModel(villaRule *model.VillaRule, req dto.VillaRuleRequest, sortOrder int) {
	villaRule.Text = req.Text
	villaRule.SortOrder = sortOrder
}
