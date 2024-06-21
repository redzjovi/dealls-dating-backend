package converter

import (
	"dealls-dating/internal/entity"
	"dealls-dating/internal/model"
)

func UserProfileToResponse(userProfile *entity.UserProfile) *model.UserProfileResponse {
	return &model.UserProfileResponse{
		ID:         userProfile.ID,
		UserID:     userProfile.UserId,
		Gender:     string(userProfile.Gender),
		Name:       userProfile.Name,
		PictureURL: userProfile.PictureURL,
		CreatedAt:  userProfile.CreatedAt,
		UpdatedAt:  userProfile.UpdatedAt,
	}
}
