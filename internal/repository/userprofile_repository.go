package repository

import (
	"dealls-dating/internal/entity"

	"gorm.io/gorm"
)

type UserProfileRepository struct {
	Repository[entity.UserProfile]
}

func NewUserProfileRepository() *UserProfileRepository {
	return &UserProfileRepository{}
}

func (u *UserProfileRepository) FindByUserId(db *gorm.DB, entity *entity.UserProfile, userId uint) error {
	return db.Where("user_id = ?", userId).Take(entity).Error
}
