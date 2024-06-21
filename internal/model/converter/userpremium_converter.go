package converter

import (
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
)

func UserPremiumToResponse(userPremium *entity.UserPremium) *model.UserPremiumResponse {
	return &model.UserPremiumResponse{
		ID:        userPremium.ID,
		StartAt:   userPremium.StartAt,
		EndAt:     userPremium.EndAt,
		CreatedAt: userPremium.CreatedAt,
		UpdatedAt: userPremium.UpdatedAt,
	}
}
